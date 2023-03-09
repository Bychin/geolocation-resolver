package geoloc

import (
	"fmt"
	"net"
	"strconv"
)

const (
	ISOCountryCodeLen = 2 // ISO 3166-1 alpha-2
)

var (
	minStringSliceLenForNewEntry = len([]string{
		"ip_address",
		"country_code",
		"country",
		"city",
		"latitude",
		"longitude",
	})
)

// Entry is a geolocation entry about an IP.
type Entry struct {
	IP          net.IP
	CountryCode string
	Country     string
	City        string
	Latitude    float64
	Longitude   float64
}

// NewEntryFromStringSlice creates new geolocation entry from raw string slice.
// It always returns either a non-nil record or a non-nil error, but not both.
func NewEntryFromStringSlice(raw []string) (*Entry, error) {
	if len(raw) < minStringSliceLenForNewEntry {
		return nil, fmt.Errorf("invalid input len")
	}

	ip := net.ParseIP(raw[0])
	if ip == nil {
		return nil, fmt.Errorf("invalid ip '%s'", raw[0])
	}

	latitude, err := strconv.ParseFloat(raw[4], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude '%s': %w", raw[4], err)
	}

	longitude, err := strconv.ParseFloat(raw[5], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude '%s': %w", raw[5], err)
	}

	return &Entry{
		IP:          ip,
		CountryCode: raw[1],
		Country:     raw[2],
		City:        raw[3],
		Latitude:    latitude,
		Longitude:   longitude,
	}, nil
}

// Validate validates that e is a valid geolocation entry.
func (e *Entry) Validate() error {
	if e.IP == nil {
		return fmt.Errorf("missing ip")
	}
	if len(e.CountryCode) != ISOCountryCodeLen {
		return fmt.Errorf("invalid country code")
	}
	if len(e.Country) == 0 {
		return fmt.Errorf("missing country")
	}
	if len(e.City) == 0 {
		return fmt.Errorf("missing city")
	}
	if e.Latitude == 0 && e.Longitude == 0 {
		return fmt.Errorf("missing coordinates")
	}

	return nil
}
