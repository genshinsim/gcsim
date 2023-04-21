package api

import (
	"context"
	"net/http"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PreviewStore interface {
	Get(context.Context, string, *model.SimulationResult) ([]byte, error)
}

func (s *Server) sendPreview(w http.ResponseWriter, r *http.Request, key string) {
	share, _, err := s.cfg.ShareStore.Read(r.Context(), key)
	if err != nil {
		if st, ok := status.FromError(err); st.Code() == codes.NotFound && ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		s.Log.Errorw("unexpected error getting share", "err", err)
		return
	}

	img, err := s.cfg.PreviewStore.Get(r.Context(), key, share)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		s.Log.Errorw("unexpected error generate img", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/octet-stream")
	w.Write(img)
}

func (s *Server) GetPreview() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "share-key")
		s.sendPreview(w, r, key)
	}
}

func (s *Server) GetPreviewByDBID() http.HandlerFunc {
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

		s.sendPreview(w, r, resp.GetData().GetShareKey())
	}
}
