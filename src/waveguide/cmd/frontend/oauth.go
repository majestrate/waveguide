package frontend

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"net/url"
)

func (r *Routes) oauthCallback() string {
	callback, _ := url.Parse(r.FrontendURL.String())
	callback.Path = "/oauth/redirect_uri"
	return callback.String()
}

func (r *Routes) ServeOAuthLogin(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, r.oauth.AuthURL(r.oauthCallback()))
}

func (r *Routes) HandleOAuthRedirect(c *gin.Context) {
	if r.oauth == nil {
		r.NotFound(c)
	} else {
		code := c.Query("code")
		user, err := r.oauth.GetUser(code, r.oauthCallback())
		if err == nil {
			r.SetCurrentUser(*user, c)
			c.Redirect(http.StatusTemporaryRedirect, "/stream/")
		} else {
			r.Error(c, err)
		}
	}
}
