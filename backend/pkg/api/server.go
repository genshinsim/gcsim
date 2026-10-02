package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/model"
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

var ErrKeyNotFound = errors.New("key does not exist")

type Store interface {
	// ReadShare returns ErrKeyNotFound if no share has this key.
	ReadShare(ctx context.Context, key string) (*model.SimulationResult, uint64, error)
	// ShareKeyByDBID returns ErrKeyNotFound if no db entry has this id.
	ShareKeyByDBID(ctx context.Context, id string) (string, error)
}

type Config struct {
	Store Store
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
	if s.cfg.Store == nil {
		return nil, fmt.Errorf("no store provided")
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
