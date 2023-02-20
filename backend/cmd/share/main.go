package main

import (
	"log"
	"net"

	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"github.com/genshinsim/gcsim/backend/pkg/services/share/mock"
	"google.golang.org/grpc"
)

func main() {
	shareStore, err := mock.NewServer()
	if err != nil {
		panic(err)
	}

	server, err := share.New(share.Config{
		Store: shareStore,
	})
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	share.RegisterShareStoreServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
