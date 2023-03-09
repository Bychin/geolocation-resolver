package resolver

import (
	"time"
)

const (
	DefaultResolverWriteTimeout = 1 * time.Second
	DefaultResolverReadTimeout  = 1 * time.Second
	DefaultImporterSeparator    = ","
)

// Config is a config for Resolver
type Config struct {
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	Importer     struct {
		Separator string `yaml:"separator"`
		Verbose   bool   `yaml:"verbose"`
	} `yaml:"importer"`
}

// WithDefaults sets the default values for config, which is a copy of c.
func (c *Config) WithDefaults() (config Config) {
	if c != nil {
		config = *c
	}

	if config.WriteTimeout == 0 {
		config.WriteTimeout = DefaultResolverWriteTimeout
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = DefaultResolverReadTimeout
	}
	if config.Importer.Separator == "" {
		config.Importer.Separator = DefaultImporterSeparator
	}

	return
}
