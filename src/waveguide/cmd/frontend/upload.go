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

func (r *Routes) HandleUpload(c *gin.Context) {
	var videoURL string
	u := r.GetCurrentUser(c)
	video, err := c.FormFile("video")
	if err == nil {
		ctype := mime.TypeByExtension(filepath.Ext(video.Filename))
		if ctype != "" {
			mtype, _, _ := mime.ParseMediaType(ctype)
			if strings.HasPrefix(mtype, "video/") {
				var vidID int64
				vidID, err = r.DB.NextVideoID()
				if err == nil {

					info := &model.VideoInfo{
						UserID:     u.UserID,
						VideoID:    vidID,
						Title:      video.Filename,
						UploadedAt: time.Now().Unix(),
					}
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
