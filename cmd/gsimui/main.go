// Package gsim does something
package main

import (
	"embed"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/genshinsim/gcsim/pkg/server"
)

var openWin bool

//go:embed build/*
var content embed.FS

func main() {

	s, err := server.New(
		func(s *server.Server) error {
			s.Cfg.ConfigDir = "./config"
			s.Cfg.Port = 8081
			return nil
		},

		//add routes
		func(s *server.Server) error {
			s.Router.NotFound(NotFoundHandler)
			return nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if openWin {
		go checkOpen()
	}

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8081", s.Router))
}

var ErrDir = errors.New("path is dir")

func tryRead(fs embed.FS, prefix, requestedPath string, w http.ResponseWriter) error {
	f, err := fs.Open(path.Join(prefix, requestedPath))
	if err != nil {
		return err
	}
	defer f.Close()

	stat, _ := f.Stat()
	if stat.IsDir() {
		return ErrDir
	}

	contentType := mime.TypeByExtension(filepath.Ext(requestedPath))
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, f)
	return err
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := tryRead(content, "build", r.URL.Path, w)
	if err == nil {
		return
	}
	err = tryRead(content, "build", "index.html", w)
	if err != nil {
		panic(err)
	}
}

func checkOpen() {
	for {
		time.Sleep(time.Second)

		log.Println("Checking if started...")
		resp, err := http.Get("http://localhost:8081")
		if err != nil {
			log.Println("Failed:", err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Println("Not OK:", resp.StatusCode)
			continue
		}

		// Reached this point: server is up and running!
		break
	}
	log.Println("SERVER UP AND RUNNING!")
	open("http://localhost:8081")
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
