package main

import (
	"log"
	"os"
	"regexp"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/genshinsim/gcsim/backend/pkg/discord"
	"github.com/genshinsim/gcsim/backend/pkg/discord/backend"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"github.com/genshinsim/gcsim/backend/pkg/services/submission"
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

	store, err := backend.New(backend.Config{
		LinkValidationRegex: regexp.MustCompile(`https:\/\/\S+.app\/\S+\/(\S+)$`),
		ShareStore:          makeShareStore(),
		SubmissionStore:     makeSubStore(),
	}, func(s *backend.Store) error {
		s.Log = sugar
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	b, err := discord.New(discord.Config{
		Token:   os.Getenv("DISCORD_BOT_TOKEN"),
		Backend: store,
	}, func(b *discord.Bot) error {
		b.Log = sugar
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = b.Run()
	if err != nil {
		log.Fatal(err)
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

func makeDBStore() api.DBStore {
	store, err := db.NewClient(db.ClientCfg{
		Addr: os.Getenv("DB_STORE_URL"),
	})
	if err != nil {
		panic(err)
	}
	return store
}

func makeSubStore() api.SubmissionStore {
	store, err := submission.NewClient(os.Getenv("SUBMISSION_STORE_URL"))
	if err != nil {
		panic(err)
	}
	return store
}
