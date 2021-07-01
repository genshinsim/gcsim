package server

import (
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ConfigDir string
	Port      int
}

var port = 8080

type Server struct {
	Router *chi.Mux
	Log    *zap.SugaredLogger
	Cfg    Config

	wsClosed     chan struct{}
	wsBroadcast  chan []byte
	wsRegister   chan wsClient
	wsUnregister chan *websocket.Conn
}

func New(cfg ...func(*Server) error) (*Server, error) {
	s := &Server{
		wsClosed:     make(chan struct{}),
		wsBroadcast:  make(chan []byte, 5),
		wsRegister:   make(chan wsClient),
		wsUnregister: make(chan *websocket.Conn),
	}

	s.Router = chi.NewRouter()
	r := s.Router
	r.Use(middleware.Logger)

	r.Handle("/ws", s.handleUpgrade())

	//custom configs
	for _, f := range cfg {
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

	port = s.Cfg.Port

	go s.wsHub()

	return s, nil
}

func handleErr(r wsRequest, s int, msg string) {
	e := wsResponse{
		ID:      r.ID,
		Status:  s,
		Payload: msg,
	}
	data, _ := json.Marshal(e)
	r.client.send <- data
}

// func (s *Server) ni(ctx context.Context, r wsRequest) {
// 	e := wsResponse{
// 		ID:      r.ID,
// 		Status:  http.StatusOK,
// 		Payload: "function not implemented",
// 	}
// 	msg, _ := json.Marshal(e)
// 	r.client.send <- msg
// }
