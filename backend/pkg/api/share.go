package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (s *Server) GetShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "share-key")
		share, expiresAt, err := s.cfg.Store.ReadShare(r.Context(), key)
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
		w.Header().Set("x-gcsim-ttl", strconv.FormatUint(expiresAt, 10))
		w.WriteHeader(http.StatusOK)
		w.Write(d)
	}
}
