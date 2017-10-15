package frontend

import (
	"net/url"
	"waveguide/lib/model"
)

func (r *Routes) VideoDoneCallbackURL(v *model.VideoInfo) *url.URL {
	// TODO: session nounce
	nounce := ""
	u, _ := url.Parse(r.workerURL)
	return v.VideoReadyURL(u, nounce)
}
