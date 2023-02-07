package compute

import (
	context "context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	APIKey        string
	ResultHandler ResultHandler
	Timeout       time.Duration
	WorkerCount   int
	Iterations    int
}

type ResultHandler interface {
	Post(key string, callbackURL string, result *model.SimulationResult, err error)
}

type Server struct {
	Log *zap.SugaredLogger
	UnimplementedComputeServer
	cfg          Config
	mu           sync.Mutex
	runnerIsBusy bool
	work         chan *RunRequest
}

func New(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		cfg:  cfg,
		work: make(chan *RunRequest),
	}

	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Log == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	//sanity check
	if cfg.ResultHandler == nil {
		return nil, errors.New("no result handler supplied")
	}
	if s.cfg.WorkerCount <= 0 {
		s.cfg.WorkerCount = runtime.NumCPU() * 2
	}
	if s.cfg.Iterations <= 0 {
		s.cfg.Iterations = 1000
	}

	go s.pool()

	return s, nil
}

func (s *Server) Run(ctx context.Context, req *RunRequest) (*RunResponse, error) {
	s.Log.Infow("run request received", "key", req.GetKey())
	//check api key
	if req.ApiKey != s.cfg.APIKey {
		s.Log.Warnw("invalid key api key in request", "key", req.GetKey(), "api_key", req.GetApiKey())
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	//check if busy
	s.mu.Lock()
	s.Log.Infow("checking if runner is busy", "key", req.GetKey())
	if s.runnerIsBusy {
		s.mu.Unlock()
		s.Log.Infow("runner is currently busy", "key", req.GetKey())
		return &RunResponse{}, status.Error(codes.Unavailable, "resource busy")
	}
	s.mu.Unlock()

	s.Log.Infow("runner not busy; scheduling work", "key", req.GetKey())

	//do work
	s.work <- req

	return &RunResponse{}, nil

}

func (s *Server) pool() {
	done := make(chan interface{})
next:
	for {
		select {
		case req := <-s.work:
			s.mu.Lock()
			if s.runnerIsBusy {
				s.mu.Unlock()
				s.Log.Infow("unexpected runner is busy in pool", "key", req.GetKey())
				continue next
			}
			s.runnerIsBusy = true
			s.mu.Unlock()
			go s.runner(req, done)
		case <-done:
			s.mu.Lock()
			s.runnerIsBusy = false
			s.mu.Unlock()
		}
	}

}

func (s *Server) runner(work *RunRequest, done chan interface{}) {
	defer func() {
		done <- nil
	}()
	key := work.GetKey()
	s.Log.Infow("runner started", "key", key, "iterations", s.cfg.Iterations, "workers", s.cfg.WorkerCount)
	// compute result
	simcfg, err := simulator.Parse(work.GetConfig())
	if err != nil {
		s.Log.Infow("runner could parse config provided", "key", key, "err", err)
		s.cfg.ResultHandler.Post(key, work.CallbackUrl, nil, err)
	}
	simcfg.Settings.Iterations = s.cfg.Iterations //force it to 1000 iterations
	simcfg.Settings.NumberOfWorkers = s.cfg.WorkerCount

	//TODO: timeout should be adjusted to something more reasonable
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
	defer cancel()

	//TODO: need to insert version, builddate
	result, err := simulator.RunWithConfig(work.GetConfig(), simcfg, simulator.Options{}, time.Now(), ctx)
	if err != nil {
		s.Log.Infow("runner encounted error running sim", "key", key, "err", err)
		s.cfg.ResultHandler.Post(key, work.CallbackUrl, nil, err)
	}

	s.cfg.ResultHandler.Post(key, work.CallbackUrl, result, err)
}
