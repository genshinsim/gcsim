package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/genshinsim/gcsim/services/pkg/embed"
	"github.com/genshinsim/gcsim/services/pkg/store"
)

func main() {
	//spin up a test server to test embeds
	//curl following:
	//curl http://localhost:3001/embed/test -X POST
	//curl http://localhost:3001/embed/test -o result.png

	s, err := embed.New(embed.Config{
		DataFolder: "./db",
	}, func(s *embed.Server) error {
		s.Store = &db{}
		s.Generator = &generator{}
		return nil
	})

	if err != nil {
		panic(err)
	}
	log.Println("Starting to listen at port 3001")
	log.Fatal(http.ListenAndServe(":3001", s.Router))
}

type db struct {
}

func (d *db) Fetch(url string) (store.Simulation, error) {
	return store.Simulation{}, nil
}

type generator struct {
}

func (g *generator) Generate(sim store.Simulation, filepath string) error {

	source, err := os.Open("./source.png")
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer destination.Close()
	//make copy of result.pgn
	_, err = io.Copy(destination, source)

	return err
}
