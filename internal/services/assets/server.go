package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//go:embed static/*
var staticAssets embed.FS

type AssetType int

const (
	AssetTypeInvalid AssetType = iota
	AssetTypeAny
	AssetTypeAvatars
	AssetTypeWeapons
	AssetTypeArtifacts
	EndAssetType
)

var assetTypeStr = []string{"", "*", "avatar", "weapons", "artifacts"}

func (a AssetType) String() string {
	return assetTypeStr[a]
}

func AssetTypeFromString(s string) (AssetType, error) {
	for i, v := range assetTypeStr {
		if s == v {
			return AssetType(i), nil
		}
	}
	return AssetTypeInvalid, nil
}

type ServerCfg func(s *Server) error

type Server struct {
	*chi.Mux

	hosts        map[AssetType][]*url.URL // external source for images
	assetsPrefix string                   // assets prefix
	cacheDir     string                   // where to save cache files

	assetNameMapping []map[string]string
	httpClient       *http.Client
	logger           *slog.Logger
}

func New() (*Server, error) {
	s := &Server{
		hosts:            make(map[AssetType][]*url.URL),
		assetsPrefix:     "/api/assets",
		cacheDir:         "/cache",
		assetNameMapping: make([]map[string]string, EndAssetType),
		httpClient: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
		Mux: chi.NewRouter(),
	}
	s.assetNameMapping[AssetTypeAvatars] = avatarMap
	s.assetNameMapping[AssetTypeWeapons] = weaponMap
	s.assetNameMapping[AssetTypeArtifacts] = artfactMap

	return s, nil
}

func (s *Server) Init() error {
	if s.logger == nil {
		s.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	err := s.routes()
	if err != nil {
		return err
	}
	// make sure cache dirs exists
	s.logger.Info("creating cache dirs...")
	os.MkdirAll(path.Join(s.cacheDir, "/"+AssetTypeAvatars.String()), fs.ModePerm)
	os.MkdirAll(path.Join(s.cacheDir, "/"+AssetTypeWeapons.String()), fs.ModePerm)
	os.MkdirAll(path.Join(s.cacheDir, "/"+AssetTypeArtifacts.String()), fs.ModePerm)

	return nil
}

func (s *Server) Shutdown() error {
	s.logger.Info("shutting down...")
	return nil
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

func WithAssetsPrefix(prefix string) ServerCfg {
	return func(s *Server) error {
		s.assetsPrefix = prefix
		return nil
	}
}

func WithCacheDir(cacheDir string) ServerCfg {
	return func(s *Server) error {
		s.cacheDir = cacheDir
		return nil
	}
}

func WithCustomTimeout(timeout time.Duration) ServerCfg {
	return func(s *Server) error {
		s.httpClient.Timeout = timeout
		return nil
	}
}

func WithAssetSource(assetType AssetType, host string) ServerCfg {
	return func(s *Server) error {
		u, err := url.Parse(host)
		if err != nil {
			return err
		}
		switch assetType {
		case AssetTypeAny:
			s.hosts[AssetTypeAvatars] = append(s.hosts[AssetTypeAvatars], u)
			s.hosts[AssetTypeWeapons] = append(s.hosts[AssetTypeWeapons], u)
			s.hosts[AssetTypeArtifacts] = append(s.hosts[AssetTypeArtifacts], u)
			return nil
		case AssetTypeAvatars:
		case AssetTypeWeapons:
		case AssetTypeArtifacts:
		default:
			return fmt.Errorf("unrecognized asset type: %v", assetType)
		}
		s.hosts[assetType] = append(s.hosts[assetType], u)
		return nil
	}
}

type myFS struct {
	content embed.FS
}

func (c myFS) Open(name string) (fs.File, error) {
	return c.content.Open(path.Join("static", name))
}

func (s *Server) routes() error {
	s.Use(middleware.Logger)
	s.Use(middleware.Recoverer)
	s.Use(middleware.RequestID)

	path := strings.TrimSuffix(s.assetsPrefix, "/")
	s.logger.Info("adding assets route", "path", path)
	s.Handle(fmt.Sprintf("%v/%v/{key}.png", path, AssetTypeAvatars.String()), s.handleGetData(AssetTypeAvatars))
	s.Handle(fmt.Sprintf("%v/%v/{key}.png", path, AssetTypeWeapons.String()), s.handleGetData(AssetTypeWeapons))
	s.Handle(fmt.Sprintf("%v/%v/{key}.png", path, AssetTypeArtifacts.String()), s.handleGetData(AssetTypeArtifacts))

	fs := http.FileServer(http.FS(&myFS{content: staticAssets}))
	s.Handle("/api/assets/*", http.StripPrefix("/api/assets/", fs))

	s.Handle("/online", s.handleOnlineCheck())

	s.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return nil
}
