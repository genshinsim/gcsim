package queue

import (
	context "context"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type NotifyService interface {
	Notify(topic string, message string)
}

type WorkSource interface {
	GetWork(context.Context) ([]*model.ComputeWork, error)
}

type Config struct {
	DBWork  WorkSource
	SubWork WorkSource
	Timeout time.Duration
}

type Server struct {
	Config
	Log *zap.SugaredLogger
	UnimplementedWorkQueueServer
	getWork      chan getWorkReq
	completeWork chan completeWorkReq
}

type getWorkReq struct {
	resp chan *model.ComputeWork
}

type completeWorkReq struct {
	id   string
	resp chan error
}

type Work struct {
	Key  string
	Work any
}

func NewQueue(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		Config:       cfg,
		getWork:      make(chan getWorkReq),
		completeWork: make(chan completeWorkReq),
	}
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

		s.Log = sugar
	}

	if s.Timeout <= 0 {
		s.Timeout = 5 * time.Minute //default 5 min
	}

	go s.queueCtrl()

	s.Log.Infow("queue service started")
	return s, nil
}

type wipWork struct {
	w      *model.ComputeWork
	Expiry time.Time
}

func (s *Server) Get(ctx context.Context, req *GetReq) (*GetResp, error) {
	s.Log.Infow("get work request received")
	resp := make(chan *model.ComputeWork)
	s.getWork <- getWorkReq{
		resp: resp,
	}
	res := <-resp
	return &GetResp{
		Data: res,
	}, nil
}

func (s *Server) Complete(ctx context.Context, req *CompleteReq) (*CompleteResp, error) {
	id := req.GetId()
	s.Log.Infow("complete work request", "id", id)
	resp := make(chan error)
	s.completeWork <- completeWorkReq{
		id:   id,
		resp: resp,
	}
	err := <-resp
	s.Log.Infow("completed", "err", err)
	if err != nil {
		return nil, err
	}
	return &CompleteResp{}, nil
}
