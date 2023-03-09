package pool

import (
	"github.com/gocql/gocql"
)

// Writer is a worker who executes write requests.
type Writer struct {
	session *gocql.Session

	queue <-chan *WriteGeolocEntryRequest

	done chan struct{}
}

// NewWriter creates new worker to execute write requests from q.
func NewWriter(s *gocql.Session, q <-chan *WriteGeolocEntryRequest) *Writer {
	return &Writer{
		session: s,
		queue:   q,
		done:    make(chan struct{}),
	}
}

// Init starts writer's routine to consume write queue.
func (w *Writer) Init() {
	go w.routine()
}

// Done returns internal done chan whose closure indicates that writer has
// exited its routine.
func (w *Writer) Done() <-chan struct{} {
	return w.done
}

func (w *Writer) routine() {
	for req := range w.queue {
		w.writeGeolocEntry(req)
	}

	close(w.done)
}
