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

func (v *VideoInfo) VideoUploadRequest(worker, callback_url, filename string, body io.ReadCloser) *http.Request {
	u, _ := url.Parse(worker)
	u.Path = "/api/" + api.EncodeVideo
	q := u.Query()
	q.Add(api.ParamFilename, filename)
	q.Add(api.ParamCallbackURL, callback_url)
	u.RawQuery = q.Encode()
	buff := new(util.Buffer)
	json.NewEncoder(buff).Encode(v)
	boundary := util.RandStr(16)
	requestBody := util.NewMultipartPipe(boundary,
		[]util.MimePart{
			util.MimePart{
				Body:     body,
				PartName: api.ParamVideoFile,
			},
			util.MimePart{
				Body:     buff,
				PartName: api.ParamVideoInfoJSON,
			},
		})
	return &http.Request{
		URL:    u,
		Method: "POST",
		Header: map[string][]string{
			"Content-Type": []string{fmt.Sprintf("multipart/mixed; boundary=%s", boundary)},
		},
		Body: requestBody,
	}
}

// generate an http request that does the video ready callback
func (v *VideoInfo) VideoReadyURL(baseurl *url.URL, nounce string) *url.URL {
	u, _ := url.Parse(baseurl.String())
	u.Path = "/callback"
	q := u.Query()
	q.Add(api.ParamNounce, nounce)
	q.Add(api.ParamVideoID, fmt.Sprintf("%d", v.VideoID))
	q.Add(api.ParamState, api.StateVideoReady)
	q.Add(api.ParamAction, api.ActionVideoState)
	u.RawQuery = q.Encode()
	return u
}
