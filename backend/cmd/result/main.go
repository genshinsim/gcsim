package main

import (
	"log"
	"net"
	"os"

	"github.com/genshinsim/gcsim/backend/pkg/services/result"
	"google.golang.org/grpc"
)

func main() {

	store, err := result.New(result.Config{
		DBPath: os.Getenv("RESULT_DATA_PATH"),
	})

	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	result.RegisterResultStoreServer(s, store)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
