// package servermode implements a local server that the web ui can connect to
// for running simulation on a local compute instead of using wasm
package servermode

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type ConfigOpt func(cfg *Config) error

type Config struct {
	ShareKey      string
	Timeout       time.Duration
	WorkerCount   int
	FlushInterval int
	Log           *slog.Logger
}

type Server struct {
	Router *chi.Mux
	Config

	// track work
	sync.Mutex                    // lock for when we need to update results i.e. at flush or final
	pool       map[string]*worker // tracker workers
}

func New(opts ...ConfigOpt) (*Server, error) {
	cfg := &Config{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	s := &Server{
		Config: *cfg,
		pool:   make(map[string]*worker),
	}
	s.Router = chi.NewRouter()
	s.routes()

	s.Log.Info("server initialized")

	return s, nil
}

func WithDefaults() ConfigOpt {
	return func(cfg *Config) error {
		cfg.Log = slog.New(slog.NewTextHandler(os.Stdout, nil))
		cfg.Timeout = time.Minute * 10
		cfg.WorkerCount = 10
		cfg.FlushInterval = 25
		return nil
	}
}

func WithShareKey(key string) ConfigOpt {
	return func(cfg *Config) error {
		cfg.ShareKey = key
		return nil
	}
}

func (s *Server) routes() {
	s.Log.Debug("setting up server routes for preview generation server")
	s.Router.Use(middleware.Logger)
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-GCSIM-SHARE-AUTH"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	s.Router.Get("/ready/{id}", s.ready())
	s.Router.Get("/running/{id}", s.running())
	s.Router.Post("/validate/{id}", s.validate())
	s.Router.Post("/sample/{id}", s.sample())
	s.Router.Post("/run/{id}", s.run())
	s.Router.Get("/results/{id}", s.latest())
	s.Router.Post("/cancel/{id}", s.cancel())
	s.Router.Get("/info", s.info())
}
