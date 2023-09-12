package main

import (
	"log"
	"os"
	"regexp"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/genshinsim/gcsim/backend/pkg/discord"
	"github.com/genshinsim/gcsim/backend/pkg/discord/backend"
	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"github.com/genshinsim/gcsim/pkg/model"
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
		LinkValidationRegex: regexp.MustCompile(`https://\S+.app/\S+/(\S+)$`),
		ShareStore:          makeShareStore(),
		DBgRPCAddr:          os.Getenv("DB_STORE_URL"),
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
		//TODO: consider moving this mapping to models maybe?
		TagMapping: map[string]model.DBTag{
			"1080228340427927593": model.DBTag_DB_TAG_GCSIM,
			"1118916799153582170": model.DBTag_DB_TAG_TESTING,
			"1120875165346177024": model.DBTag_DB_TAG_KQM_GUIDE,
			"1120878673952788500": model.DBTag_DB_TAG_GEO_SIMPS,
			"1120878739786571866": model.DBTag_DB_TAG_ITTO_SIMPS,
			"1148370191185616948": model.DBTag_DB_TAG_RANDOM_DELAYS,
		},
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
