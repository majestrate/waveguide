package frontend

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (r *Routes) HandleUpload(c *gin.Context) {
	_ = r.GetCurrentUser(c)
	video, err := c.FormFile("video")
	if err == nil {
		var body io.ReadCloser
		body, err = video.Open()
		if err == nil {
			callbackURL := r.GenerateCallbackURL()
			err = r.api.UploadVideo(callbackURL.String(), video.Filename, body)
		}
	}
	if err == nil {
		c.HTML(http.StatusOK, "upload_processing.html", map[string]interface{}{})
	} else {
		r.Error(c, err)
	}
}
