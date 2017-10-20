package database

import (
	"waveguide/lib/model"
)

type Database interface {
	Init() error
	GetFrontpageVideos() (model.VideoList, error)
	RegisterVideo(v *model.VideoInfo) error
	GetVideoInfo(id int64) (*model.VideoInfo, error)
	NextVideoID() (int64, error)
	Close() error
}

func NewDatabase(url string) Database {
	return &pgDB{
		url: url,
	}
}
