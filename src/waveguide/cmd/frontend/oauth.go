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
	c.Redirect(http.StatusTemporaryRedirect, r.adn.AuthURL(r.oauthCallback()))
}

func (r *Routes) HandleOAuthRedirect(c *gin.Context) {
	if r.adn == nil {
		r.NotFound(c)
	} else {
		code := c.Query("code")
		user, err := r.adn.GrantUser(code, r.oauthCallback())
		if err == nil {
			r.SetCurrentUser(*user, c)
			c.Redirect(http.StatusTemporaryRedirect, "/")
		} else {
			r.Error(c, err)
		}
	}
}
