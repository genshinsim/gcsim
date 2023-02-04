package api

import (
	"compress/gzip"
	"net/http"
)

func (s *Server) getWork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Log.Infow("get compute work request received")
		work, err := s.cfg.DBStore.GetComputeWork(r.Context())
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		if work == nil {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("no work"))
			return
		}

		data, err := work.MarshalJson()
		if err != nil {
			s.Log.Warnw("error get compute work - cannot marshal result", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		writer, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			s.Log.Warnw("error get comptue work - cannot write gzip result", "err", err)
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
