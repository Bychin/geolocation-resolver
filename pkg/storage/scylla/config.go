package scylla

import (
	"time"
)

const (
	DefaultServer         = "127.0.0.1:9042"
	DefaultKeyspace       = "geoloc"
	DefaultTimeout        = 1 * time.Second
	DefaultConnectTimeout = 1 * time.Second
	DefaultNumConns       = 10
	DefaultWriterPoolSize = 10
	DefaultWriteQueueSize = 10
	DefaultReaderPoolSize = 10
	DefaultReadQueueSize  = 10
)

// Config is a config for Connector for ScyllaDB
type Config struct {
	Servers        []string      `yaml:"servers"`
	Keyspace       string        `yaml:"keyspace"`
	Timeout        time.Duration `yaml:"timeout"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	NumConns       int           `yaml:"num_conns"`
	WriterPoolSize int           `yaml:"writer_pool_size"`
	WriteQueueSize int           `yaml:"write_queue_size"`
	ReaderPoolSize int           `yaml:"reader_pool_size"`
	ReadQueueSize  int           `yaml:"read_queue_size"`
}

// WithDefaults sets the default values for config, which is a copy of c.
func (c *Config) WithDefaults() (config Config) {
	if c != nil {
		config = *c
	}

	if len(config.Servers) == 0 {
		config.Servers = []string{DefaultServer}
	}
	if config.Keyspace == "" {
		config.Keyspace = DefaultKeyspace
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = DefaultConnectTimeout
	}
	if config.NumConns == 0 {
		config.NumConns = DefaultNumConns
	}
	if config.WriterPoolSize == 0 {
		config.WriterPoolSize = DefaultWriterPoolSize
	}
	if config.WriteQueueSize == 0 {
		config.WriteQueueSize = DefaultWriteQueueSize
	}
	if config.ReaderPoolSize == 0 {
		config.ReaderPoolSize = DefaultReaderPoolSize
	}
	if config.ReadQueueSize == 0 {
		config.ReadQueueSize = DefaultReadQueueSize
	}

	return
}
