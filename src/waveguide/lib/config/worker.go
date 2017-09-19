package config

import (
	"os"
	"waveguide/lib/config/parser"
)

type WorkerConfig struct {
	TempDir string
	Encoder VideoEncoderConfig
}

func (c *WorkerConfig) Load(s *parser.Section) error {
	c.TempDir = s.Get("tempfiles", os.TempDir())
	return c.Encoder.Load(s)
}
