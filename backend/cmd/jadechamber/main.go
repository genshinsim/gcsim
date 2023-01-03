package main

import (
	"log"
	"net/http"
	"os"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/backend/pkg/services/result"
	"github.com/genshinsim/gcsim/backend/pkg/user"
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

	resultStore, err := result.NewClient(result.ClientCfg{
		Addr: os.Getenv("RESULT_STORE_URL"),
	})

	if err != nil {
		panic(err)
	}

	userStore, err := user.New(user.Config{
		DBPath: os.Getenv("USER_DATA_PATH"),
	}, func(s *user.Store) error {
		s.Log = sugar
		return nil
	})

	if err != nil {
		panic(err)
	}

	dbStore, err := db.NewClient(db.ClientCfg{
		Addr: os.Getenv("DB_STORE_URL"),
	})
	if err != nil {
		panic(err)
	}

	s, err := api.New(api.Config{
		ResultStore: resultStore,
		UserStore:   userStore,
		DBStore:     dbStore,
		Discord: api.DiscordConfig{
			RedirectURL:  os.Getenv("REDIRECT_URL"),
			ClientID:     os.Getenv("DISCORD_ID"),
			ClientSecret: os.Getenv("DISCORD_SECRET"),
			JWTKey:       os.Getenv("JWT_KEY"),
		},
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
