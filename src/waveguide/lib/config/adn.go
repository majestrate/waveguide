package config

import (
	"strconv"
	"waveguide/lib/config/parser"
)

type ADNConfig struct {
	Provider     string
	ClientID     string
	ClientSecret string
	Workers      int
	Enabled      bool
}

func (c *ADNConfig) Load(s *parser.Section) (err error) {
	c.Provider = s.ValueOf("provider")
	c.ClientID = s.ValueOf("clientid")
	c.ClientSecret = s.ValueOf("clientsecret")
	c.Enabled = s.ValueOf("enabled") == "1"
	c.Workers, err = strconv.Atoi(s.Get("workers", "8"))
	return
}
