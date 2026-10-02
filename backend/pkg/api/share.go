package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/go-chi/chi"
)

type ShareStore interface {
	// Read returns ErrKeyNotFound if no share has this key.
	Read(ctx context.Context, id string) (*model.SimulationResult, uint64, error)
}

func (s *Server) GetShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "share-key")
		share, ttl, err := s.cfg.ShareStore.Read(r.Context(), key)
		if err != nil {
			if errors.Is(err, ErrKeyNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			s.Log.Errorw("unexpected error getting share", "err", err)
			return
		}
		d, err := share.MarshalJSON()
		if err != nil {
			s.Log.Errorw("unexpected error marshalling to json", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-gcsim-ttl", strconv.FormatUint(ttl, 10))
		w.WriteHeader(http.StatusOK)
		w.Write(d)
	}
}
