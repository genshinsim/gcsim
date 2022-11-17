package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/genshinsim/gcsim/backend/pkg/services/preview"
	"github.com/genshinsim/gcsim/backend/pkg/services/result"
)

//go:embed dist/*
var content embed.FS

func main() {

	resultStore, err := result.NewClient(result.ClientCfg{
		Addr: os.Getenv("RESULT_STORE_URL"),
	})

	if err != nil {
		panic(err)
	}

	s, err := preview.New(preview.Config{
		URL:          "http://localhost:3000",
		Files:        content,
		AssetsFolder: os.Getenv(("ASSETS_DATA_PATH")),
		ResultStore:  resultStore,
	})

	if err != nil {
		panic(err)
	}

	log.Println("starting img generation listener")
	log.Fatal(http.ListenAndServe(":3000", s.Router))
}
