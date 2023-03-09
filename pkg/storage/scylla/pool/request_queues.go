package pool

// ReadRequestsQueue is a queue of requests for reading geoloc entries.
type ReadRequestsQueue chan *ReadGeolocEntryRequest

// NewReadRequestsQueue creates a ReadRequestsQueue of size size.
func NewReadRequestsQueue(size int) ReadRequestsQueue {
	return make(chan *ReadGeolocEntryRequest, size)
}

// WriteRequestsQueue is a queue of requests for writing geoloc entries.
type WriteRequestsQueue chan *WriteGeolocEntryRequest

// NewWriteRequestsQueue creates a WriteRequestsQueue of size size.
func NewWriteRequestsQueue(size int) WriteRequestsQueue {
	return make(chan *WriteGeolocEntryRequest, size)
}
