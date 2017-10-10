package config

import (
	"waveguide/lib/config/parser"
)

type FrontendConfig struct {
	StaticDir   string
	TemplateDir string
	DB          DBConfig
	WorkerURL   string
}

func (c *FrontendConfig) Load(s *parser.Section) error {
	c.TemplateDir = s.ValueOf("templates")
	c.StaticDir = s.ValueOf("staticfiles")
	c.WorkerURL = s.ValueOf("worker")
	return c.DB.Load(s)
}
