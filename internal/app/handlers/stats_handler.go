package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type StatsHandler struct {
	log     zap.SugaredLogger
	storage URLStorageInterface
}

type statsData struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//get URLs count
	urls, err := h.storage.GetShortURLCount(r.Context())
	if err != nil {
		h.log.Errorf("cant get a short URLs count: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//get users count
	users, err := h.storage.GetUserCount(r.Context())
	if err != nil {
		h.log.Errorf("cant get a user count: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
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
