package frontend

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
	"waveguide/lib/model"
)

func (r *Routes) HandleUpload(c *gin.Context) {
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
				callbackURL := r.VideoDoneCallbackURL(info)
				err = r.api.UploadVideo(callbackURL.String(), video.Filename, body)
			}
		}
	}
	if err == nil {
		c.HTML(http.StatusOK, "upload_processing.html", map[string]interface{}{})
	} else {
		r.Error(c, err)
	}
}
