package expire

import (
	"waveguide/lib/model"
)

const DefaultCapacity = 10

type ExpirePolicy interface {
	GetExpiredVideos(capacity uint64) ([]model.VideoInfo, error)
}
