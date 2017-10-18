package worker

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"waveguide/lib/api"
	"waveguide/lib/log"
)

func (w *Worker) InformCallback(callback *url.URL, err error) {
	if err != nil {
		q := callback.Query()
		q.Add(api.ParamError, err.Error())
		callback.RawQuery = q.Encode()
	}
	log.Debugf("inform callback %s", callback.String())
	_, err = http.Get(callback.String())
	if err != nil {
		log.Errorf("error informing callback: %s", err.Error())
	}
}

func (w *Worker) ServeCallback(context *gin.Context) {
}
