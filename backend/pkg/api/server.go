package api

import (
	"fmt"

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

type Config struct {
	ShareStore  ShareStore
	DBShareKeys DBShareKeyStore
}

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

	// sanity checks
	if s.cfg.ShareStore == nil {
		return nil, fmt.Errorf("no result store provided")
	}
	if s.cfg.DBShareKeys == nil {
		return nil, fmt.Errorf("no db share key store provided")
	}

	s.Log.Info("server is ready")

	return s, nil
}

func (s *Server) routes() {
	s.Log.Debugw("setting up server routes")
	r := s.Router

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/share/{share-key}", s.GetShare())
		r.Get("/dbshare/{id}", s.GetDBShareKey())
	})
}
