package main

import (
	"log"
	"os"

	"github.com/genshinsim/gcsim/services/pkg/bot"
)

func main() {
	token := os.Getenv("DBMANAGE_DISCORD_TOKEN")
	path := os.Getenv("DBMANAGE_DB_PATH")
	adminChan := os.Getenv("DBMANAGE_ADMIN_CHAN")
	port := os.Getenv("DBMANAGE_PG_PORT")

	cfg := bot.Config{
		Token:          token,
		DBPath:         path,
		AdminChannelID: adminChan,
		PostgRESTPort:  port,
	}

	err := bot.Run(cfg)

	if err != nil {
		log.Fatal(err)
	}

}
