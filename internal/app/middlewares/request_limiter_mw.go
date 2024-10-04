package middlewares

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

var errLimitReached error = errors.New("limit reached")

func NewLimitReachedError() error {
	return errLimitReached
}

// Requestik is a list node with time and pointer to a next node.
type Requestik struct {
	time time.Time
	next *Requestik
}

type RequestManager struct {
	mutex     sync.RWMutex
	amount    int
	head      *Requestik
	last      *Requestik
	limit     int
	timeLimit time.Duration
}

// NewRequestManager returns new RequestManager
func NewRequestManager(limit int, timeLimit time.Duration) *RequestManager {
	return &RequestManager{limit: limit, head: nil, last: nil, timeLimit: timeLimit}
}

// RequestManager.add adds new Requestik to a map if limit isn`t reached.
func (r *RequestManager) add() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.amount >= r.limit {
		return NewLimitReachedError()
	} else if r.amount == 0 {
		r.head = &Requestik{
			time: time.Now(),
			next: nil,
		}
		r.last = r.head
		r.amount++
		return nil

	} else {
		r.last.next = &Requestik{
			time: time.Now(),
			next: nil,
		}
		r.last = r.last.next
		r.amount++
		return nil
	}
}

// RequestManager.clean cleans ONLY old requestiks (.time > timeLimit), returns amount of deleted requestiks.
func (r *RequestManager) clean() (cleaned int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.head != nil {
		for time.Since(r.head.time) > r.timeLimit {
			r.head = r.head.next
			cleaned++
			if r.head == nil {
				break
			}
		}
	}
	r.amount -= cleaned
	return cleaned
}

// RequestLimiterMW Limits requests from all users by time. Returns a http.StatusTooManyRequests if limit reached.
func RequestLimiterMW(logger zap.SugaredLogger, manager *RequestManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			manager.clean()
			if err := manager.add(); errors.Is(err, NewLimitReachedError()) {
				w.WriteHeader(http.StatusTooManyRequests)
				logger.Infow("Too many requests", zap.String("Request URI", r.RequestURI))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
