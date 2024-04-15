package embedgenerator

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/redis/go-redis/v9"
)

type serverCfg func(s *Server) error

type Server struct {
	*chi.Mux

	logger     *slog.Logger
	rdb        redis.UniversalClient
	work       chan string
	l          *launcher.Launcher
	previewURL string
	authKey    string

	// static asset; mandatory to serve from
	staticDir string

	// additional local assets
	useLocalAssets bool
	assetsPrefix   string // cannot be blank
	assetsDir      string

	// proxy requests
	useProxy     bool
	proxyPrefix  string
	skipInsecure bool

	// for proxying api requests
	proxy       *httputil.ReverseProxy
	proxyTarget *url.URL

	// timeouts
	generateTimeout time.Duration
	cacheTTL        time.Duration

	// browser for navigating to pages
	browser *rod.Browser
}

func New(staticDir string, connOpt redis.UniversalOptions, launcherURL, previewURL, authKey string) (*Server, error) {
	s := &Server{
		staticDir:       staticDir,
		work:            make(chan string),
		l:               launcher.MustNewManaged(launcherURL),
		Mux:             chi.NewRouter(),
		previewURL:      previewURL,
		generateTimeout: 90 * time.Second,
		cacheTTL:        15 * time.Minute,
		authKey:         authKey,
	}
	s.rdb = redis.NewUniversalClient(&connOpt)
	_, err := s.rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	s.browser = rod.New().Client(s.l.MustClient())
	err = s.browser.Connect()
	if err != nil {
		return nil, fmt.Errorf("error connecting to browser: %w", err)
	}

	return s, nil
}

func (s *Server) SetOpts(opts ...serverCfg) error {
	for _, f := range opts {
		err := f(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Init() error {
	if s.logger == nil {
		s.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	err := s.routes()
	if err != nil {
		return err
	}
	go s.listen()

	return nil
}

func (s *Server) Shutdown() error {
	return s.browser.Close()
}

func WithLogger(logger *slog.Logger) serverCfg {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

func WithLocalAssets(prefix, dir string) serverCfg {
	return func(s *Server) error {
		s.useLocalAssets = true
		s.assetsPrefix = prefix
		s.assetsDir = dir
		return nil
	}
}

func WithProxy(prefix, target string) serverCfg {
	return func(s *Server) error {
		s.useProxy = true
		s.proxyPrefix = prefix
		host, err := url.Parse(target)
		if err != nil {
			return fmt.Errorf("error parsing target %v: %w", target, err)
		}
		s.proxyTarget = host
		return nil
	}
}

func WithSkipTLSVerify() serverCfg {
	return func(s *Server) error {
		s.skipInsecure = true
		return nil
	}
}

func WithCacheTTL(ttl int) serverCfg {
	return func(s *Server) error {
		if ttl <= 0 {
			return fmt.Errorf("invalid cache ttl <= 0: %v", ttl)
		}
		s.cacheTTL = time.Duration(ttl * int(time.Second))
		return nil
	}
}

func WithGenerateTimeout(timeout int) serverCfg {
	return func(s *Server) error {
		if timeout <= 0 {
			return fmt.Errorf("invalid timeout <= 0: %v", timeout)
		}
		s.generateTimeout = time.Duration(timeout * int(time.Second))
		return nil
	}
}

func (s *Server) routes() error {
	s.Use(middleware.Logger)
	s.Use(middleware.Recoverer)
	s.Use(middleware.RequestID)
	s.With(s.authKeyCheck).Route("/", func(r chi.Router) {
		if s.useLocalAssets {
			localAssetsFS := http.FileServer(http.Dir(s.assetsDir))
			r.Handle(fmt.Sprintf("%v/*", s.assetsPrefix), http.StripPrefix(s.assetsPrefix+"/", localAssetsFS))
		}

		if s.useProxy {
			path := strings.TrimSuffix(s.proxyPrefix, "/")
			r.Handle(path+"/*", s.handleProxy(path))

			s.proxy = httputil.NewSingleHostReverseProxy(s.proxyTarget)
			if s.skipInsecure {
				s.proxy.Transport = &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
			}
		}

		r.Handle("/generate/db/{id}", s.handleImageRequest("db"))
		r.Handle("/generate/sh/{id}", s.handleImageRequest("sh"))

		r.NotFound(s.notFoundHandler())
	})

	return nil
}
