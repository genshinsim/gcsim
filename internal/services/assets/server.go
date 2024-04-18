package assets

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ServerCfg func(s *Server) error

type Server struct {
	*chi.Mux

	hosts       map[string][]*url.URL // external source for images
	assetPrefix string                // assets prefix
	cacheDir    string                // where to save cache files
	specialDir  string                // where to locate special files

	httpClient *http.Client
	logger     *slog.Logger
}

func New() (*Server, error) {
	s := &Server{
		hosts: make(map[string][]*url.URL),
		httpClient: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
	}

	return s, nil
}

func (s *Server) Init() error {
	if s.logger == nil {
		s.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	err := s.routes()
	if err != nil {
		return err
	}

	return nil
}

func WithProxyTimeout(timeout time.Duration) ServerCfg {
	return func(s *Server) error {
		s.httpClient.Timeout = timeout
		return nil
	}
}

func WithAssetSource(assetType, host string) ServerCfg {
	return func(s *Server) error {
		u, err := url.Parse(host)
		if err != nil {
			return err
		}
		switch assetType {
		case "*":
			s.hosts["avatar"] = append(s.hosts["avatar"], u)
			s.hosts["weapons"] = append(s.hosts["weapons"], u)
			s.hosts["artifacts"] = append(s.hosts["artifacts"], u)
			return nil
		case "avatar":
		case "weapons":
		case "artifacts":
		default:
			return fmt.Errorf("unrecognized asset type: %v", assetType)
		}
		s.hosts[assetType] = append(s.hosts[assetType], u)
		return nil
	}
}

func (s *Server) routes() error {
	s.Use(middleware.Logger)
	s.Use(middleware.Recoverer)
	s.Use(middleware.RequestID)

	path := strings.TrimSuffix(s.assetPrefix, "/")
	s.Route(path+"/", func(r chi.Router) {
		// usually /api/assets/avatar/<name>.png
		r.Handle("/avatar/{key}.png", s.handleGetData("avatar"))
	})
	return nil
}

func (s *Server) handleGetData(t string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		s.logger.Info("request for avatar image", "key", key)

		// special assets are handled separately
		if _, ok := specialAssets[key]; ok {
			s.loadFromSpecial(key, w)
			return
		}

		// try cache first
		data, err := s.loadFromFile(fmt.Sprintf("/%v/%v.png", t, key))
		switch {
		case data != nil && err == nil:
			// load successful
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		case err != nil:
			s.logger.Info("unexpected error trying to read from cache", "t", t, "key", key)
		case data == nil:
			s.logger.Info("cache data not found", "t", t, "key", key)
		}

		hosts := s.hosts[t]

		// else try external 1 at a time
		for i, v := range hosts {
			// joinedUrl := v.ResolveReference(&url.URL{Path: })
			// data, err := s.proxyImageRequest()
			s.logger.Info("trying external image source", "host", v.String(), "try", i, "key", key)
		}

		// save to cache before serving
	}
}
