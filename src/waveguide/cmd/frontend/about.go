package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *Routes) ServeAbout(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}
