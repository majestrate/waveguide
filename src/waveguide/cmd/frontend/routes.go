package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"waveguide/lib/api"
	"waveguide/lib/database"
	"waveguide/lib/model"
)

type Routes struct {
	DB          database.Database
	api         *api.Client
	workerURL   string
	FrontendURL *url.URL
}

func (r *Routes) Error(c *gin.Context, err error) {
	c.HTML(http.StatusInternalServerError, "error.html", map[string]string{
		"Error": err.Error(),
	})
}

func (r *Routes) NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error_404.html", nil)
}

func (r *Routes) ServeIndex(c *gin.Context) {
	videos, err := r.DB.GetFrontpageVideos()
	if err != nil {
		r.Error(c, err)
		return
	}
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"Videos": videos,
	})
}

func (r *Routes) ServeVideo(c *gin.Context) {
	var info *model.VideoInfo
	videoID := c.Param("VideoID")
	id, err := strconv.ParseInt(videoID, 10, 64)
	if err == nil {
		info, err = r.DB.GetVideoInfo(id)
		if err == nil {
			c.HTML(http.StatusOK, "video.html", map[string]*model.VideoInfo{
				"Video": info,
			})
			return
		}
	}
	r.NotFound(c)
}

func (r *Routes) ServeUser(c *gin.Context) {
	// TODO: implement
}

func (r *Routes) ServeUpload(c *gin.Context) {
	u := r.GetCurrentUser(c)
	c.HTML(http.StatusOK, "upload.html", map[string]interface{}{
		"User": u,
	})
}
