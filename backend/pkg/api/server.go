package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	Router *chi.Mux
	Log    *zap.SugaredLogger
	cfg    Config
}

type DiscordConfig struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
	JWTKey       string
}

type Config struct {
	ResultStore       ResultStore
	UserStore         UserStore
	Discord           DiscordConfig
	DBStore           DBStore
	AESDecryptionKeys map[string][]byte
}

type APIContextKey string

const (
	TTLContextKey  APIContextKey = "ttl"
	UserContextKey APIContextKey = "user"
)

func New(cfg Config, cust ...func(*Server) error) (*Server, error) {

	s := &Server{
		cfg: cfg,
	}

	s.Router = chi.NewRouter()
	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Log == nil {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err := config.Build()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	s.routes()

	//sanity checks
	if s.cfg.ResultStore == nil {
		return nil, fmt.Errorf("no result store provided")
	}
	if s.cfg.DBStore == nil {
		return nil, fmt.Errorf("no db store provided")
	}

	return s, nil
}

func (s *Server) routes() {
	s.Log.Debugw("setting up server routes")
	r := s.Router

	// r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(s.tokenCheck)

	r.Route("/api", func(r chi.Router) {

		r.Route("/share", func(r chi.Router) {
			r.Post("/", s.CreateShare())        // share a sim
			r.Get("/{share-key}", s.GetShare()) // get a shared sim
			r.Get("/random", s.GetRandomShare())
			r.Get("/preview/{share-key}", s.notImplemented()) // preview (embed) for a shared sim
		})

		r.Get("/login", s.Login())

		r.Route("/user", func(r chi.Router) {
			r.Post("/save", s.UserSave())
		})

		r.Route("/db", func(r chi.Router) {
			r.Get("/", s.getDB())
		})
	})

}

func (s *Server) notImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
