package frontend

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
	"waveguide/lib/model"
	"waveguide/lib/util"
)

var ErrBadMediaType = errors.New("bad media type")
var ErrBadContentType = errors.New("bad content type")
var ErrNoWebseedFileName = errors.New("no title provided")

func (r *Routes) ApiUpload(c *gin.Context) {
	var err error
	var videoURL string
	u := r.GetCurrentUser(c)
	info := &model.VideoInfo{
		UserID:     u.UserID,
		UploadedAt: time.Now().Unix(),
	}
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
			ctype := mime.TypeByExtension(filepath.Ext(video.Filename))
			if ctype != "" {
				mtype, _, _ := mime.ParseMediaType(ctype)
				if strings.HasPrefix(mtype, "video/") {

					info.Title = video.Filename
					err = r.DB.RegisterVideo(info)
					if err == nil {
						videoURL = info.GetURL(r.FrontendURL).String()
						ext := filepath.Ext(video.Filename)
						tmpFile := util.TempFileName(r.TempDir, ext)
						c.SaveUploadedFile(video, tmpFile)
						fileURL := &url.URL{
							Scheme: "file",
							Path:   tmpFile,
						}
						if err == nil {
							err = r.api.Do(info.VideoUploadRequest(fileURL, video.Filename))
						}
					}
				} else {
					err = ErrBadMediaType
				}
			} else {
				err = ErrBadContentType
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
