package model

type VideoInfo struct {
	UserID      int64
	VideoID     int64
	Title       string
	Description string
	UploadedAt  int64
	WebSeeds    []string
	TorrentURL  string
}
