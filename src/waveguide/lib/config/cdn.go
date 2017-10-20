package config

import (
	"strings"
	"waveguide/lib/config/parser"
)

type CDNConfig struct {
	WebseedServers  []string
	MetainfoServers []string
}

func (c *CDNConfig) Load(s *parser.Section) error {
	c.WebseedServers = strings.Split(s.ValueOf("webseed_servers"), ",")
	c.MetainfoServers = strings.Split(s.ValueOf("metainfo_servers"), ",")
	return nil
}
