package frontend

import (
	"net/url"
)

func (r *Routes) GenerateCallbackURL() *url.URL {
	u, _ := url.Parse(r.workerURL)
	u.Path = "/callback"
	return u
}
