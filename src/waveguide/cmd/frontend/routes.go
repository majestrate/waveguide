package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"waveguide/lib/config"
	"waveguide/lib/database"
	"waveguide/lib/oauth"
	pomf "waveguide/lib/pomf/api"
	"waveguide/lib/streaming"
	"waveguide/lib/worker/api"
)

type Routes struct {
	DB          database.Database
	api         *api.Client
	FrontendURL *url.URL
	TempDir     string
	Streaming   *streaming.Client
	oauth       *oauth.Client
	Pomf        *pomf.Server
}

func (r *Routes) Close() error {
	r.DB.Close()
	r.api.Close()
	return nil
}

func (r *Routes) Configure(c *config.Config) error {
	return r.configure(c, false)
}

func (r *Routes) Reconfigure(c *config.Config) error {
	return r.configure(c, true)
}

func (r *Routes) configure(c *config.Config, reload bool) (err error) {
	r.Streaming = streaming.NewClient(c)
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

	r.Pomf = pomf.NewServer(r)

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
	chatID := r.ChatIDForVideo(videoID)
	if err == nil && info != nil {
		c.HTML(http.StatusOK, "video.html", map[string]interface{}{
			"Video":  info,
			"User":   r.GetCurrentUser(c),
			"ChatID": chatID,
		})
		return
	}
	r.NotFound(c)
}

func (r *Routes) ServeUpload(c *gin.Context) {
	if r.CurrentUserLoggedIn(c) {
		u := r.GetCurrentUser(c)
		c.HTML(http.StatusOK, "upload.html", map[string]interface{}{
			"User": u,
		})
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/login/")
	}
}
