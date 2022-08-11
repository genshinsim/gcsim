package main

import (
	"log"
	"os"

	"github.com/genshinsim/gcsim/services/pkg/bot"
	"github.com/genshinsim/gcsim/services/pkg/store"
)

func main() {

	cfg := bot.Config{
		Token:          os.Getenv("DBMANAGE_DISCORD_TOKEN"),
		DBPath:         os.Getenv("DATA_PATH"),
		AdminChannelID: os.Getenv("DBMANAGE_ADMIN_CHAN"),
	}

	err := bot.Run(cfg, &store.PostgRESTStore{
		URL: os.Getenv("POSTGREST_URL"),
	})

	if err != nil {
		log.Fatal(err)
	}

}
