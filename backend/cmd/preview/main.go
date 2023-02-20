package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/genshinsim/gcsim/backend/pkg/services/preview"
	"github.com/genshinsim/gcsim/backend/pkg/services/share"
)

//go:embed dist/*
var content embed.FS

func main() {

	shareStore, err := share.NewClient(share.ClientCfg{
		Addr: os.Getenv("SHARE_STORE_URL"),
	})

	if err != nil {
		panic(err)
	}

	s, err := preview.New(preview.Config{
		URL:          "http://localhost:3000",
		Files:        content,
		AssetsFolder: os.Getenv(("ASSETS_DATA_PATH")),
		ShareStore:   shareStore,
	})

	if err != nil {
		panic(err)
	}

	log.Println("starting img generation listener")
	log.Fatal(http.ListenAndServe(":3000", s.Router))
}
