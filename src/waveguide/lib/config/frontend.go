package config

import (
	"waveguide/lib/config/parser"
)

type FrontendConfig struct {
	TemplateDir string
	DB          DBConfig
}

func (c *FrontendConfig) Load(s *parser.Section) error {
	c.TemplateDir = s.ValueOf("templates")
	return c.DB.Load(s)
}
