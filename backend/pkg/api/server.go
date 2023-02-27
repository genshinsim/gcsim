package api

import (
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	Router     *chi.Mux
	Log        *zap.SugaredLogger
	mqttClient mqtt.Client
	cfg        Config
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
	DBStore           DBStore
	SubmissionStore   SubmissionStore
	AESDecryptionKeys map[string][]byte
	//mqtt for notification purposes
	MQTTConfig MQTTConfig
	//queue service for getting work etc
	QueueService QueueService
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
	if s.cfg.ShareStore == nil {
		return nil, fmt.Errorf("no result store provided")
	}
	if s.cfg.DBStore == nil {
		return nil, fmt.Errorf("no db store provided")
	}

	//connect to mqtt
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
			r.Post("/submit", s.submitEntry())

			r.Get("/work", s.getWork())
			r.Post("/work", s.computeCallback())
		})
	})

}

func (s *Server) notImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
