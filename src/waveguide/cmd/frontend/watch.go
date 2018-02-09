package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *Routes) ServeWatch(c *gin.Context) {
	id, _ := c.GetQuery("u")
	if id == "" {
		r.NotFound(c)
	} else {
		chatID := r.ChatIDForStream(id)
		info := r.Streaming.Find(id)
		c.HTML(http.StatusOK, "video_live.html", map[string]interface{}{
			"User":     r.GetCurrentUser(c),
			"ChatID":   chatID,
			"StreamID": id,
			"Stream":   info,
		})
	}
}
