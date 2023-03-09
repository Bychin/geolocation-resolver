package pool

import (
	"github.com/gocql/gocql"
)

// Reader is a worker who executes read requests.
type Reader struct {
	session *gocql.Session

	queue <-chan *ReadGeolocEntryRequest

	done chan struct{}
}

// NewReader creates new worker to execute read requests from q.
func NewReader(s *gocql.Session, q <-chan *ReadGeolocEntryRequest) *Reader {
	return &Reader{
		session: s,
		queue:   q,
		done:    make(chan struct{}),
	}
}

// Init starts reader's routine to consume read queue.
func (r *Reader) Init() {
	go r.routine()
}

// Done returns internal done chan whose closure indicates that reader has
// exited its routine.
func (r *Reader) Done() <-chan struct{} {
	return r.done
}

func (r *Reader) routine() {
	for req := range r.queue {
		r.readGeolocEntry(req)
	}

	close(r.done)
}
