package streaming

import (
	"time"
	"waveguide/lib/adn"
)

const StreamUpdateTimeout = time.Minute
const PlaceholderThumbnail = "https://beam.sapphire.moe/assets/gitgudtv/GitGud.tv.png"

type StreamInfo struct {
	ID         adn.UID    `json:"id"`
	Username   string     `json:"username"`
	URLS       [5]string  `json:"urls"`
	Token      string     `json:"-"`
	LastUpdate time.Time  `json:"-"`
	Segments   uint64     `json:"segments"`
	ChatID     adn.ChanID `json:"chat_id"`
}

func (i *StreamInfo) ThumbnailURL() string {
	last := i.LastTorrent()
	if len(last) > 0 {
		return last[:len(last)-7] + "jpeg"
	}
	return PlaceholderThumbnail
}

func (i *StreamInfo) LastTorrent() string {
	return i.URLS[0]
}

func (i *StreamInfo) OldestTorrent() string {
	return i.URLS[4]
}

func (i *StreamInfo) AddTorrent(url string) {
	i.URLS[4] = i.URLS[3]
	i.URLS[3] = i.URLS[2]
	i.URLS[2] = i.URLS[1]
	i.URLS[1] = i.URLS[0]
	i.URLS[0] = url
	i.LastUpdate = time.Now()
}

func (i *StreamInfo) IsExpired() bool {
	now := time.Now()
	return now.Sub(i.LastUpdate) > StreamUpdateTimeout
}
