package model

import (
	"fmt"
	"net/url"
	"waveguide/lib/api"
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

func (v *VideoInfo) VideoUploadRequest(fileURL *url.URL, filename string) *api.Request {
	return &api.Request{
		Method: api.EncodeVideo,
		Args: map[string]interface{}{
			api.ParamVideoInfoJSON: v,
			api.ParamFilename:      filename,
			api.ParamFileURL:       fileURL.String(),
		},
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
