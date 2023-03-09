package hashmap

import (
	"context"
	"fmt"
	"net"
	"sync"

	"geolocation-resolver/pkg/geoloc"
)

const (
	initialStorageCap = 1024
)

var (
	ErrUninitialized = fmt.Errorf("storage is uninitialized")
)

// Storage is a simple hashmap with IP as a key and geoloc entry as a value. It
// is safe for concurrent use.
type Storage struct {
	storage   map[string]geoloc.Entry
	storageMu sync.RWMutex
}

// NewStorage creates new hashmap storage.
func NewStorage() *Storage {
	return &Storage{
		storage:   make(map[string]geoloc.Entry, initialStorageCap),
		storageMu: sync.RWMutex{},
	}
}

// Write synchronously writes entry to hashmap storage.
func (s *Storage) Write(_ context.Context, entry *geoloc.Entry) error {
	if s == nil {
		return ErrUninitialized
	}

	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	s.storage[entry.IP.String()] = *entry
	return nil
}

// Read synchronously reads entry by ip from hashmap storage. If the entry is
// not found, nil error is returned.
func (s *Storage) Read(_ context.Context, ip net.IP) (*geoloc.Entry, error) {
	if s == nil {
		return nil, ErrUninitialized
	}

	s.storageMu.RLock()
	defer s.storageMu.RUnlock()

	entry, ok := s.storage[ip.String()]
	if !ok {
		return nil, nil
	}

	return &entry, nil
}
