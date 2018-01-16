package config

import (
	"waveguide/lib/config/parser"
)

type OAuthConfig struct {
	Provider     string
	ClientID     string
	ClientSecret string
	Enabled      bool
}

func (c *OAuthConfig) Load(s *parser.Section) error {
	c.Provider = s.ValueOf("provider")
	c.ClientID = s.ValueOf("clientid")
	c.ClientSecret = s.ValueOf("clientsecret")
	c.Enabled = s.ValueOf("enabled") == "1"
	return nil
}
