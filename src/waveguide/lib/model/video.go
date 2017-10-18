package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"waveguide/lib/api"
	"waveguide/lib/util"
)

const VideoURLBase = "/v"

type VideoInfo struct {
	UserID      int64
	VideoID     int64
	Title       string
	Description string
	UploadedAt  int64
	WebSeeds    []string
	TorrentURL  string
}

func (v *VideoInfo) GetURL(frontend *url.URL) *url.URL {
	u, _ := url.Parse(frontend.String())
	u.Path = fmt.Sprintf("%s/%d/", VideoURLBase, v.VideoID)
	return u
}

func (v *VideoInfo) VideoUploadRequest(worker, filename string, body io.ReadCloser) *http.Request {
	u, _ := url.Parse(worker)
	u.Path = "/api/" + api.EncodeVideo
	buff := new(util.Buffer)
	json.NewEncoder(buff).Encode(v)
	return &http.Request{
		URL: u,
		GetBody: func() (io.ReadCloser, error) {
			return util.NewMultipartPipe(
				[]util.MimePart{
					util.MimePart{
						Body:     body,
						PartName: api.ParamVideoFile,
					},
					util.MimePart{
						Body:     buff,
						PartName: api.ParamVideoInfoJSON,
					}}), nil
		},
	}
}

// generate an http request that does the video ready callback
func (v *VideoInfo) VideoReadyRequest(baseurl, callback *url.URL, nounce string) *http.Request {
	u, _ := url.Parse(baseurl.String())
	u.Path = "/api/" + api.RegisterVideo
	q := u.Query()
	q.Add(api.ParamNounce, nounce)
	q.Add(api.ParamCallbackURL, callback.String())
	u.RawQuery = q.Encode()
	return &http.Request{
		Method: "POST",
		URL:    u,
		GetBody: func() (io.ReadCloser, error) {
			buff := new(util.Buffer)
			err := json.NewEncoder(buff).Encode(v)
			return buff, err
		},
	}
}
