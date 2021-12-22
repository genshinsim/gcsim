package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/genshinsim/gcsim"
	"github.com/genshinsim/gcsim/internal/logtohtml"
	"github.com/genshinsim/gcsim/pkg/calcqueue"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
)

func main() {
	cfgFile := flag.String("c", "config.txt", "which profile to use")
	flag.Parse()
	src, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	parser := parse.New("single", string(src))
	cfg, opts, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	chars := make([]string, len(cfg.Characters.Profile))

	for i, v := range cfg.Characters.Profile {
		chars[i] = v.Base.Key.String()
	}

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout = w

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	opts.LogDetails = true
	opts.Debug = true
	opts.DebugPaths = []string{"stdout"}

	result, err := gcsim.Run(string(src), opts, func(s *gcsim.Simulation) error {
		var err error
		s.C.Queue, err = createQueue(cfg, s)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	w.Close()
	os.Stdout = old
	out := <-outC

	err = logtohtml.WriteString(out, "./debug.html", cfg.Characters.Initial.String(), chars)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(result)

	fmt.Print(result.PrettyPrint())

}

func createQueue(cfg core.Config, s *gcsim.Simulation) (core.QueueHandler, error) {
	cust := make(map[string]int)
	for i, v := range cfg.Rotation {
		if v.Name != "" {
			cust[v.Name] = i
		}
		// log.Println(v.Conditions)
	}
	for _, v := range cfg.Rotation {
		if _, ok := s.C.CharByName(v.Target); !ok {
			return nil, fmt.Errorf("invalid char in rotation %v", v.Target)
		}
	}

	r := calcqueue.New(s.C)
	r.SetActionList(cfg.Rotation)

	return r, nil
}
