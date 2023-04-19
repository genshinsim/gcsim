package backend

import (
	"context"
	"errors"
	"regexp"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
)

type ShareReader interface {
	Read(context.Context, string) (*model.SimulationResult, uint64, error)
}

type SubmissionStore interface {
	Submit(context.Context, *model.Submission) (string, error)
}

type Config struct {
	LinkValidationRegex *regexp.Regexp
	ShareStore          ShareReader
	SubmissionStore     SubmissionStore
}

type Store struct {
	Config
	Log *zap.SugaredLogger
}

func New(cfg Config, cust ...func(*Store) error) (*Store, error) {
	s := &Store{
		Config: cfg,
	}
	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Log == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	if s.ShareStore == nil {
		return nil, errors.New("share store cannot be nil")
	}

	if s.SubmissionStore == nil {
		return nil, errors.New("submission store cannot be nil")
	}

	return s, nil
}
