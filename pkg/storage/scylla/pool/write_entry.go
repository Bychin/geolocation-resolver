package pool

import (
	"geolocation-resolver/pkg/geoloc"
)

// WriteGeolocEntryRequest is a request for writing geoloc entry to ScyllaDB.
type WriteGeolocEntryRequest struct {
	entry *geoloc.Entry

	out chan *WriteGeolocEntryResponse
}

// NewWriteGeolocEntryRequest creates new request for writing geoloc entry.
func NewWriteGeolocEntryRequest(entry *geoloc.Entry) *WriteGeolocEntryRequest {
	return &WriteGeolocEntryRequest{
		entry: entry,
		out:   make(chan *WriteGeolocEntryResponse, 1),
	}
}

// Response returns a chan with a response to the request w.
func (w *WriteGeolocEntryRequest) Response() <-chan *WriteGeolocEntryResponse {
	return w.out
}

// WriteGeolocEntryResponse is a response to WriteGeolocEntryRequest.
type WriteGeolocEntryResponse struct {
	Err error
}

// writeGeolocEntry actually writes geoloc entry to ScyllaDB.
func (w *Writer) writeGeolocEntry(req *WriteGeolocEntryRequest) {
	defer close(req.out)

	query := w.session.Query(writeGeolocEntryQuery).Bind(
		req.entry.IP,
		req.entry.CountryCode,
		req.entry.Country,
		req.entry.City,
		req.entry.Latitude,
		req.entry.Longitude,
	)
	defer query.Release()

	iter := query.Iter()
	queryErr := iter.Close()

	req.out <- &WriteGeolocEntryResponse{
		Err: queryErr,
	}
}
