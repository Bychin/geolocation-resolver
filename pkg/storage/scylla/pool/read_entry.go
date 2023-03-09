package pool

import (
	"net"

	"geolocation-resolver/pkg/geoloc"
)

// ReadGeolocEntryRequest is a request for reading geoloc entry from ScyllaDB.
type ReadGeolocEntryRequest struct {
	ip net.IP

	out chan *ReadGeolocEntryResponse
}

// NewReadGeolocEntryRequest creates new request for reading geoloc entry by IP.
func NewReadGeolocEntryRequest(ip net.IP) *ReadGeolocEntryRequest {
	return &ReadGeolocEntryRequest{
		ip:  ip,
		out: make(chan *ReadGeolocEntryResponse, 1),
	}
}

// Response returns a chan with a response to the request r.
func (r *ReadGeolocEntryRequest) Response() <-chan *ReadGeolocEntryResponse {
	return r.out
}

// ReadGeolocEntryResponse is a response to ReadGeolocEntryRequest.
type ReadGeolocEntryResponse struct {
	Entry *geoloc.Entry
	Err   error
}

// readGeolocEntry actually reads geoloc entry from ScyllaDB.
func (r *Reader) readGeolocEntry(req *ReadGeolocEntryRequest) {
	defer close(req.out)

	query := r.session.Query(readGeolocEntryQuery, req.ip)
	defer query.Release()

	entry := &geoloc.Entry{
		IP: req.ip,
	}
	if err := query.Scan(
		&entry.CountryCode,
		&entry.Country,
		&entry.City,
		&entry.Latitude,
		&entry.Longitude,
	); err != nil {
		req.out <- &ReadGeolocEntryResponse{
			Err: err,
		}
		return
	}

	req.out <- &ReadGeolocEntryResponse{
		Entry: entry,
	}
}
