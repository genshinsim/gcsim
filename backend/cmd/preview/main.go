package main

import (
	"embed"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/genshinsim/gcsim/backend/pkg/services/preview"
	"google.golang.org/grpc"
)

//go:embed dist/*
var content embed.FS

func main() {

	server, err := preview.New(preview.Config{
		URL:          "http://localhost:3001",
		Files:        content,
		AssetsFolder: os.Getenv(("ASSETS_DATA_PATH")),
	})

	if err != nil {
		panic(err)
	}

	log.Println("starting img generation listener")
	go log.Fatal(http.ListenAndServe(":3001", server.Router))

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	preview.RegisterEmbedServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
