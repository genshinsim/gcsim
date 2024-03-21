package api

import (
	"compress/gzip"
	"encoding/json"
	"net/http"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type dbGetOpt struct {
	Query   map[string]interface{} `json:"query"`
	Sort    map[string]interface{} `json:"sort"`
	Project map[string]interface{} `json:"project"`
	Skip    int64                  `json:"skip"`
	Limit   int64                  `json:"limit"`
}

func marshalOptions() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		AllowPartial:    true,
		UseEnumNumbers:  true, // TODO: prob better if we set to false?
		EmitUnpopulated: false,
	}
}

func (s *Server) getDB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var o dbGetOpt
		queryStr := r.URL.Query().Get("q")
		s.Log.Infow("db query request received", "q", queryStr)
		if queryStr != "" {
			err := json.Unmarshal([]byte(queryStr), &o)
			if err != nil {
				s.Log.Infow("error querying db - bad request", "err", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

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
		opt := &db.QueryOpt{
			Query:   query,
			Sort:    sort,
			Project: project,
			Skip:    o.Skip,
			Limit:   o.Limit,
		}

		s.Log.Infow("forwarding request to db", "opt", opt.String())

		res, err := s.dbClient.Get(r.Context(), &db.GetRequest{Query: opt})
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
		data, err := marshalOptions().Marshal(res.GetData())
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

func (s *Server) getByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "id")

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
		data, err := marshalOptions().Marshal(resp.GetData())
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
	}
}
