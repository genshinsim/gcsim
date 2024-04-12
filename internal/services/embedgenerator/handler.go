package embedgenerator

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/redis/go-redis/v9"
)

func (s *Server) handleProxy(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("proxying request", "prefix", prefix, "url", r.URL)
		r.Host = s.proxyTarget.Host
		s.proxy.ServeHTTP(w, r)
	}
}

var ErrDir = errors.New("path is dir")

func (s *Server) tryRead(requestedPath string, w http.ResponseWriter) error {
	f, err := s.staticFS.Open(path.Join("dist", requestedPath))
	if err != nil {
		return err
	}
	defer f.Close()

	stat, _ := f.Stat()
	if stat.IsDir() {
		return ErrDir
	}

	contentType := mime.TypeByExtension(filepath.Ext(requestedPath))
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, f)
	return err
}

func (s *Server) notFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("received request on not found handler", "path", r.URL)
		err := s.tryRead(r.URL.Path, w)
		if err == nil {
			s.logger.Info("file found; serving", "path", r.URL)
			return
		}
		s.logger.Info("file not found; trying index.html next", "path", r.URL, "err", err)
		err = s.tryRead("index.html", w)
		if err != nil {
			s.logger.Info("error reading index.html", "path", r.URL, "err", err)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (s *Server) authKeyCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.authKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		authKey := r.Header.Get("X-CUSTOM-AUTH-KEY")
		if authKey != s.authKey {
			s.logger.Info("unauthorized request", "authkey", authKey)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleImageRequest(src string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		s.logger.Info("generate request", "src", src, "id", id)
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// change id to include src
		id = src + "/" + id
		ctx, cancel := context.WithTimeout(r.Context(), s.generateTimeout)
		defer cancel()

		pubsub := s.rdb.Subscribe(ctx, id)
		defer pubsub.Close()

		res := s.rdb.Get(ctx, id)
		s.logger.Info("got get from redis", "res_length", len(res.Val()))
		switch res.Err() {
		case nil:
			if val := res.Val(); !strings.HasPrefix(val, "wip") {
				s.logger.Info("id already in wip")
				s.handleResult(val, w)
				return
			}
			// wait for existing result
		case redis.Nil:
			s.logger.Info("id no result; starting new")
			s.rdb.Set(ctx, id, "wip", s.generateTimeout)
			go s.do(id)
		default:
			// exception case where something goes wrong with redis
			s.logger.Info("unexpected get with non nil err", "id", id, "err", res.Err())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// wip, sub to topic == id and wait for ok or time out
		ch := pubsub.Channel()
		select {
		case <-ctx.Done():
			s.logger.Info("context done before receiving msg", "reason", ctx.Err())
			w.WriteHeader(http.StatusRequestTimeout)
			return
		case msg := <-ch:
			requestID := ctx.Value(middleware.RequestIDKey)
			s.logger.Info("received msg from redis", "request_id", requestID, "msg", msg.Payload)
			if msg.Payload != "done" {
				// this is an error message so just skip getting
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(msg.Payload))
				return
			}
			// key must be available now?
			res := s.rdb.Get(ctx, id)
			if err := res.Err(); err != nil {
				s.logger.Info("redis get unexpected err", "request_id", requestID, "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// res can't be wip still
			val := res.Val()
			if strings.HasPrefix(val, "wip") {
				s.logger.Info("unexpected res is still wip", "request_id", requestID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			s.handleResult(val, w)
			return
		}
	}
}

func (s *Server) handleResult(val string, w http.ResponseWriter) {
	if strings.HasPrefix(val, "error: ") {
		val = strings.TrimPrefix(val, "error: ")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(val))
		return
	}
	data, err := decode(val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
