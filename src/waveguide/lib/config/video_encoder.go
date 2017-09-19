package config

import (
	"waveguide/lib/config/parser"
)

const ExternalFFMPEG = "ffmpeg"

type VideoEncoderConfig struct {
	Dialect string
	FFMPEG  FFMPEGConfig
}

func (c *VideoEncoderConfig) Load(s *parser.Section) error {
	c.Dialect = s.ValueOf("encoder")
	if c.Dialect == ExternalFFMPEG {
		return c.FFMPEG.Load(s)
	}
	return nil
}

type FFMPEGConfig struct {
	Path   string
	Params string
}

func (c *FFMPEGConfig) Load(s *parser.Section) error {
	c.Path = s.ValueOf("ffmpeg_path")
	c.Params = s.ValueOf("ffmpeg_params")
	return nil
}
