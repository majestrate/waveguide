package frontend

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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
	magnet := ""
	status := http.StatusNotFound
	key := c.GetString("UserID")

	stream := r.Streaming.Find(key)
	if stream != nil {
		status = http.StatusOK
		magnet = stream.LastMagnet()
	}
	c.String(status, magnet)
}

func (r *Routes) ApiStreamUpdate(c *gin.Context) {
	u := r.GetCurrentUser(c)
	stream := r.Streaming.Ensure(u.UserID)
	status := http.StatusOK
	buff := new(bytes.Buffer)
	io.Copy(buff, c.Request.Body)
	c.Request.Body.Close()
	stream.Add(buff.String())
	c.String(status, "")
}
