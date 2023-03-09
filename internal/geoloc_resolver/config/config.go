package config

import (
	"fmt"
	"log"
	"os"

	"geolocation-resolver/internal/geoloc_resolver/resolver"
	"geolocation-resolver/pkg/storage/scylla"

	"gopkg.in/yaml.v3"
)

const (
	StorageOptionHashmap = "hashmap"
	StorageOptionScylla  = "scylla"

	DefaultStorageOption = StorageOptionHashmap
	DefaultHTTPAddr      = "localhost:8080"
)

// Config is a config for geoloc_resolver app.
type Config struct {
	PathToCSV     string `yaml:"path_to_csv"`
	StorageOption string `yaml:"storage_option"`
	HTTP          struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`
	Resolver resolver.Config `yaml:"resolver"`
	Scylla   scylla.Config   `yaml:"scylla"`
}

// NewConfig parses yaml config from file at path.
func NewConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open yaml file: %w", err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	config := &Config{}
	if err := d.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode yaml content: %w", err)
	}

	return config, nil
}

// WithDefaults sets the default values for config, which is a copy of c.
func (c *Config) WithDefaults() (config Config) {
	if c != nil {
		config = *c
	}

	if config.StorageOption == "" {
		config.StorageOption = DefaultStorageOption
	}
	if config.HTTP.Addr == "" {
		config.HTTP.Addr = DefaultHTTPAddr
	}

	config.Resolver = config.Resolver.WithDefaults()
	config.Scylla = config.Scylla.WithDefaults()

	return
}

// Print prints config in yaml format.
func (c *Config) Print() {
	out, err := yaml.Marshal(c)
	if err != nil {
		log.Printf("[E] config: failed to marshal config: %s\n", err)
	}

	log.Printf("[I] config: using the following args for geoloc_resolver:\n%s\n", string(out))
}
