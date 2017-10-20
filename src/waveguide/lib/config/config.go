package config

import (
	"waveguide/lib/config/parser"
)

type Configurable interface {
	Load(s *parser.Section) error
}

type Config struct {
	Worker   WorkerConfig
	Frontend FrontendConfig
	DB       DBConfig
	MQ       MQConfig
	Storage  StorageConfig
	CDN      CDNConfig
}

func (c *Config) Load(fname string) error {
	sects := map[string]Configurable{
		"worker":   &c.Worker,
		"frontend": &c.Frontend,
		"database": &c.DB,
		"rabbitmq": &c.MQ,
		"storage":  &c.Storage,
		"cdn":      &c.CDN,
	}
	cfg, err := parser.Read(fname)
	if err != nil {
		return err
	}

	for name := range sects {
		sect, err := cfg.Section(name)
		if err == nil {
			err = sects[name].Load(sect)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
