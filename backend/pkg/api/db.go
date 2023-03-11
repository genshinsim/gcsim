package api

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type DBStore interface {
	Get(context.Context, *model.DBQueryOpt) (*model.DBEntries, error)
	GetOne(ctx context.Context, id string) (*model.DBEntry, error)
	Update(ctx context.Context, id string, result *model.SimulationResult) error
	ApproveTag(context.Context, string, model.DBTag) error
	RejectTag(context.Context, string, model.DBTag) error
}

type dbGetOpt struct {
	Query   map[string]interface{} `json:"query"`
	Sort    map[string]interface{} `json:"sort"`
	Project map[string]interface{} `json:"project"`
	Skip    int64                  `json:"skip"`
	Limit   int64                  `json:"limit"`
}

func (s *Server) getDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Log.Infow("db query request received")
		var o dbGetOpt
		queryStr := r.URL.Query().Get("q")
		if queryStr != "" {
			err := json.Unmarshal([]byte(queryStr), &o)
			if err != nil {
				s.Log.Infow("error querying db - bad request", "err", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		s.Log.Infow("db query - query string parsed ok", "query_string", o.Query)

		query, err := structpb.NewStruct(o.Query)
		if err != nil {
			s.Log.Warnw("error querying db - could not convert query to structpb", "err", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		sort, err := structpb.NewStruct(o.Sort)
		if err != nil {
			s.Log.Warnw("error querying db - could not convert sort to structpb", "err", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		project, err := structpb.NewStruct(o.Project)
		if err != nil {
			s.Log.Warnw("error querying db - could not convert project to structpb", "err", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		opt := &model.DBQueryOpt{
			Query:   query,
			Sort:    sort,
			Project: project,
			Skip:    o.Skip,
			Limit:   o.Limit,
		}

		res, err := s.cfg.DBStore.Get(r.Context(), opt)
		if err != nil {
			s.Log.Warnw("error querying db", "err", err)
			if st, ok := status.FromError(err); ok {
				// Error was not a status error
				switch st.Code() {
				case codes.NotFound:
					http.Error(w, "internal server error", http.StatusInternalServerError)
				default:
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
				return
			}

			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		data, err := res.MarshalJson()
		if err != nil {
			s.Log.Warnw("error query db - cannot marshal result", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		writer, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			s.Log.Warnw("error query db - cannot write gzip result", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer writer.Close()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		writer.Write(data)
		// w.Write(data)

	}
}

func (s *Server) approveTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.Context().Value(DBTagContextKey).(model.DBTag)
		key := chi.URLParam(r, "db-key")
		err := s.cfg.DBStore.ApproveTag(r.Context(), key, tag)
		if err != nil {
			if st, ok := status.FromError(err); st.Code() == codes.NotFound && ok {
				http.Error(w, "id not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) rejectTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.Context().Value(DBTagContextKey).(model.DBTag)
		key := chi.URLParam(r, "db-key")
		err := s.cfg.DBStore.RejectTag(r.Context(), key, tag)
		if err != nil {
			if st, ok := status.FromError(err); st.Code() == codes.NotFound && ok {
				http.Error(w, "id not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
