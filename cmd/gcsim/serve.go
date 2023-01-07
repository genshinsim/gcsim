package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	hash, err := hashFromCompressed(compressed)
	if err != nil {
		log.Printf("error generating secure has from gz data: %v\n", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return false
	}

	log.Println("Received results request, sending response...")
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Encoding", "deflate")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Expose-Headers", "X-GCSIM-SHARE-AUTH")
	resp.Header().Set("X-GCSIM-SHARE-AUTH", string(hash))
	resp.WriteHeader(http.StatusOK)
	resp.Write(compressed)

	if f, ok := resp.(http.Flusher); ok {
		f.Flush()
	}
	return true
}

func hashFromCompressed(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	json.Unmarshal(b, &res)
	b, _ = json.Marshal(res)

	h := sha256.New()
	h.Write(b)
	bs := h.Sum(nil)

	//shareKey should be of the format id:key
	id, hexkey, ok := strings.Cut(shareKey, ":")
	if !ok {
		return nil, fmt.Errorf("invalid share key")
	}
	key, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, fmt.Errorf("invalid share key")
	}

	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	hash := gcm.Seal(nonce, nonce, bs, nil)
	hashStr := base64.StdEncoding.EncodeToString(hash)

	return []byte(id + ":" + hashStr), nil
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
