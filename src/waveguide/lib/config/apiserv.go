package config

import (
	"waveguide/lib/config/parser"
)

type ApiServerConfig struct {
	Addr    string
	Enabled bool
}

func (c *ApiServerConfig) Load(s *parser.Section) error {
	c.Addr = s.ValueOf("bind")
	c.Enabled = s.ValueOf("enabled") == "1"
	return nil
}
