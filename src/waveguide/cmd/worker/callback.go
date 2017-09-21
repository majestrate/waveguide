package worker

import (
	"net/http"
	"net/url"
	"waveguide/lib/log"
)

func (w *Worker) InformCallback(callback *url.URL, err error) {
	if err != nil {
		q := callback.Query()
		q.Add("error", err.Error())
		callback.RawQuery = q.Encode()
	}
	log.Debugf("inform callback %s", callback.String())
	_, err = http.Get(callback.String())
	if err != nil {
		log.Errorf("error informing callback: %s", err.Error())
	}
}
