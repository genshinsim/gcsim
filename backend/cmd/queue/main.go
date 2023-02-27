package main

import (
	"log"
	"net"
	"time"

	"github.com/genshinsim/gcsim/backend/pkg/services/queue"
	"google.golang.org/grpc"
)

func main() {

	server, err := queue.NewQueue(queue.Config{
		Timeout: 1 * time.Minute,
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
