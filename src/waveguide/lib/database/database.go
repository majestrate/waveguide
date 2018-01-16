package database

import (
	"waveguide/lib/model"
)

type Database interface {
	Init() error
	GetFrontpageVideos() (model.VideoList, error)
	RegisterVideo(v *model.VideoInfo) error
	SetVideoMetaInfo(id string, url string) error
	AddWebseed(id string, url string) error
	GetVideoInfo(id string) (*model.VideoInfo, error)
	//GetUserByName(name string) (*model.UserInfo, error)
	//GetUserByID(id int64) (*model.UserInfo, error)
	GetVideosForUserByName(name string) (*model.VideoFeed, error)
	Close() error
}

func NewDatabase(url string) Database {
	return &pgDB{
		url: url,
	}
}
