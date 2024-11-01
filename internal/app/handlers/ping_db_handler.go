package handlers

import (
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=mocks/mock_DBInterface.go -package=mocks github.com/Lesnoi3283/url_shortener/internal/app/handlers DBInterface

// PingDBHandler is a handler struct. Use it`s ServeHTTP func.
type PingDBHandler struct {
	DB  logic.URLStorageInterface
	log zap.SugaredLogger
}

// ServeHTTP returns http.StatusOK if database is active.
func (p *PingDBHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := logic.PingDB(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.log.Errorf("DB ping err: %v", err)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
