package http

import (
	"log"
	"net"
	"net/http"

	"geolocation-resolver/pkg/geoloc"

	"github.com/francoispqt/gojay"
)

// Resolver provides access to geolocation data by IP.
type Resolver interface {
	// ResolveIP resolves geolocation data for ip address. If geolocation data
	// for ip is not found, nil error must be returned.
	ResolveIP(net.IP) (*geoloc.Entry, error)
}

// Handler handles one type of requests that, given an IP address, returns
// information about the IP address' location (country and city).
type Handler struct {
	resolver Resolver
}

// NewHandler creates Handler that handles resolve IP requests using Resolver.
func NewHandler(resolver Resolver) *Handler {
	return &Handler{
		resolver: resolver,
	}
}

func (h *Handler) respond(
	w http.ResponseWriter,
	code int,
	response gojay.MarshalerJSONObject,
) {
	data, err := gojay.MarshalJSONObject(response)
	if err != nil {
		log.Printf("[E] http: failed to marshal JSON: %s\n", err)

		code = http.StatusInternalServerError
		data = marshalErrResponse
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		log.Printf("[E] http: failed to write response: %s\n", err)
	}
}

func (h *Handler) handleResolveIPRequest(w http.ResponseWriter, r *http.Request) {
	rawIP := r.FormValue("ip")
	if len(rawIP) == 0 {
		h.respond(w, http.StatusBadRequest, &errorResponse{
			Message: "missing ip",
		})
		return
	}

	ip := net.ParseIP(rawIP)
	if len(ip) == 0 {
		h.respond(w, http.StatusBadRequest, &errorResponse{
			Message: "invalid ip",
		})
		return
	}

	entry, err := h.resolver.ResolveIP(ip)
	if err != nil {
		log.Printf("[E] http: failed to resolve IP: %s\n", err)

		h.respond(w, http.StatusInternalServerError, &errorResponse{
			Message: "resolve error",
		})
		return
	}
	if entry == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Printf("[D] http: got entry [%v] for '%s'\n", entry, ip)

	h.respond(w, http.StatusOK, &resolveIPResponse{
		City:    entry.City,
		Country: entry.Country,
	})
}
