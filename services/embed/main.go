package main

import (
	"log"
	"net/http"

	"github.com/genshinsim/gcsim/services/pkg/embed"
)

func main() {
	s, err := embed.New(embed.Config{
		PostgRESTPort: 3000,
	})
	if err != nil {
		panic(err)
	}
	log.Println("Starting to listen at port 3001")
	log.Fatal(http.ListenAndServe(":3001", s.Router))
}
