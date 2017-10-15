package model

import (
	"net/url"
	"waveguide/lib/api"
)

type VideoInfo struct {
	UserID      int64
	VideoID     int64
	Title       string
	Description string
	UploadedAt  int64
	WebSeeds    []string
	TorrentURL  string
}

func (v *VideoInfo) VideoReadyURL(baseurl *url.URL, nounce string) *url.URL {
	u, _ := url.Parse(baseurl.String())
	u.Path += api.VideoReady
	q := u.Query()
	q.Add(api.ParamNounce, nounce)
	u.RawQuery = q.Encode()
	return u
}
