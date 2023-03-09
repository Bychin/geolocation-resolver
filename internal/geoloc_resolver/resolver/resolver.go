package resolver

import (
	"context"
	"fmt"
	"log"
	"net"

	"geolocation-resolver/pkg/geoloc"
	"geolocation-resolver/pkg/importer"
)

// GeoDB is an interface that represents storage for reading and writing
// geolocation data.
type GeoDB interface {
	// Write writes entry to storage.
	Write(context.Context, *geoloc.Entry) error

	// Read reads entry by ip from storage. If the entry is not found, nil
	// error must be returned.
	Read(context.Context, net.IP) (*geoloc.Entry, error)
}

// Resolver can import CSV files containing the raw geolocation data, parse it
// and store it in db. It also provides access to geolocation data stored in db.
type Resolver struct {
	config *Config
	db     GeoDB
}

// NewResolver creates new Resolver from config and database for geolocation data.
func NewResolver(config *Config, db GeoDB) *Resolver {
	return &Resolver{
		config: config,
		db:     db,
	}
}

func (r *Resolver) storeRoutine(entries <-chan *geoloc.Entry, done chan<- struct{}) {
	defer close(done)

	for entry := range entries {
		ctx, cancel := context.WithTimeout(context.Background(), r.config.WriteTimeout)
		if err := r.db.Write(ctx, entry); err != nil {
			log.Printf("[E] resolver: failed to write entry to DB: %s\n", err)
		}
		cancel()
	}
}

// ImportCSV imports CSV file at path, parses it and stores geolocation data in db.
func (r *Resolver) ImportCSV(path string) error {
	if len(r.config.Importer.Separator) == 0 {
		return fmt.Errorf("invalid separator")
	}

	comma := rune(r.config.Importer.Separator[0])
	importer := importer.NewCSVImporter(path, comma, r.config.Importer.Verbose)

	done := make(chan struct{})
	go r.storeRoutine(importer.Entries(), done)

	stats, err := importer.Exec()
	if err != nil {
		return err
	}
	<-done

	if r.config.Importer.Verbose {
		log.Printf("[I] resolver: parsed %d records successfully, skipped %d, elapsed time %s\n",
			stats.ParsedLines, stats.ReadLines-stats.ParsedLines, stats.ElapsedTime)
	}
	return nil
}

// ResolveIP resolves geolocation data for ip address. If geolocation data for
// ip is not found in db, nil error is returned.
func (r *Resolver) ResolveIP(ip net.IP) (*geoloc.Entry, error) {
	if ip == nil {
		return nil, fmt.Errorf("missing ip")
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.config.ReadTimeout)
	defer cancel()

	entry, err := r.db.Read(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("failed to read entry from DB: %w", err)
	}

	return entry, nil
}
