package scylla

import (
	"context"
	"errors"
	"fmt"
	"net"

	"geolocation-resolver/pkg/geoloc"
	"geolocation-resolver/pkg/storage/scylla/pool"

	"github.com/gocql/gocql"
)

// Connector is a connector for ScyllaDB with functions for reading and writing
// geolocation data.
type Connector struct {
	session *gocql.Session
	config  *Config

	writers    []*pool.Writer
	writeQueue pool.WriteRequestsQueue

	readers   []*pool.Reader
	readQueue pool.ReadRequestsQueue
}

// NewConnector creates Connector from config. By doing this, it also
// establishes a connection to ScyllaDB.
func NewConnector(config *Config) (*Connector, error) {
	cluster := gocql.NewCluster(config.Servers...)
	cluster.Consistency = gocql.One
	cluster.Timeout = config.Timeout
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.NumConns = config.NumConns
	cluster.ConnectTimeout = config.ConnectTimeout
	cluster.Keyspace = config.Keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &Connector{
		session: session,
		config:  config,

		writers:    make([]*pool.Writer, config.WriterPoolSize),
		writeQueue: pool.NewWriteRequestsQueue(config.WriteQueueSize),
		readers:    make([]*pool.Reader, config.ReaderPoolSize),
		readQueue:  pool.NewReadRequestsQueue(config.ReadQueueSize),
	}, nil
}

// Init initializes internal read and write workers.
func (c *Connector) Init() {
	for i := 0; i < len(c.writers); i++ {
		c.writers[i] = pool.NewWriter(c.session, c.writeQueue)
		c.writers[i].Init()
	}

	for i := 0; i < len(c.readers); i++ {
		c.readers[i] = pool.NewReader(c.session, c.readQueue)
		c.readers[i].Init()
	}
}

// Shutdown synchronously shutdowns internal read and write workers and closes
// the connection to ScyllaDB.
func (c *Connector) Shutdown() {
	close(c.writeQueue)
	close(c.readQueue)

	for i := 0; i < len(c.writers); i++ {
		<-c.writers[i].Done()
	}
	for i := 0; i < len(c.readers); i++ {
		<-c.readers[i].Done()
	}

	c.session.Close()
}

// Write synchronously writes entry to ScyllaDB.
func (c *Connector) Write(ctx context.Context, entry *geoloc.Entry) error {
	req := pool.NewWriteGeolocEntryRequest(entry)
	select {
	case c.writeQueue <- req:
		// pass
	case <-ctx.Done():
		return fmt.Errorf("failed to push req into queue: %w", ctx.Err())
	}

	select {
	case resp := <-req.Response():
		return resp.Err
	case <-ctx.Done():
		return fmt.Errorf("failed to get response: %w", ctx.Err())
	}
}

// Read synchronously reads entry by ip from ScyllaDB. If the entry is not
// found, nil error is returned.
func (c *Connector) Read(ctx context.Context, ip net.IP) (*geoloc.Entry, error) {
	req := pool.NewReadGeolocEntryRequest(ip)
	select {
	case c.readQueue <- req:
		// pass
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to push req into queue: %w", ctx.Err())
	}

	select {
	case resp := <-req.Response():
		if errors.Is(resp.Err, gocql.ErrNotFound) {
			return nil, nil
		}
		return resp.Entry, resp.Err
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to get response: %w", ctx.Err())
	}
}
