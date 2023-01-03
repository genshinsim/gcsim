package api

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type DBStore interface {
	Create(context.Context, *model.DBEntry) (string, error)
	Get(context.Context, *structpb.Struct, int64) (*model.DBEntries, error)
}

type dbGetOpt struct {
	Query   map[string]interface{} `json:"query"`
	Sort    map[string]interface{} `json:"sort"`
	Project map[string]interface{} `json:"project"`
	Skip    int                    `json:"skip"`
	Limit   int                    `json:"limit"`
}

func (s *Server) getDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Log.Infow("db query request received")
		var o dbGetOpt
		queryStr := r.URL.Query().Get("q")
		if queryStr != "" {
			err := json.Unmarshal([]byte(queryStr), &o.Query)
			if err != nil {
				s.Log.Infow("error querying db - bad request", "err", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		query, err := structpb.NewStruct(o.Query)
		if err != nil {
			s.Log.Warnw("error querying db - could not convert to structpb", "err", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		res, err := s.cfg.DBStore.Get(r.Context(), query, 1)
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
		data, err := protojson.Marshal(res)
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
