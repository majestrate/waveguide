package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"waveguide/lib/api"
)

func (r *Routes) ServeStream(c *gin.Context) {
	if r.CurrentUserLoggedIn(c) {
		u := r.GetCurrentUser(c)
		c.HTML(http.StatusOK, "stream.html", map[string]interface{}{
			"User": u,
		})
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/login/")
	}
}

func (r *Routes) ApiStreamsOnline(c *gin.Context) {
	c.JSON(http.StatusOK, r.Streaming.Online())
}

func (r *Routes) ApiStreamMagnets(c *gin.Context) {
	status := http.StatusNotFound
	key, ok := c.GetQuery("u")
	if ok {
		stream := r.Streaming.Find(key)
		if stream != nil {
			c.String(http.StatusOK, stream.LastTorrent())
			return
		}
	}
	c.String(status, "")
}

func (r *Routes) ApiStreamUpdate(c *gin.Context) {
	status := http.StatusOK
	c.String(status, "")
}

func (r *Routes) ApiStreamURL(c *gin.Context) {
	u := r.GetCurrentUser(c)
	// TODO: use app token not user token
	if r.CurrentUserLoggedIn(c) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"streamkey": u.UserID + api.StreamKeyDelim + u.Token,
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"streamkey": nil,
		})
	}
}
