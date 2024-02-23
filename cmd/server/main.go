package main

import (
	"log"
	"net/http"
	"os"

	"github.com/genshinsim/gcsim/pkg/servermode"
)

var (
	shareKey string
)

func main() {
	if shareKey == "" {
		shareKey = os.Getenv("GCSIM_SHARE_KEY")
	}

	server, err := servermode.New(
		servermode.WithDefaults(),
		servermode.WithShareKey(shareKey),
	)

	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe("localhost:54321", server.Router))
}
