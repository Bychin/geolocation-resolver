package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"geolocation-resolver/pkg/geoloc"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	entriesChanCap      = 10000
	uniqueIPsInitialCap = 1000
)

// CSVImporter imports data one time from a CSV file located at path, using
// comma delimiter. Import results are sent to the internal channel.
// One CSVImporter instance can only import one file at most once.
type CSVImporter struct {
	path  string
	comma rune

	entriesChan chan *geoloc.Entry
	once        *sync.Once

	stats   *Stats
	verbose bool
}

// NewCSVImporter creates CSVImporter for file located at path, using comma
// delimiter. If verbose, specific line parse errors are logged.
func NewCSVImporter(path string, comma rune, verbose bool) *CSVImporter {
	return &CSVImporter{
		path:        path,
		comma:       comma,
		entriesChan: make(chan *geoloc.Entry, entriesChanCap),
		once:        &sync.Once{},
		stats:       &Stats{},
		verbose:     verbose,
	}
}

// Entries returns chan with import results.
func (c *CSVImporter) Entries() <-chan *geoloc.Entry {
	return c.entriesChan
}

// Exec synchronously imports data one time from a CSV file. It returns import
// process stats and an error, if any.
func (c *CSVImporter) Exec() (stats Stats, err error) {
	c.once.Do(func() {
		err = c.run()
	})

	return *c.stats, err
}

func (c *CSVImporter) run() error {
	defer close(c.entriesChan)

	f, err := os.Open(c.path)
	if err != nil {
		return fmt.Errorf("could not open CSV file: %w", err)
	}

	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	reader.Comma = c.comma
	reader.ReuseRecord = true

	start := time.Now()
	defer func() {
		c.stats.ElapsedTime = time.Since(start)
	}()

	uniqueIPs := make(map[string]struct{}, uniqueIPsInitialCap)

	for i := 0; ; i++ {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		c.stats.ReadLines++

		entry, err := geoloc.NewEntryFromStringSlice(record)
		if err != nil {
			if c.verbose {
				log.Printf("[W] importer: could not parse geoloc entry from record on line %d: %s\n", i, err)
			}
			continue
		}

		if _, ok := uniqueIPs[entry.IP.String()]; ok {
			if c.verbose {
				log.Printf("[W] importer: got duplicate record on line %d\n", i)
			}
			continue
		}
		uniqueIPs[entry.IP.String()] = struct{}{}

		if err := entry.Validate(); err != nil {
			if c.verbose {
				log.Printf("[W] importer: invalid geoloc entry from record on line %d: %s\n", i, err)
			}
			continue
		}

		c.stats.ParsedLines++
		c.entriesChan <- entry
	}

	return nil
}
