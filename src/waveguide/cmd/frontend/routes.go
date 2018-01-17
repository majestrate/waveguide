package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"waveguide/lib/api"
	"waveguide/lib/config"
	"waveguide/lib/database"
	"waveguide/lib/model"
	"waveguide/lib/oauth"
	"waveguide/lib/streaming"
)

type Routes struct {
	DB          database.Database
	api         *api.Client
	FrontendURL *url.URL
	TempDir     string
	Streaming   *streaming.Context
	oauth       *oauth.Client
}

func (r *Routes) Close() error {
	r.DB.Close()
	r.api.Close()
	return nil
}

func (r *Routes) Configure(c *config.Config) error {
	r.Streaming = streaming.NewContext()
	return r.configure(c, false)
}

func (r *Routes) Reconfigure(c *config.Config) error {
	return r.configure(c, true)
}

func (r *Routes) configure(c *config.Config, reload bool) (err error) {
	if c.OAuth.Enabled {
		r.oauth = oauth.NewClient(c.OAuth)
	} else {
		r.oauth = nil
	}
	r.TempDir = c.Storage.TempDir

	r.FrontendURL, err = url.Parse(c.Frontend.FrontendURL)
	if err != nil {
		return
	}

	if r.DB != nil {
		r.DB.Close()
	}

	r.DB = database.NewDatabase(c.DB.URL)

	err = r.DB.Init()
	if err != nil {
		return
	}

	if r.api != nil {
		r.api.Close()
	}

	r.api, err = api.NewClient(&c.MQ)
	return
}

func (r *Routes) Error(c *gin.Context, err error) {
	c.HTML(http.StatusInternalServerError, "error.html", map[string]string{
		"Error": err.Error(),
	})
}

func (r *Routes) NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error_404.html", nil)
}

func (r *Routes) ServeIndex(c *gin.Context) {
	videos, err := r.DB.GetFrontpageVideos()
	if err != nil {
		r.Error(c, err)
		return
	}
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"Videos":  videos.Videos,
		"Streams": r.Streaming.Online(),
	})
}

func (r *Routes) ServeVideo(c *gin.Context) {
	videoID := c.Param("id")
	info, err := r.DB.GetVideoInfo(videoID)
	if err == nil && info != nil {
		c.HTML(http.StatusOK, "video.html", map[string]*model.VideoInfo{
			"Video": info,
		})
		return
	}
	r.NotFound(c)
}

func (r *Routes) ServeUpload(c *gin.Context) {
	u := r.GetCurrentUser(c)
	c.HTML(http.StatusOK, "upload.html", map[string]interface{}{
		"User": u,
	})
}
