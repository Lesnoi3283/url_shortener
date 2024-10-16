package handlers

import (
	"net/http"
)

//go:generate mockgen -destination=mocks/mock_DBInterface.go -package=mocks github.com/Lesnoi3283/url_shortener/internal/app/handlers DBInterface

// PingDBHandler is a handler struct. Use it`s ServeHTTP func.
type PingDBHandler struct {
	DB URLStorageInterface
	//todo: log
}

// ServeHTTP returns http.StatusOK if database is active.
func (p *PingDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := p.DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//стоит ли тут в логгер выводить, что бд не работает?
		//ответ: стоит.
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
