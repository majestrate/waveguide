package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"waveguide/lib/adn"
)

func (r *Routes) ServeWatch(c *gin.Context) {
	id, _ := c.GetQuery("u")
	if id == "" {
		r.NotFound(c)
	} else {
		info := r.Streaming.Find(id)
		var chatID adn.ChanID
		if info == nil {
			chatID = adn.ChanID(5)
		} else {
			chatID = info.ChatID
		}
		c.HTML(http.StatusOK, "video_live.html", map[string]interface{}{
			"User":     r.GetCurrentUser(c),
			"ChatID":   chatID,
			"StreamID": id,
			"Stream":   info,
		})
	}
}
