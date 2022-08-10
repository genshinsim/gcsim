package main

import (
	"log"
	"net/http"

	"github.com/genshinsim/gcsim/services/embed/pkg/server"
)

func main() {
	s, err := server.New(server.Config{
		PostgRESTPort: 3000,
	})
	if err != nil {
		panic(err)
	}
	log.Println("Starting to listen at port 3001")
	log.Fatal(http.ListenAndServe(":3001", s.Router))
}
