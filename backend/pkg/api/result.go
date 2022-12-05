package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type ResultStore interface {
	ResultReader
	Create(data []byte, ctx context.Context) (string, error)
	Update(id string, data []byte, ctx context.Context) error
	SetTTL(id string, ctx context.Context) error
	Delete(id string, ctx context.Context) error
	Random(ctx context.Context) (string, error)
}

type ResultReader interface {
	Read(id string, ctx context.Context) ([]byte, uint64, error)
}

var ErrKeyNotFound = errors.New("key does not exist")

const DefaultTLL = 24 * 14

func (s *Server) CreateShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			s.Log.Errorw("unexpected error reading request body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !json.Valid(data) {
			s.Log.Info("request is not valid json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		uuid, err := s.cfg.ResultStore.Create(data, context.WithValue(r.Context(), TTLContextKey, DefaultTLL))

		if err != nil {
			s.Log.Errorw("unexpected error saving result", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(uuid))
	}
}

func (s *Server) GetShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "share-key")

		share, ttl, err := s.cfg.ResultStore.Read(key, r.Context())
		switch err {
		case nil:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("x-gcsim-ttl", strconv.FormatUint(ttl, 10))
			w.WriteHeader(http.StatusOK)
			w.Write(share)
		case ErrKeyNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			s.Log.Errorw("unexpected error getting share", "err", err)
		}

	}
}

func (s *Server) GetRandomShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		share, err := s.cfg.ResultStore.Random(r.Context())
		switch err {
		case nil:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(share))
		case ErrKeyNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			s.Log.Errorw("unexpected error getting share", "err", err)
		}

	}
}
