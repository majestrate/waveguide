package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"waveguide/lib/database"
	"waveguide/lib/model"
)

type Routes struct {
	DB database.Database
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

}
