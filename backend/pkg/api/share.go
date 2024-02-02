package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShareStore interface {
	Read(ctx context.Context, id string) (*model.SimulationResult, uint64, error)
	Create(context.Context, *model.SimulationResult, uint64, string) (string, error)
	SetTTL(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Random(ctx context.Context) (string, error)
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

		str := r.Header.Get("X-GCSIM-SHARE-AUTH")
		if str == "" {
			s.Log.Infow("create share request failed - no hash received", "header", r.Header)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		s.Log.Infow("create share request received", "hash", str)

		err = s.validateSigning(data, str)
		if err != nil {
			s.Log.Infow("create share request - validation failed", "err", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		res := &model.SimulationResult{}
		err = res.UnmarshalJSON(data)
		if err != nil {
			s.Log.Infow("create share request - unmarshall failed", "err", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		user := ""
		userCtx := r.Context().Value(UserContextKey)
		if userCtx != nil {
			user = userCtx.(string)
		}

		expiresAt := uint64(time.Now().Unix() + DefaultTLL)
		id, err := s.cfg.ShareStore.Create(r.Context(), res, expiresAt, user)

		if err != nil {
			s.Log.Errorw("unexpected error saving result", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.Log.Infow("create share request success", "key", id)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(id))
	}
}

func (s *Server) sendShare(w http.ResponseWriter, r *http.Request, key string) {
	share, ttl, err := s.cfg.ShareStore.Read(r.Context(), key)
	if err != nil {
		if st, ok := status.FromError(err); st.Code() == codes.NotFound && ok {
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

func (s *Server) GetShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "share-key")
		s.sendShare(w, r, key)
	}
}

func (s *Server) GetShareByDBID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "db-key")

		resp, err := s.dbClient.GetOne(r.Context(), &db.GetOneRequest{Id: key})
		if err != nil {
			if st, ok := status.FromError(err); st.Code() == codes.NotFound && ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			s.Log.Errorw("unexpected error getting share", "err", err)
			return
		}
		s.sendShare(w, r, resp.GetData().GetShareKey())
	}
}

func (s *Server) GetRandomShare() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		share, err := s.cfg.ShareStore.Random(r.Context())
		switch {
		case err == nil:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(share))
		case errors.Is(err, ErrKeyNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			s.Log.Errorw("unexpected error getting share", "err", err)
		}
	}
}
