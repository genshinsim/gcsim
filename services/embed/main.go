package main

import (
	"log"
	"net/http"
	"os"

	"github.com/genshinsim/gcsim/services/pkg/embed"
	"github.com/genshinsim/gcsim/services/pkg/pygenerator"
	"github.com/genshinsim/gcsim/services/pkg/store"
)

func main() {
	pgStore := &store.PostgRESTStore{URL: os.Getenv("POSTGREST_URL")}

	s, err := embed.New(embed.Config{
		DataFolder: os.Getenv("DATA_PATH"),
	}, func(s *embed.Server) error {
		s.Store = pgStore
		s.Generator = pygenerator.New("./scripts/embed.py")
		return nil
	})

	if err != nil {
		panic(err)
	}
	log.Println("Starting to listen at port 3001")
	log.Fatal(http.ListenAndServe(":3001", s.Router))
}
