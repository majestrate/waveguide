package frontend

import (
	"net/url"
)

func (r *Routes) getNextWorkerURL() *url.URL {
	u, _ := url.Parse(r.workerURL)
	return u
}
