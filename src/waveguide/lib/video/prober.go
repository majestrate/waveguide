package video

import (
	"errors"
	"waveguide/lib/config"
)

var ErrNoProber = errors.New("no such prober dialect")

type Prober interface {
	Init() error
	/** return true if the file at fname path needs to be encoded */
	VideoNeedsEncoding(fname string, wanted Info) (bool, error)
}

func NewProber(conf *config.VideoEncoderConfig) (p Prober, err error) {
	switch conf.Dialect {
	case config.ExternalFFMPEG:
		p = &FFProbe{
			Path: conf.FFMPEG.FFprobePath,
		}
	}
	if p == nil {
		err = ErrNoProber
	} else {
		err = p.Init()
	}
	return
}
