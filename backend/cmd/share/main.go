package main

import (
	"log"
	"net"
	"os"

	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"github.com/genshinsim/gcsim/backend/pkg/services/share/mongo"
	"google.golang.org/grpc"
)

func main() {
	mongoCfg := mongo.Config{
		URL:        os.Getenv("MONGODB_URL"),
		Database:   os.Getenv("MONGODB_DATABASE"),
		Collection: os.Getenv("MONGODB_COLLECTION"),
		Username:   os.Getenv("MONGODB_USERNAME"),
		Password:   os.Getenv("MONOGDB_PASSWORD"),
	}
	log.Println(os.Getenv("MONGODB_URL"))
	log.Printf("Cfg: %v\n", mongoCfg)
	shareStore, err := mongo.NewServer(mongoCfg)
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
