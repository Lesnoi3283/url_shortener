package handlers

import (
	"encoding/json"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic"
	"go.uber.org/zap"
	"net/http"
)

type StatsHandler struct {
	log     zap.SugaredLogger
	storage logic.URLStorageInterface
}

type statsData struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urls, users, err := logic.GetStats(r.Context(), h.storage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Errorf("cant get stats, err: %v", err)
		return
	}

	//prepare and marshal data
	stats := statsData{
		URLs:  urls,
		Users: users,
	}
	JSONStats, err := json.Marshal(stats)
	if err != nil {
		h.log.Errorf("cant marshal stats data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//make a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(JSONStats)
}
