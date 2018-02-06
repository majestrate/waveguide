package frontend

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"waveguide/lib/model"
	"waveguide/lib/util"
)

var ErrBadMediaType = errors.New("bad media type")
var ErrBadContentType = errors.New("bad content type")
var ErrNoWebseedFileName = errors.New("no title provided")

func NewUpload(u *model.UserInfo) *model.VideoInfo {
	return &model.VideoInfo{
		UserID:     u.UserID,
		UploadedAt: time.Now().Unix(),
	}
}

func (r *Routes) ApiUpload(c *gin.Context) {
	var err error
	var videoURL string
	u := r.GetCurrentUser(c)
	info := NewUpload(u)
	webseed, ok := c.GetPostForm("webseed")
	if ok {
		info.Title, ok = c.GetPostForm("title")
		if ok {
			err = r.DB.RegisterVideo(info)
			if err == nil {
				videoURL = info.GetURL(r.FrontendURL).String()
				var webseedURL *url.URL
				webseedURL, err = url.Parse(webseed)
				if err == nil {
					err = r.DB.AddWebseed(info.VideoID, webseedURL.String())
					if err == nil {
						err = r.api.Do(info.WebseedUploadRequest(webseedURL))
					}
				}
			}
		} else {
			err = ErrNoWebseedFileName
		}
	} else {
		video, err := c.FormFile("video")
		if err == nil {
			info.Title = video.Filename
			var f multipart.File
			f, err = video.Open()
			if err == nil {
				videoURL, err = r.UploadVideoFile(video.Filename, f, info)
				f.Close()
			}
		}
	}
	if err == nil {
		c.JSON(http.StatusCreated, map[string]interface{}{
			"url":   videoURL,
			"error": nil,
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"url":   videoURL,
			"error": err.Error(),
		})
	}
}

func (r *Routes) UploadVideoFile(filename string, body io.Reader, info *model.VideoInfo) (videoURL string, err error) {
	ctype := mime.TypeByExtension(filepath.Ext(filename))
	if ctype != "" {
		mtype, _, _ := mime.ParseMediaType(ctype)
		if strings.HasPrefix(mtype, "video/") {
			err = r.DB.RegisterVideo(info)
			if err == nil {
				videoURL = info.GetURL(r.FrontendURL).String()
				ext := filepath.Ext(filename)
				tmpFile := util.TempFileName(r.TempDir, ext)
				var osf *os.File
				osf, err = os.Create(tmpFile)
				if err == nil {
					var buff [65536]byte
					_, err = io.CopyBuffer(osf, body, buff[:])
					osf.Close()

					fileURL := &url.URL{
						Scheme: "file",
						Path:   tmpFile,
					}
					if err == nil {
						err = r.api.Do(info.VideoUploadRequest(fileURL, filename))
					}
				}
			}
		} else {
			err = ErrBadMediaType
		}
	} else {
		err = ErrBadContentType
	}
	return
}
