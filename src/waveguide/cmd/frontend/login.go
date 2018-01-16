package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func (r *Routes) ServeLogin(c *gin.Context) {
	if r.oauth == nil {
		r.NotFound(c)
	} else {
		callback, _ := url.Parse(r.FrontendURL.String())
		callback.Path = "/oauth/redirect_uri"
		c.Redirect(http.StatusTemporaryRedirect, r.oauth.AuthURL(callback.String()))
	}
}

func (r *Routes) ApiLogin(c *gin.Context) {

}
