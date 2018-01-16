package frontend

import (
	"github.com/gin-gonic/gin"
)

func (r *Routes) ServeLogin(c *gin.Context) {
	if r.oauth == nil {
		r.NotFound(c)
	} else {
		r.ServeOAuthLogin(c)
	}
}

func (r *Routes) ApiLogin(c *gin.Context) {

}
