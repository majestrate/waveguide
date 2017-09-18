package config

import (
	"waveguide/lib/config/parser"
)

type WorkerConfig struct {
}

func (c *WorkerConfig) Load(s *parser.Section) error {
	return nil
}
