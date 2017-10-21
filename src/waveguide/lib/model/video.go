package model

import (
	"fmt"
	"net/url"
	"waveguide/lib/api"
	"waveguide/lib/log"
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
	u, err := url.Parse(frontend.String())
	if err != nil {
		log.Errorf("waveguide/lib/model/Video.GetURL: %s", err)
	}
	u.Path = fmt.Sprintf("%s/%d/", VideoURLBase, v.VideoID)
	return u
}

func (v *VideoInfo) WebseedUploadRequest(remoteFile *url.URL) *api.Request {
	return &api.Request{
		Method: api.MakeTorrent,
		Args: map[string]interface{}{
			api.ParamVideoID:  v.VideoID,
			api.ParamFilename: v.Title,
			api.ParamFileURL:  remoteFile.String(),
		},
	}
}

func (v *VideoInfo) VideoUploadRequest(fileURL *url.URL, filename string) *api.Request {
	return &api.Request{
		Method: api.EncodeVideo,
		Args: map[string]interface{}{
			api.ParamVideoID:  v.VideoID,
			api.ParamFilename: filename,
			api.ParamFileURL:  fileURL.String(),
		},
	}
}
