package config

import (
	"waveguide/lib/config/parser"
)

type MQConfig struct {
	URL string
}

func (c *MQConfig) Load(s *parser.Section) error {
	c.URL = s.ValueOf("url")
	return nil
}
