package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type viewerResults struct {
	Data string `json:"data"`
}

func serve(connectionsClosed chan struct{}, resultPath string, samplePath string, keepAlive bool) {
	server := &http.Server{Addr: address}

	http.HandleFunc("/data", func(resp http.ResponseWriter, req *http.Request) {
		success := handleResults(resp, req, resultPath)
		if success && !keepAlive {
			shutdown()
		}
	})

	http.HandleFunc("/sample", func(resp http.ResponseWriter, req *http.Request) {
		success := handleSample(resp, req, samplePath)
		if success && !keepAlive {
			shutdown()
		}
	})

	go interruptShutdown(server, connectionsClosed)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndSever Error: %v", err)
		}
	}()
}

func handleResults(resp http.ResponseWriter, req *http.Request, path string) bool {
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

	encoded := base64.StdEncoding.EncodeToString(compressed)
	results, err := json.Marshal(viewerResults{Data: encoded})
	if err != nil {
		log.Printf("error marshal json: %v\n", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return false
	}

	log.Println("Received results request, sending response...")
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.WriteHeader(http.StatusOK)
	resp.Write(results)

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
	encoded := base64.StdEncoding.EncodeToString(compressed)

	log.Println("Received sample request, sending response...")
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.WriteHeader(http.StatusOK)
	io.WriteString(resp, encoded)

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

func interruptShutdown(server *http.Server, connectionsClosed chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown Error: %v", err)
	}
	close(connectionsClosed)
}

func shutdown() {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Signal(os.Interrupt); err != nil {
		log.Fatal(err)
	}
}
