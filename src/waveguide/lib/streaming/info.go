package streaming

import (
	"time"
)

const StreamUpdateTimeout = time.Minute

type StreamInfo struct {
	Magnet     string
	LastUpdate time.Time
}

func (i *StreamInfo) LastMagnet() string {
	return i.Magnet
}

func (i *StreamInfo) Add(url string) {
	i.Magnet = url
	i.LastUpdate = time.Now()
}

func (i *StreamInfo) IsExpired() bool {
	now := time.Now()
	return now.Sub(i.LastUpdate) > StreamUpdateTimeout
}
