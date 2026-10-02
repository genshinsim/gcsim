package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
)

type DBShareKeyStore interface {
	// ShareKeyByDBID returns ErrKeyNotFound if no db entry has this id.
	ShareKeyByDBID(ctx context.Context, id string) (string, error)
}

func (s *Server) GetDBShareKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		key, err := s.cfg.DBShareKeys.ShareKeyByDBID(r.Context(), id)
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
