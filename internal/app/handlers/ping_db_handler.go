package handlers

import (
	"net/http"
)

//go:generate mockgen -destination=mocks/mock_DBInterface.go -package=mocks github.com/Lesnoi3283/url_shortener/internal/app/handlers DBInterface

type PingDBHandler struct {
	DB URLStorageInterface
	//todo: log
}

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
