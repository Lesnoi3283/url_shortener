package middlewares

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

var errLimitReached error = errors.New("limit reached")

func NewLimitReachedError() error {
	return errLimitReached
}

type Requestik struct {
	time time.Time
	next *Requestik
}

// Есть 2 идеи:
//  1. Очищать просроченные запросы в самой мидлваре
//  2. Запустить горутину, которая будет следить за первым запросом и удалять его по истечению времени
type RequestManager struct {
	mutex     sync.RWMutex
	amount    int
	head      *Requestik
	last      *Requestik
	limit     int
	timeLimit time.Duration
}

func NewRequestManager(limit int, timeLimit time.Duration) *RequestManager {
	return &RequestManager{limit: limit, head: nil, last: nil, timeLimit: timeLimit}
}

func (r *RequestManager) add() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.amount == 0 {
		r.head = &Requestik{
			time: time.Now(),
			next: nil,
		}
		r.last = r.head
		r.amount++
		return nil

	} else if r.amount >= r.limit {
		return NewLimitReachedError()
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

// returns amount of deleted requestiks
func (r *RequestManager) clean() (cleaned int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.head != nil {
		for time.Since(r.head.time) > r.timeLimit {
			r.head = r.head.next
			cleaned++
		}
	}
	r.amount -= cleaned
	return cleaned
}

//func (r *RequestManager) isFull() bool {
//	r.mutex.RLock()
//	defer r.mutex.RUnlock()
//
//	return (r.amount >= r.limit)
//}

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
