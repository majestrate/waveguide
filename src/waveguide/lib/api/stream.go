package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"waveguide/lib/log"
	"waveguide/lib/util"
	"waveguide/lib/worker/api"
)

const StreamKeyDelim = "|"

func extractUserToken(streamkey string) (user, token string) {
	parts := strings.Split(streamkey, StreamKeyDelim)
	if len(parts) == 2 {
		user, token = parts[0], parts[1]
	}
	return
}

func (s *Server) deleteTorrent(oldest string) {
	utorrent, _ := url.Parse(oldest)
	ufile, _ := url.Parse(oldest[:len(oldest)-8])
	cdnURL, _ := url.Parse(s.conf.Worker.UploadURL)
	utorrent.Scheme = "http"
	ufile.Scheme = "http"
	utorrent.Host = cdnURL.Host
	ufile.Host = cdnURL.Host
	api.DoHTTP(api.DeleteRequest(utorrent))
	api.DoHTTP(api.DeleteRequest(ufile))
}

func (s *Server) APIStreamPublish(c *gin.Context) {
	user, token := extractUserToken(c.PostForm("name"))
	if user != "" && token != "" {
		err := s.oauth.AnnounceStream(token, "streaming live at https://gitgud.tv/watch/?u="+user)
		if err == nil {
			s.ctx.Ensure(user)
		} else {
			c.String(http.StatusForbidden, err.Error())
		}
	} else {
		c.String(http.StatusForbidden, "")
	}
}

func (s *Server) APIStreamJoin(c *gin.Context) {
}

func (s *Server) APIStreamPart(c *gin.Context) {
}

func (s *Server) APIStreamDone(c *gin.Context) {
	user, _ := extractUserToken(c.PostForm("name"))
	info := s.ctx.Find(user)
	if info != nil {
		for _, u := range info.URLS[:] {
			if u != "" {
				s.deleteTorrent(u)
			}
		}
	}
}

func (s *Server) APIStreamSegment(c *gin.Context) {
	user, _ := extractUserToken(c.PostForm("name"))
	infile := c.PostForm("path")
	ext := filepath.Ext(infile)
	outfile := util.TempFileName(os.TempDir(), ext)
	err := os.Rename(infile, outfile)
	if err != nil {
		log.Errorf("failed to move file: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer os.Remove(outfile)
	defer os.Remove(infile)
	videoURL, _ := url.Parse(s.MakeVideoUploadUrl(outfile))
	torrentURL, _ := url.Parse(s.MakeTorrentUploadURL(outfile))
	webseedURL := s.MakeWebseedURL(outfile)
	publicURL := s.MakePublicURL(outfile)

	var videoFile *os.File
	videoFile, err = os.Open(outfile)
	if err == nil {
		defer videoFile.Close()
		torrent := new(util.Buffer)
		err = s.torrent.MakeSingleWithWebseed(filepath.Base(outfile), webseedURL, videoFile, torrent)
		if err == nil {
			videoFile.Seek(0, 0)
			err = api.DoHTTP(api.UploadRequest(videoURL, videoFile))
			if err == nil {
				err = api.DoHTTP(api.UploadRequest(torrentURL, torrent))
			}
		}
	}
	if err == nil {
		info := s.ctx.Find(user)
		if info == nil {
			log.Errorf("failed to update non existing stream %s", user)
		} else {
			oldest := info.OldestTorrent()
			if oldest != "" {
				s.deleteTorrent(oldest)
			}
			info.AddTorrent(publicURL)
		}
	}
}

func (s *Server) APIStreamInfo(c *gin.Context) {
	u := c.Query("u")
	info := s.ctx.Find(u)
	c.JSON(http.StatusOK, info)
}
