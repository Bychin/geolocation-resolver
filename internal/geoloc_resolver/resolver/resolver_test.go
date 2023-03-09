package resolver

import (
	"context"
	"net"
	"reflect"
	"testing"

	"geolocation-resolver/pkg/geoloc"
)

type dummyDB map[string]geoloc.Entry

func newDummyDB() *dummyDB {
	db := dummyDB(make(map[string]geoloc.Entry))
	return &db
}

func (d *dummyDB) Write(_ context.Context, entry *geoloc.Entry) error {
	(*d)[entry.IP.String()] = *entry
	return nil
}

func (d *dummyDB) Read(_ context.Context, ip net.IP) (*geoloc.Entry, error) {
	entry, ok := (*d)[ip.String()]
	if !ok {
		return nil, nil
	}

	return &entry, nil
}

func TestImportCSVAndResolveIP(t *testing.T) {
	dataPath := "testdata/data_dump.csv"

	ip := net.ParseIP("200.106.141.15")

	expectedEntriesNum := 1
	expectedEntry := &geoloc.Entry{
		IP:          ip,
		CountryCode: "SI",
		Country:     "Nepal",
		City:        "DuBuquemouth",
		Latitude:    -84.87503094689836,
		Longitude:   7.206435933364332,
	}

	cfg := &Config{
		Importer: struct {
			Separator string "yaml:\"separator\""
			Verbose   bool   "yaml:\"verbose\""
		}{
			Separator: ",",
			Verbose:   true,
		},
	}
	db := newDummyDB()

	r := NewResolver(cfg, db)
	if err := r.ImportCSV(dataPath); err != nil {
		t.Errorf("failed to import CSV: %s", err)
	}

	if len(*db) != expectedEntriesNum {
		t.Errorf("wrong amount of entries in db: have %d, want %d",
			len(*db), expectedEntriesNum)
	}

	resolvedEntry, err := r.ResolveIP(ip)
	if err != nil {
		t.Errorf("failed to resolve ip: %s", err)
	}

	if !reflect.DeepEqual(resolvedEntry, expectedEntry) {
		t.Errorf("resolved entry is not equal to expected one, have: %+v, want: %+v",
			resolvedEntry, expectedEntry)
	}
}
