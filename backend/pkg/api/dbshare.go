package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) GetDBShareKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		key, err := s.cfg.Store.ShareKeyByDBID(r.Context(), id)
		switch {
		case errors.Is(err, ErrKeyNotFound) || (err == nil && key == ""):
			http.Error(w, "not found", http.StatusNotFound)
			return
		case err != nil:
			s.Log.Errorw("unexpected error looking up db share key", "id", id, "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			ShareKey string `json:"share_key"`
		}{key})
	}
}
