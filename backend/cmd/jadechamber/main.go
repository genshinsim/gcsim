package main

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sha1ver string

func main() {
	setHash()

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	sugar.Debugw("jadechamber started", "sha1ver", sha1ver)

	s, err := api.New(api.Config{
		ShareStore:  makeShareStore(),
		DBShareKeys: makeDBShareKeys(),
	}, func(s *api.Server) error {
		s.Log = sugar
		return nil
	})
	if err != nil {
		panic(err)
	}

	log.Println("API gateway starting to listen at port 3000")
	log.Fatal(http.ListenAndServe(":3000", s.Router))
}

func setHash() {
	info, _ := debug.ReadBuildInfo()
	for _, bs := range info.Settings {
		if bs.Key == "vcs.revision" {
			sha1ver = bs.Value
		}
	}
}

func makeShareStore() api.ShareStore {
	shareStore, err := share.NewClient(share.ClientCfg{
		Addr: os.Getenv("SHARE_STORE_URL"),
	})
	if err != nil {
		panic(err)
	}
	return shareStore
}
