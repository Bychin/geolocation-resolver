package http

import (
	"github.com/francoispqt/gojay"
)

type resolveIPResponse struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

func (r *resolveIPResponse) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddStringKey("country", r.Country)
	enc.AddStringKey("city", r.City)
}

func (r *resolveIPResponse) IsNil() bool {
	return r == nil
}

// marshalErrResponse is used to respond to client when marshal error occurs.
var marshalErrResponse = []byte(`{"message":"response marshaling error"}`)

type errorResponse struct {
	Message string `json:"message"`
}

func (e *errorResponse) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddStringKey("message", e.Message)
}

func (e *errorResponse) IsNil() bool {
	return e == nil
}
