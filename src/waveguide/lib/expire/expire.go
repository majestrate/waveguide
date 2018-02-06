package expire

import (
	"waveguide/lib/model"
)

const DefaultCapacity = 50

type ExpirePolicy interface {
	GetExpiredVideos(capacity uint64) ([]model.VideoInfo, error)
}
