package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/genshinsim/gsim/internal/logtohtml"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/parse"
	"github.com/pkg/profile"
)

func main() {

	var src []byte

	var err error

	debug := flag.Bool("debug", false, "show debug?")
	seconds := flag.Int("s", 90, "how many seconds to run the sim for")
	cfgFile := flag.String("c", "config.txt", "which profile to use")
	detailed := flag.Bool("d", true, "log combat details")
	// f := flag.String("o", "debug.log", "detailed log file")
	// hp := flag.Float64("hp", 0, "hp mode: how much hp to deal damage to")
	// showCaller := flag.Bool("caller", false, "show caller in debug low")
	// fixedRand := flag.Bool("noseed", false, "use 0 for rand seed always - guarantee same results every time; only in single mode")
	// avgMode := flag.Bool("a", false, "run sim multiple times and calculate avg damage (smooth out randomness). default false. note that there is no debug log in this mode")
	w := flag.Int("w", 0, "number of workers to run when running multiple iterations; default 24")
	i := flag.Int("i", 0, "number of iterations to run if we're running multiple")
	multi := flag.String("m", "", "mutiple config mode")
	// t := flag.Int("t", 1, "target multiplier")

	flag.Parse()

	if *multi != "" {
		content, err := ioutil.ReadFile(*multi)
		if err != nil {
			log.Fatal(err)
		}
		files := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
		// lines := strings.Split(string(content), `\n`)
		runMulti(files, *w, *i)
		return
	}

	src, err = ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	parser := parse.New("single", string(src))
	cfg, opts, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	if *i > 0 {
		opts.Iteration = *i
	}
	if *w > 0 {
		opts.Workers = *w
	}
	if *debug {
		opts.Debug = true
	}
	if *seconds > 0 {
		opts.Duration = *seconds
	}
	if *detailed {
		opts.LogDetails = true
	}

	log.Println(opts)

	defer elapsed("simulation completed")()
	defer profile.Start(profile.ProfilePath("./")).Stop()

	//if debug we're going to capture the logs
	if opts.Debug {

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

		opts.DebugPaths = []string{"stdout"}

		result, err := combat.Run(string(src), opts)
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

		fmt.Print(result.PrettyPrint())

		// fmt.Print(out)

	} else {
		result, err := combat.Run(string(src), opts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(result.PrettyPrint())
	}

}

func runMulti(files []string, w, i int) {
	fmt.Print("Filename                                                     |      Mean|       Min|       Max|   Std Dev|   HP Mode|     Iters|\n")
	fmt.Print("--------------------------------------------------------------------------------------------------------------------------------\n")
	for _, f := range files {
		if f == "" || f[0] == '#' {
			continue
		}
		source, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		parser := parse.New("single", string(source))
		_, opts, err := parser.Parse()
		if err != nil {
			log.Fatal(err)
		}
		if w > 0 {
			opts.Workers = w
		}
		if i > 0 {
			opts.Iteration = i
		}
		opts.Debug = false
		opts.LogDetails = false

		fmt.Printf("%60.60v |", f)
		r, err := combat.Run(string(source), opts)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println(r)
		fmt.Printf("%10.2f|%10.2f|%10.2f|%10.2f|%10.10v|%10d|\n", r.DPS.Mean, r.DPS.Min, r.DPS.Max, r.DPS.SD, r.IsDamageMode, r.Iterations)
	}
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}
