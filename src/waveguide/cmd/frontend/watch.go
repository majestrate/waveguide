package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *Routes) ServeWatch(c *gin.Context) {
	id, _ := c.GetQuery("u")
	if id == "" {
		r.NotFound(c)
	} else {
		c.HTML(http.StatusOK, "watch.html", map[string]interface{}{
			"StreamID": id,
		})
	}
}
