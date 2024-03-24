package handlers

import (
	"net/http"
)

//go:generate mockgen -destination=mocks/mock_DBInterface.go -package=mocks github.com/Lesnoi3283/url_shortener/internal/app/handlers DBInterface

type DBInterface interface {
	Ping() error
}

type pingDBHandler struct {
	db DBInterface
}

func (p *pingDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := p.db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//стоит ли тут в логгер выводить, что бд не работает?
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
