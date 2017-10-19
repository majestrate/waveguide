package config

import (
	"os"
	"waveguide/lib/config/parser"
)

type WorkerConfig struct {
	TempDir   string
	UploadURL string
	Encoder   VideoEncoderConfig
}

func (c *WorkerConfig) Load(s *parser.Section) error {
	c.TempDir = s.Get("tempfiles", os.TempDir())
	c.UploadURL = s.ValueOf("cdn_upload_url")
	return c.Encoder.Load(s)
}
