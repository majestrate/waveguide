package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"waveguide/lib/log"
	"waveguide/lib/oauth"
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
	thumb, _ := url.Parse(oldest[:len(oldest)-7] + "jpeg")
	cdnURL, _ := url.Parse(s.conf.Worker.UploadURL)
	utorrent.Scheme = "http"
	ufile.Scheme = "http"
	thumb.Scheme = "http"
	utorrent.Host = cdnURL.Host
	ufile.Host = cdnURL.Host
	thumb.Host = cdnURL.Host
	api.DoHTTP(api.DeleteRequest(utorrent))
	api.DoHTTP(api.DeleteRequest(ufile))
	api.DoHTTP(api.DeleteRequest(thumb))
}

func (s *Server) APIStreamPublish(c *gin.Context) {
	if s.Anon() {
		s.ctx.Ensure("1", "anon", "5")
		return
	}
	user, token := extractUserToken(c.PostForm("name"))
	if user != "" && token != "" {
		u, err := s.oauth.GetUser(token)
		if err == nil {
			if user == u.ID {
				chatid, err := s.oauth.EnsureChat(token)
				if err == nil {
					s.ctx.Ensure(user, u.Username, chatid)
				} else {
					log.Errorf("failed to create chat for stream: %s", err.Error())
					c.String(http.StatusInternalServerError, err.Error())
				}
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
		s.ctx.Remove(user)
		for _, u := range info.URLS[:] {
			if u != "" {
				s.deleteTorrent(u)
			}
		}
		s.oauth.SubmitPost(token, info.ChatID, oauth.Post{
			Text: fmt.Sprintf("%s has ended streaming, press F to pay respects", info.Username),
		})
		s.oauth.StreamOffline(token, info.ID)
	}
}

func (s *Server) APIStreamSegment(c *gin.Context) {
	user, token := extractUserToken(c.PostForm("name"))
	infile := c.PostForm("path")
	info := s.ctx.Find(user)
	if info != nil && info.Segments == 0 && s.oauth != nil {
		// got first segment
		s.oauth.SubmitPost(token, info.ChatID, oauth.Post{
			Text: "aw yeh, now heckin' streaming at http://gitgud.tv/watch/?u=" + user,
		})
		s.oauth.StreamOnline(token, user)
	}
	if info == nil {
		log.Errorf("non existing stream %s", user)
		os.Remove(infile)
		c.String(http.StatusNotFound, "no such stream")
		return
	} else {
		info.Segments++
	}

	outfile := util.TempFileName(os.TempDir(), fmt.Sprintf("-%s-stream.mp4", user))
	defer os.Remove(outfile)
	defer os.Remove(infile)
	err := s.encoder.Transcode(infile, outfile)
	if err != nil {
		log.Errorf("failed to transcode file: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	thumbURL, _ := url.Parse(s.MakeThumbnailUploadURL(outfile))
	thumbfile := filepath.Join(os.TempDir(), filepath.Base(thumbURL.Path))
	defer os.Remove(thumbfile)
	err = s.encoder.Thumbnail(outfile, thumbfile)
	if err != nil {
		log.Errorf("failed to generate stream thumbnail: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	videoURL, _ := url.Parse(s.MakeVideoUploadUrl(outfile))
	torrentURL, _ := url.Parse(s.MakeTorrentUploadURL(outfile))
	webseedURL := s.MakeWebseedURL(outfile)
	publicURL := s.MakePublicURL(outfile)

	if !s.ctx.Has(info.ID) {
		// stream is gone
		log.Errorf("stream %s is gone, will not publish files", info.ID)
		c.String(http.StatusNotFound, "stream is gone")
		return
	}

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
		// upload thumbnail
		var thumb *os.File
		thumb, err = os.Open(thumbfile)
		if err == nil {
			err = api.DoHTTP(api.UploadRequest(thumbURL, thumb))
			thumb.Close()
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
	} else {
		log.Errorf("api error: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
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
