package frontend

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (r *Routes) ServeStream(c *gin.Context) {
	u := r.GetCurrentUser(c)
	c.HTML(http.StatusOK, "stream.html", map[string]interface{}{
		"User": u,
	})
}

func (r *Routes) ApiStreamMagnets(c *gin.Context) {
	var magnets []string
	status := http.StatusNotFound
	key := c.GetInt64("UserID")

	stream := r.Streaming.Find(key)
	if stream != nil {
		magnets = append(magnets, stream.Magnets...)
		status = http.StatusOK
	}
	c.AbortWithStatusJSON(status, magnets)
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
