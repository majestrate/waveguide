package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *Routes) ServeLogin(c *gin.Context) {
	if r.adn == nil {
		c.HTML(http.StatusOK, "login.html", map[string]interface{}{})
	} else {
		r.ServeOAuthLogin(c)
	}
}

func (r *Routes) ApiLogout(c *gin.Context) {
	r.ResetUser(c)
	c.Redirect(http.StatusFound, "/")
}

func (r *Routes) ApiLogin(c *gin.Context) {

}
