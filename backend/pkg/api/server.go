package api

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Router     *chi.Mux
	Log        *zap.SugaredLogger
	mqttClient mqtt.Client
	cfg        Config
	dbClient   db.DBStoreClient
}

type DiscordConfig struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
	JWTKey       string
}

type MQTTConfig struct {
	MQTTUser string
	MQTTPass string
	MQTTHost string
}

type Config struct {
	ShareStore        ShareStore
	UserStore         UserStore
	Discord           DiscordConfig
	DBAddr            string
	PreviewStore      PreviewStore
	RoleCheck         RoleChecker
	AESDecryptionKeys map[string][]byte
	// mqtt for notification purposes
	MQTTConfig MQTTConfig
}

type ContextKey string

const (
	TTLContextKey   ContextKey = "ttl"
	UserContextKey  ContextKey = "user"
	ShareContextKey ContextKey = "share"
	DBTagContextKey ContextKey = "db-tag"
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

	// sanity checks
	if s.cfg.ShareStore == nil {
		return nil, fmt.Errorf("no result store provided")
	}
	if s.cfg.UserStore == nil {
		return nil, fmt.Errorf("no user store provided")
	}

	// connect to db
	conn, err := grpc.Dial(cfg.DBAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	s.dbClient = db.NewDBStoreClient(conn)

	// connect to mqtt
	opts := mqttOpts(cfg)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := client.Subscribe("gcsim/#", 1, s.handlePublishedMsgs); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	s.mqttClient = client

	s.Log.Info("server is ready")

	return s, nil
}

func mqttOpts(cfg Config) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.SetKeepAlive(60 * time.Second)
	opts.AddBroker(cfg.MQTTConfig.MQTTHost)
	opts.SetClientID("gcsim-api-server")
	if cfg.MQTTConfig.MQTTUser != "" {
		opts.SetUsername(cfg.MQTTConfig.MQTTUser)
		opts.SetPassword(cfg.MQTTConfig.MQTTPass)
	}
	return opts
}

func (s *Server) handlePublishedMsgs(client mqtt.Client, msg mqtt.Message) {
	s.Log.Infow("mqtt msg received", "msg", string(msg.Payload()))
}

func (s *Server) routes() {
	s.Log.Debugw("setting up server routes")
	r := s.Router

	// r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(s.tokenCheck)

	r.Route("/api", func(r chi.Router) {
		r.Route("/preview", func(r chi.Router) {
			r.Get("/{share-key}", s.GetPreview())
			r.Get("/db/{db-key}", s.GetPreviewByDBID())
		})

		r.Route("/share", func(r chi.Router) {
			r.Post("/", s.CreateShare())        // share a sim
			r.Get("/{share-key}", s.GetShare()) // get a shared sim
			r.Get("/random", s.GetRandomShare())
			r.Get("/db/{db-key}", s.GetShareByDBID())
		})

		r.Get("/login", s.Login())

		r.Route("/user", func(r chi.Router) {
			r.Post("/save", s.UserSave())
		})

		r.Route("/db", func(r chi.Router) {
			r.Get("/", s.getDB())
			r.Get("/id/{id}", s.getByID())
		})
	})
}
