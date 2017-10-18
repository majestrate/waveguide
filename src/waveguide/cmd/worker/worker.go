package worker

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"waveguide/lib/api"
	"waveguide/lib/database"
	"waveguide/lib/log"
	"waveguide/lib/torrent"
	"waveguide/lib/video"
)

type Worker struct {
	WorkerURL string
	Encoder   video.Encoder
	Torrent   torrent.Factory
	TempDir   string
	DB        database.Database
}

// get next worker in worker pool
func (w *Worker) GetNextWorkerURL() *url.URL {
	u, _ := url.Parse(w.WorkerURL)
	return u
}

func (w *Worker) APIAccepted(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"Error": nil,
	})
}

func (w *Worker) APIError(c *gin.Context, err error) {
	log.Errorf("worker error: %s", err.Error())
	c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"Error": err.Error(),
	})
}

func (w *Worker) ServeAPI(c *gin.Context) {
	callbackUrl := c.Query(api.ParamCallbackURL)
	u, err := url.Parse(callbackUrl)
	if err != nil {
		w.APIError(c, err)
		return
	}

	method := c.Param("Method")
	switch method {
	case api.EncodeVideo:
		err = w.VideoEncode(c, u)
	case api.MakeTorrent:
		err = w.MakeTorrent(c, u)
	default:
		err = api.ErrNoSuchMethod
	}
	if err == nil {
		w.APIAccepted(c)
	} else {
		w.APIError(c, err)
	}
}
