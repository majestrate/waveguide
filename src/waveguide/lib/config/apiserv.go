package config

import (
	"waveguide/lib/config/parser"
)

type ApiServerConfig struct {
	Addr    string
	Enabled bool
	Anon    bool
}

func (c *ApiServerConfig) Load(s *parser.Section) error {
	c.Addr = s.ValueOf("bind")
	c.Enabled = s.ValueOf("enabled") == "1"
	c.Anon = s.ValueOf("anon") == "1"
	return nil
}
