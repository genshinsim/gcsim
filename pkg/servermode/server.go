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

type Config struct {
	ShareKey    string
	Timeout     time.Duration
	WorkerCount int
}

type Server struct {
	Router *chi.Mux
	Log    *slog.Logger
	cfg    Config

	// track work
	sync.Mutex                    // lock for when we need to update results i.e. at flush or final
	pool       map[string]*worker // tracker workers
}

func New(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		cfg:  cfg,
		pool: make(map[string]*worker),
	}
	s.Router = chi.NewRouter()
	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}
	if s.Log == nil {
		s.Log = slog.New(slog.NewTextHandler(os.Stdout, nil))
		s.Log.Info("logger initiated")
	}

	s.routes()

	return s, nil
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
