package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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
		u, err := s.oauth.GetUser(token)
		if err == nil {
			if user == u.ID {
				s.ctx.Ensure(user, u.Username)
				s.oauth.AnnounceStream(token, "now live streaming at http://gitgud.tv/watch/?u="+user)
			} else {
				log.Errorf("user id missmatch, '%s' != '%s'", user, u.ID)
				c.String(http.StatusForbidden, "")
			}
		} else {
			log.Errorf("failed to auth stream with oauth: %s", err.Error())
			c.String(http.StatusForbidden, err.Error())
		}
	} else {
		c.String(http.StatusForbidden, "")
	}
}

func (s *Server) APIStreamJoin(c *gin.Context) {
	// TODO: implement
}

func (s *Server) APIStreamPart(c *gin.Context) {
	// TODO: implement
}

func (s *Server) APIStreamDone(c *gin.Context) {
	user, token := extractUserToken(c.PostForm("name"))
	info := s.ctx.Find(user)
	if info != nil {
		for _, u := range info.URLS[:] {
			if u != "" {
				s.deleteTorrent(u)
			}
		}
		s.oauth.AnnounceStream(token, "stream is now offline, bai.")
		s.ctx.Remove(user)
	}
}

func (s *Server) APIStreamSegment(c *gin.Context) {
	user, _ := extractUserToken(c.PostForm("name"))
	infile := c.PostForm("path")
	outfile := util.TempFileName(os.TempDir(), ".mp4")
	defer os.Remove(outfile)
	defer os.Remove(infile)
	err := s.encoder.Transcode(infile, outfile)
	if err != nil {
		log.Errorf("failed to transcode file: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
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
			// delete oldest torrent first
			oldest := info.OldestTorrent()
			if oldest != "" {
				s.deleteTorrent(oldest)
			}
			// then add new torrent
			info.AddTorrent(publicURL)
		}
	}
}

func (s *Server) APIStreamInfo(c *gin.Context) {
	u := c.Param("key")
	info := s.ctx.Find(u)
	if info == nil {
		// stream not found
		c.JSON(http.StatusNotFound, map[string]interface{}{})
	} else {
		c.JSON(http.StatusOK, info)
	}
}

func (s *Server) APIListStreams(c *gin.Context) {
	limit := 10
	limit_str := c.Query("limit")
	if limit_str != "" {
		li, err := strconv.Atoi(limit_str)
		if err == nil {
			limit = li
		}
	}
	streams := s.ctx.Online(limit)
	c.JSON(http.StatusOK, streams)
}
