package importer

import "time"

// Stats contains import process statistics.
type Stats struct {
	// ReadLines is the number of successfully read lines from file.
	ReadLines int
	// ParsedLines is the number of successfully parsed lines from file.
	ParsedLines int
	// ElapsedTime is the amount of time it took to import.
	ElapsedTime time.Duration
}
