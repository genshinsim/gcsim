package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type viewerResults struct {
	Data string `json:"data"`
}

func serve(
	connectionsClosed chan struct{},
	resultPath string,
	hash string,
	samplePath string,
	keepAlive bool) {

	server := &http.Server{Addr: address}
	done := make(chan bool)

	http.HandleFunc("/data", func(resp http.ResponseWriter, req *http.Request) {
		success := handleResults(resp, req, resultPath, hash)
		if success && !keepAlive {
			done <- true
		}
	})

	http.HandleFunc("/sample", func(resp http.ResponseWriter, req *http.Request) {
		success := handleSample(resp, req, samplePath)
		if success && !keepAlive {
			done <- true
		}
	})

	go interruptShutdown(server, done, connectionsClosed)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndSever Error: %v", err)
		}
	}()
}

func handleResults(resp http.ResponseWriter, req *http.Request, path string, hash string) bool {
	if req.Method == "OPTIONS" {
		log.Println("OPTIONS request received, responding...")
		optionsResponse(resp)
		return false
	}

	if req.Method != "GET" {
		log.Printf("Invalid request method: %v\n", req.Method)
		resp.WriteHeader(http.StatusForbidden)
		return false
	}

	compressed, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error reading gz data: %v\n", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return false
	}

	log.Println("Received results request, sending response...")
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Encoding", "deflate")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Expose-Headers", "X-GCSIM-SHARE-AUTH")
	resp.Header().Set("X-GCSIM-SHARE-AUTH", hash)
	resp.WriteHeader(http.StatusOK)
	resp.Write(compressed)

	if f, ok := resp.(http.Flusher); ok {
		f.Flush()
	}
	return true
}

func handleSample(resp http.ResponseWriter, req *http.Request, path string) bool {
	if req.Method == "OPTIONS" {
		log.Println("OPTIONS request received, responding...")
		optionsResponse(resp)
		return false
	}

	if req.Method != "GET" {
		log.Printf("Invalid request method: %v\n", req.Method)
		resp.WriteHeader(http.StatusForbidden)
		return false
	}

	compressed, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error reading gz data: %v\n", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return false
	}

	log.Println("Received sample request, sending response...")
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Encoding", "deflate")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.WriteHeader(http.StatusOK)
	resp.Write(compressed)

	if f, ok := resp.(http.Flusher); ok {
		f.Flush()
	}
	return true
}

func optionsResponse(resp http.ResponseWriter) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	resp.Header().Set(
		"Access-Control-Allow-Headers",
		"Accept, Access-Control-Allow-Origin, Content-Type, "+
			"Content-Length, Accept-Encoding, Authorization")
	resp.WriteHeader(http.StatusNoContent)
}

func interruptShutdown(server *http.Server, done chan bool, connectionsClosed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	select {
	case <-done:
	case <-sigint:
	}

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown Error: %v", err)
	}
	close(connectionsClosed)
}
