package config

import (
	"waveguide/lib/config/parser"
)

type FrontendConfig struct {
	StaticDir   string
	TemplateDir string
	DB          DBConfig
}

func (c *FrontendConfig) Load(s *parser.Section) error {
	c.TemplateDir = s.ValueOf("templates")
	c.StaticDir = s.ValueOf("staticfiles")
	return c.DB.Load(s)
}
