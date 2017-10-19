package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
	"waveguide/lib/model"
	"waveguide/lib/util"
)

func (r *Routes) HandleUpload(c *gin.Context) {
	var videoURL string
	u := r.GetCurrentUser(c)
	video, err := c.FormFile("video")
	if err == nil {
		var vidID int64
		vidID, err = r.DB.NextVideoID()
		if err == nil {

			info := &model.VideoInfo{
				UserID:     u.UserID,
				VideoID:    vidID,
				Title:      video.Filename,
				UploadedAt: time.Now().Unix(),
			}
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
