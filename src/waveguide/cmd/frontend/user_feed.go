package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func (r *Routes) ServeUserVideosFeed(c *gin.Context) {
	feed, err := r.DB.GetVideosForUserByName(c.Param("Username"))
	if err == nil {
		if feed != nil {
			feed.URL, err = url.Parse(r.FrontendURL.String())
			feed.URL.Path = c.Request.RequestURI
			if err == nil {
				c.XML(http.StatusOK, feed)
				return
			}
		}
	}
	c.Status(http.StatusNotFound)
}
