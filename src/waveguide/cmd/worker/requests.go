package worker

import (
	"fmt"
	"net/http"
	"waveguide/lib/log"
)

func (w *Worker) DoRequest(r *http.Request) error {
	log.Debugf("do request %s %s", r.Method, r.URL.String())
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Errorf("error doing request: %s", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return err
}
