package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/backend/pkg/services/queue"
	"github.com/genshinsim/gcsim/backend/pkg/services/submission"
	"google.golang.org/grpc"
)

func main() {
	subStore, err := submission.NewClient(os.Getenv("SUBMISSION_STORE_URL"))
	if err != nil {
		panic(err)
	}

	dbStore, err := db.NewClient(db.ClientCfg{
		Addr: os.Getenv("DB_STORE_URL"),
	})
	if err != nil {
		panic(err)
	}

	server, err := queue.NewQueue(queue.Config{
		Timeout: 1 * time.Minute,
		DBWork:  dbStore,
		SubWork: subStore,
	})

	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen on port 3000: %v", err)
	}

	s := grpc.NewServer()
	queue.RegisterWorkQueueServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
