package database

import (
	"waveguide/lib/model"
)

type Database interface {
	Init() error
	GetFrontpageVideos() (model.VideoList, error)
	RegisterVideo(v *model.VideoInfo) error
	SetVideoMetaInfo(id int64, url string) error
	AddWebseed(id int64, url string) error
	GetVideoInfo(id int64) (*model.VideoInfo, error)
	Close() error
}

func NewDatabase(url string) Database {
	return &pgDB{
		url: url,
	}
}
