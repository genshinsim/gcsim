package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/genshinsim/gcsim/internal/services/assets"
)

type config struct {
	Host string `env:"HOST"`
	Port string `env:"PORT" envDefault:"3000"`
	// auth key for checking incoming requests
	AssetsPrefix string `env:"ASSETS_PREFIX"`
	CacheDir     string `env:"CACHE_DIR"`
	// timeouts
	ExternalGetTimeoutInMS int `env:"EXTERNAL_GET_TIMEOUT_IN_MS"`
	// sources
	SourceType []string `env:"SOURCE_TYPE" envSeparator:""`
	SourceHost []string `env:"SOURCE_HOST" envSeparator:""`
}

func main() {
	var err error

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)

	if len(cfg.SourceType) != len(cfg.SourceHost) {
		log.Fatal("unmatch source type and source hosts; not equal in number of entries")
	}

	// append mhy defaults
	cfg.SourceType = append(cfg.SourceType, "avatar", "weapons", "artifacts")
	cfg.SourceHost = append(cfg.SourceHost,
		"https://upload-os-bbs.mihoyo.com/game_record/genshin/character_icon/",
		"https://upload-os-bbs.mihoyo.com/game_record/genshin/equip/",
		"https://upload-os-bbs.mihoyo.com/game_record/genshin/equip/",
	)

	log.Println("running with config ", cfg)

	server, err := assets.New()

	panicErr(err)

	for i, v := range cfg.SourceType {
		t, err := assets.AssetTypeFromString(v)
		panicErr(err)
		err = server.SetOpts(assets.WithAssetSource(t, cfg.SourceHost[i]))
		panicErr(err)
	}

	if cfg.ExternalGetTimeoutInMS > 0 {
		panicErr(server.SetOpts(assets.WithCustomTimeout(time.Millisecond * time.Duration(cfg.ExternalGetTimeoutInMS))))
	}
	if cfg.AssetsPrefix != "" {
		panicErr(server.SetOpts(assets.WithAssetsPrefix(cfg.AssetsPrefix)))
	}
	if cfg.CacheDir != "" {
		panicErr(server.SetOpts(assets.WithCacheDir(cfg.CacheDir)))
	}

	err = server.Init()
	if err != nil {
		log.Fatal(err)
	}

	httpServer := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: server,
	}

	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	defer func() {
		log.Println("Shutting down server: ", server.Shutdown())
	}()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Println("HTTP graceful shutdown encountered error, forcing shutdown")
		// force shut down
		err := httpServer.Close()
		log.Println("Force shut down completed with error: ", err)
		log.Panicf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
