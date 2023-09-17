package main

import (
	"log"
	"os"

	"github.com/genshinsim/gcsim/pkg/gcs/repl"
)

func main() {
	f := ""
	if len(os.Args) > 1 {
		f = os.Args[1]
	}
	l := log.New(os.Stdout, "eval log ", log.LstdFlags)

	if f != "" {
		b, err := os.ReadFile("./" + f)
		if err != nil {
			panic(err)
		}
		repl.Eval(string(b), l)
		return
	}

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
