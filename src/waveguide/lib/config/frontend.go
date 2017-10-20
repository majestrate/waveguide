package config

import (
	"waveguide/lib/config/parser"
)

type FrontendConfig struct {
	StaticDir   string
	TemplateDir string
	WorkerURL   string
	FrontendURL string
	Addr        string
}

func (c *FrontendConfig) Load(s *parser.Section) error {
	c.TemplateDir = s.ValueOf("templates")
	c.StaticDir = s.ValueOf("staticfiles")
	c.WorkerURL = s.ValueOf("worker")
	c.FrontendURL = s.ValueOf("url")
	c.Addr = s.ValueOf("addr")
	return nil
}
