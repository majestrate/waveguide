package config

import (
	"os"
	"waveguide/lib/config/parser"
)

type WorkerConfig struct {
	TempDir   string
	UploadURL string
	Encoder   VideoEncoderConfig
	Torrent   TorrentConfig
}

func (c *WorkerConfig) Load(s *parser.Section) error {
	c.TempDir = s.Get("tempfiles", os.TempDir())
	c.UploadURL = s.ValueOf("cdn_upload_url")
	err := c.Torrent.Load(s)
	if err != nil {
		return err
	}
	return c.Encoder.Load(s)
}
