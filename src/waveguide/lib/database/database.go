package database

import (
	"waveguide/lib/model"
)

type Database interface {
	Init() error
	GetFrontpageVideos() (model.VideoList, error)
}

func NewDatabase(url string) Database {
	return &pgDB{
		url: url,
	}
}
