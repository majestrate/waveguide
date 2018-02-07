package config

import (
	"waveguide/lib/config/parser"
)

type RTMPConfig struct {
	BaseURL  string
	VideoDir string
	Enabled  bool
}

func (c *RTMPConfig) Load(s *parser.Section) error {
	c.Enabled = s.ValueOf("enabled") == "1"
	if c.Enabled {
		c.BaseURL = s.ValueOf("url")
		c.VideoDir = s.ValueOf("video_dir")
	}
	return nil
}
