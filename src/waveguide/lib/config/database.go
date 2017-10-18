package config

import (
	"waveguide/lib/config/parser"
)

type DBConfig struct {
	URL string
}

func (c *DBConfig) Load(s *parser.Section) error {
	c.URL = s.ValueOf("url")
	return nil
}
