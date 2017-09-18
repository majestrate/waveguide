package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"waveguide/lib/database"
)

type Routes struct {
	DB database.Database
}

func (r *Routes) Error(c *gin.Context, err error) {
	c.HTML(http.StatusInternalServerError, "errors.html", map[string]string{
		"Error": err.Error(),
	})
}

func (r *Routes) Index(c *gin.Context) {
	videos, err := r.DB.GetFrontpageVideos()
	if err != nil {
		r.Error(c, err)
		return
	}
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"Videos": videos,
	})
}
