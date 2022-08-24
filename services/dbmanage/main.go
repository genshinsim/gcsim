package main

import (
	"log"
	"os"

	"github.com/genshinsim/gcsim/services/pkg/db"
	"github.com/genshinsim/gcsim/services/pkg/store"
)

func main() {

	cfg := db.Config{
		Token:          os.Getenv("DBMANAGE_DISCORD_TOKEN"),
		DBPath:         os.Getenv("DATA_PATH"),
		AdminChannelID: os.Getenv("DBMANAGE_ADMIN_CHAN"),
	}

	err := db.Run(cfg, store.NewPostgRESTStore(os.Getenv("POSTGREST_URL")))

	if err != nil {
		log.Fatal(err)
	}

}
