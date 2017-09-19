package video

import (
	"errors"
	"strings"
	"waveguide/lib/config"
)

var ErrNoEncoder = errors.New("no such encoder")

type Encoder interface {
	Init() error
	EncodeFile(infile, outfile string) error
}

func NewEncoder(conf *config.VideoEncoderConfig) (enc Encoder, err error) {
	switch conf.Dialect {
	case config.ExternalFFMPEG:
		enc = &FFMPEGEncoder{
			Path:   conf.FFMPEG.Path,
			Params: strings.Split(conf.FFMPEG.Params, " "),
		}
	}
	if enc == nil {
		err = ErrNoEncoder
	} else {
		err = enc.Init()
	}
	return
}
