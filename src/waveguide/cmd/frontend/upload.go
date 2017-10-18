package frontend

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
	"waveguide/lib/model"
)

func (r *Routes) HandleUpload(c *gin.Context) {
	var videoURL string
	u := r.GetCurrentUser(c)
	video, err := c.FormFile("video")
	if err == nil {
		var body io.ReadCloser
		body, err = video.Open()
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
				videoURL = info.GetURL(r.FrontendURL).String()
				body, err = video.Open()
				if err == nil {
					err = r.api.Do(info.VideoUploadRequest(r.workerURL, video.Filename, body))
				}
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
