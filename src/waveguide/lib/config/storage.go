package config

import (
	"waveguide/lib/config/parser"
)

type StorageConfig struct {
	TempDir string
}

func (c *StorageConfig) Load(s *parser.Section) error {
	c.TempDir = s.ValueOf("tempdir")
	return nil
}
