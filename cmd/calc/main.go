package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/genshinsim/gsim/internal/logtohtml"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/parse"
)

func main() {
	cfgFile := flag.String("c", "config.txt", "which profile to use")
	flag.Parse()
	src, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	parser := parse.New("single", string(src))
	cfg, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	chars := make([]string, len(cfg.Characters.Profile))

	for i, v := range cfg.Characters.Profile {
		chars[i] = v.Base.Name
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

	result, err := combat.Run(string(src), true, true, func(s *combat.Simulation) error {
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

	err = logtohtml.WriteString(out, "./debug.html", cfg.Characters.Initial, chars)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println(result)

	fmt.Print(result.PrettyPrint())

}

func createQueue(cfg core.Config, s *combat.Simulation) (*queue, error) {
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

	return &queue{
		prio: cfg.Rotation,
		core: s.C,
	}, nil
}

type queue struct {
	core *core.Core
	prio []core.Action
	ind  int
}

func (q *queue) SetActionList(a []core.Action) {
	q.prio = a
}

func (q *queue) Next() ([]core.ActionItem, error) {

	var r []core.ActionItem
	active := q.core.Chars[q.core.ActiveChar].Name()
	for {
		if q.ind >= len(q.prio) {
			q.core.Log.Debugw(
				"no more actions",
				"frame", q.core.F,
				"event", core.LogQueueEvent,
			)
			return nil, nil
		}
		//we only accept action+=?? target=character wait=150
		//also, go down the list 1 at a time
		v := q.prio[q.ind]

		if v.IsSeq {
			//ignore and move on
			q.ind++
			continue
		}

		//check wait
		if v.Wait > q.core.F {
			q.core.Log.Debugw(
				"item on wait",
				"frame", q.core.F,
				"event", core.LogQueueEvent,
				"wait", v.Wait,
				"index", q.ind,
				"name", v.Name,
				"target", v.Target,
				"is seq", v.IsSeq,
				"strict", v.IsStrict,
				"exec", v.Exec,
				"once", v.Once,
				"post", v.PostAction.String(),
				"swap_to", v.SwapTo,
				"raw", v.Raw,
			)
			return nil, nil
		}

		//check if we need to swap
		if active != v.Target {
			r = append(r, core.ActionItem{
				Target: v.Target,
				Typ:    core.ActionSwap,
			})
		}

		r = append(r, v.Exec[0])

		q.core.Log.Debugw(
			"item queued",
			"frame", q.core.F,
			"event", core.LogQueueEvent,
			"index", q.ind,
			"name", v.Name,
			"target", v.Target,
			"is seq", v.IsSeq,
			"strict", v.IsStrict,
			"exec", v.Exec,
			"once", v.Once,
			"post", v.PostAction.String(),
			"swap_to", v.SwapTo,
			"raw", v.Raw,
		)

		q.ind++
		return r, nil
	}

}
