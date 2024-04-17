package embedgenerator

import (
	context "context"
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

type ServerCfg func(s *Server) error

type Server struct {
	*chi.Mux

	// server context for graceful shutdown
	serverClosed chan bool

	// status for if service is ready
	redisReady bool
	rodReady   bool

	logger      *slog.Logger
	rdb         redis.UniversalClient
	work        chan string
	l           *launcher.Launcher
	launcherURL string
	previewURL  string
	authKey     string

	// static asset; mandatory to serve from
	staticDir string

	// proxy for assets
	useAssetsProxy    bool
	assetsProxyPrefix string // cannot be blank
	assetsProxyTarget *url.URL
	assetsProxy       *httputil.ReverseProxy

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
		launcherURL:     launcherURL,
		Mux:             chi.NewRouter(),
		previewURL:      previewURL,
		generateTimeout: 90 * time.Second,
		cacheTTL:        15 * time.Minute,
		authKey:         authKey,
		serverClosed:    make(chan bool),
		rdb:             redis.NewUniversalClient(&connOpt),
	}

	go s.handleRedis()
	go s.handleLauncher()

	return s, nil
}

func (s *Server) SetOpts(opts ...ServerCfg) error {
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
	close(s.serverClosed)
	return s.browser.Close()
}

func WithLogger(logger *slog.Logger) ServerCfg {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

func WithAssetsProxy(prefix, target string) ServerCfg {
	return func(s *Server) error {
		s.useAssetsProxy = true
		s.assetsProxyPrefix = prefix
		host, err := url.Parse(target)
		if err != nil {
			return fmt.Errorf("error parsing target %v: %w", target, err)
		}
		s.assetsProxyTarget = host
		return nil
	}
}

func WithProxy(prefix, target string) ServerCfg {
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

func WithSkipTLSVerify() ServerCfg {
	return func(s *Server) error {
		s.skipInsecure = true
		return nil
	}
}

func WithCacheTTL(ttl int) ServerCfg {
	return func(s *Server) error {
		if ttl <= 0 {
			return fmt.Errorf("invalid cache ttl <= 0: %v", ttl)
		}
		s.cacheTTL = time.Duration(ttl * int(time.Second))
		return nil
	}
}

func WithGenerateTimeout(timeout int) ServerCfg {
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
	s.Handle("/online", s.handleOnlineCheck())
	s.With(s.authKeyCheck).With(s.readyCheck).Route("/", func(r chi.Router) {
		if s.useAssetsProxy {
			path := strings.TrimSuffix(s.assetsProxyPrefix, "/")
			r.Handle(path+"/*", s.handleProxy(path))

			s.assetsProxy = httputil.NewSingleHostReverseProxy(s.assetsProxyTarget)
			if s.skipInsecure {
				s.assetsProxy.Transport = &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}
			}
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

func (s *Server) handleRedis() {
	_, err := s.rdb.Ping(context.Background()).Result()
	if err != nil {
		s.logger.Warn("redis ping failed", "err", err)
	}
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case _, ok := <-s.serverClosed:
			if !ok {
				return
			}
		case <-ticker.C:
			_, err := s.rdb.Ping(context.Background()).Result()
			if err != nil {
				s.logger.Warn("redis ping failed", "err", err)
			}
			s.redisReady = err == nil
		}
	}
}

func (s *Server) handleLauncher() {
	// initial connection
	err := s.tryConnectLauncher()
	if err != nil {
		s.logger.Warn("launcher connect failed", "err", err)
	}
	ticker := time.NewTicker(30 * time.Second)
	// keep trying to connect to launcher and establish a browser if not available
	for {
		select {
		case _, ok := <-s.serverClosed:
			if !ok {
				return
			}
		case <-ticker.C:
			// try connecting if not already connected
			err := s.tryConnectLauncher()
			if err != nil {
				s.logger.Warn("launcher connect failed", "err", err)
			}
			s.rodReady = err == nil
		}
	}
}

func (s *Server) tryConnectLauncher() error {
	if !s.rodReady {
		l, err := launcher.NewManaged(s.launcherURL)
		if err != nil {
			return fmt.Errorf("error instancing launcher: %w", err)
		}
		s.l = l
		// try connecting browser
		rc, err := l.Client()
		if err != nil {
			return fmt.Errorf("error creating client for remote browser: %w", err)
		}
		s.browser = rod.New().Client(rc)
	}
	// try checking browser version
	version, err := s.browser.Version()
	if err != nil {
		return fmt.Errorf("browser version check failed: %w", err)
	}
	s.logger.Info("browser version check ok", "version", version)
	return nil
}
