package model

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"
	"waveguide/lib/log"
	"waveguide/lib/util"
	"waveguide/lib/worker/api"
)

const VideoURLBase = "/v"

type VideoInfo struct {
	UserID      string
	VideoID     string
	Title       string
	Description string
	UploadedAt  int64
	WebSeeds    []string
	TorrentURL  string
	BaseURL     *url.URL
}

type videoInfoFeed struct {
	Title      string    `xml:"title"`
	Link       Link      `xml:"link"`
	ID         string    `xml:"id"`
	Updated    time.Time `xml:"updated"`
	Summary    string    `xml:"summary"`
	AuthorName string    `xml:"author>name"`
}

func (v VideoInfo) toFeed() *videoInfoFeed {
	u := v.GetURL(v.BaseURL)
	return &videoInfoFeed{
		Title:   v.Title,
		Link:    NewLink(u.Host, u.RequestURI()),
		ID:      util.SHA256(v.TorrentURL),
		Updated: v.CreatedAt(),
		Summary: v.Description,
	}
}

// for atom
func (v VideoInfo) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	err = e.EncodeElement(v.toFeed(), start)
	return
}

func (v VideoInfo) CreatedAt() time.Time {
	return time.Unix(v.UploadedAt, 0)
}

func (v *VideoInfo) GetURL(frontend *url.URL) *url.URL {
	u, err := url.Parse(frontend.String())
	if err != nil {
		log.Errorf("waveguide/lib/model/Video.GetURL: %s", err)
	}
	u.Path = fmt.Sprintf("%s/%s/", VideoURLBase, v.VideoID)
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
