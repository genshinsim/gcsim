package main

import (
	"log"
	"net"
	"os"
	"runtime/debug"

	"github.com/genshinsim/gcsim/backend/pkg/mongo"
	"github.com/genshinsim/gcsim/backend/pkg/notify"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"google.golang.org/grpc"
)

var (
	sha1ver string
)

func main() {
	info, _ := debug.ReadBuildInfo()
	for _, bs := range info.Settings {
		if bs.Key == "vcs.revision" {
			sha1ver = bs.Value
		}
	}
	mongoCfg := mongo.Config{
		URL:         os.Getenv("MONGODB_URL"),
		Database:    os.Getenv("MONGODB_DATABASE"),
		Collection:  os.Getenv("MONGODB_COLLECTION"),
		ValidView:   os.Getenv("MONGODB_QUERY_VIEW"),
		SubView:     os.Getenv("MONGODB_SUB_VIEW"),
		Username:    os.Getenv("MONGODB_USERNAME"),
		Password:    os.Getenv("MONOGDB_PASSWORD"),
		CurrentHash: sha1ver,
	}
	log.Println(os.Getenv("MONGODB_URL"))
	log.Printf("Cfg: %v\n", mongoCfg)
	log.Printf("Current hash: %v\n", sha1ver)
	dbStore, err := mongo.NewServer(mongoCfg)
	if err != nil {
		panic(err)
	}
	shareStore, err := share.NewClient(share.ClientCfg{
		Addr: os.Getenv("SHARE_STORE_URL"),
	})

	if err != nil {
		panic(err)
	}

	n, err := notify.New("db-notifier")
	if err != nil {
		panic(err)
	}

	server, err := db.NewServer(db.Config{
		DBStore:       dbStore,
		ShareStore:    shareStore,
		ExpectedHash:  sha1ver,
		NotifyService: n,
	})

	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	db.RegisterDBStoreServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
