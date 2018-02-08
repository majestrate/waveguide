package streaming

import (
	"time"
)

const StreamUpdateTimeout = time.Minute

type StreamInfo struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	URLS       [3]string `json:"urls"`
	Token      string    `json:"-"`
	LastUpdate time.Time `json:"-"`
}

func (i *StreamInfo) LastTorrent() string {
	return i.URLS[0]
}

func (i *StreamInfo) OldestTorrent() string {
	return i.URLS[2]
}

func (i *StreamInfo) AddTorrent(url string) {
	i.URLS[2] = i.URLS[1]
	i.URLS[1] = i.URLS[0]
	i.URLS[0] = url
	i.LastUpdate = time.Now()
}

func (i *StreamInfo) IsExpired() bool {
	now := time.Now()
	return now.Sub(i.LastUpdate) > StreamUpdateTimeout
}
