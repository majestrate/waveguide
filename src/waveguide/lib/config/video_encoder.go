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
	FFmpegPath    string
	FFmpegParams  string
	FFprobePath   string
	FFprobeParams string
}

func (c *FFMPEGConfig) Load(s *parser.Section) error {
	c.FFmpegPath = s.ValueOf("ffmpeg_path")
	c.FFmpegParams = s.ValueOf("ffmpeg_params")
	c.FFprobePath = s.ValueOf("ffprobe_path")
	c.FFprobeParams = s.ValueOf("ffprobe_params")
	return nil
}
