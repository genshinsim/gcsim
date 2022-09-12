package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	folder := os.Getenv("ASSETS_PATH")
	fs := http.FileServer(http.Dir(folder))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Print("listening on :3000 ...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
