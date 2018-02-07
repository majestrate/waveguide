package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"waveguide/lib/oauth"
)

func (r *Routes) ServeLogin(c *gin.Context) {
	if r.oauth == nil {
		c.HTML(http.StatusOK, "login.html", map[string]interface{}{})
	} else {
		r.ServeOAuthLogin(c)
	}
}

func (r *Routes) ApiLogout(c *gin.Context) {
	r.SetCurrentUser(oauth.User{}, c)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (r *Routes) ApiLogin(c *gin.Context) {

}
