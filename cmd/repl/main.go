package main

import (
	"log"
	"os"

	"github.com/genshinsim/gcsim/pkg/gcs/repl"
)

func main() {
	l := log.New(os.Stdout, "eval", log.LstdFlags)
	for {
		runReplCatchPanic(l)
	}
}

func runReplCatchPanic(l *log.Logger) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	repl.Start(os.Stdin, os.Stdout, l, false)
}
