package main

import (
	"log"
	"net/http"
	"os"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/genshinsim/gcsim/backend/pkg/result"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	sugar.Debugw("logger initiated")

	resultStore, err := result.New(result.Config{
		DBPath: os.Getenv("DATA_PATH"),
	}, func(s *result.Store) error {
		s.Log = sugar
		return nil
	})

	if err != nil {
		panic(err)
	}

	s, err := api.New(api.Config{
		ResultStore: resultStore,
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
