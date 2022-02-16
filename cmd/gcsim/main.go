package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/genshinsim/gcsim"
	"github.com/genshinsim/gcsim/internal/logtohtml"
	"github.com/genshinsim/gcsim/internal/tmpl/calcqueue"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
	"go.uber.org/zap"
)

type opts struct {
	print        bool
	js           string
	debug        bool
	debugPath    string
	debugHTML    bool
	seconds      int
	config       string
	detailed     bool
	w            int
	i            int
	multi        string
	minmax       bool
	calc         bool
	ercalc       bool
	gz           bool
	substatOptim bool
	verbose      bool
	out          string
	options      string
}

func main() {

	var opt opts

	flag.BoolVar(&opt.print, "print", true, "print output to screen? default true")
	flag.StringVar(&opt.js, "js", "", "alias for 'out' flag")
	flag.BoolVar(&opt.debug, "d", false, "show debug? default false")
	flag.StringVar(&opt.debugPath, "do", "debug.html", "file name for debug file")
	flag.BoolVar(&opt.debugHTML, "dh", true, "output debug html? default true (but only matters if debug is enabled)")
	flag.IntVar(&opt.seconds, "s", 0, "how many seconds to run the sim for")
	flag.StringVar(&opt.config, "c", "config.txt", "which profile to use")
	flag.BoolVar(&opt.detailed, "t", true, "log combat details")
	flag.BoolVar(&opt.gz, "gz", false, "gzip json results; require js flag")
	// f := flag.String("o", "debug.log", "detailed log file")
	// hp := flag.Float64("hp", 0, "hp mode: how much hp to deal damage to")
	// showCaller := flag.Bool("caller", false, "show caller in debug low")
	// fixedRand := flag.Bool("noseed", false, "use 0 for rand seed always - guarantee same results every time; only in single mode")
	// avgMode := flag.Bool("a", false, "run sim multiple times and calculate avg damage (smooth out randomness). default false. note that there is no debug log in this mode")
	flag.IntVar(&opt.w, "w", 0, "number of workers to run when running multiple iterations; default 24")
	flag.IntVar(&opt.i, "i", 0, "number of iterations to run if we're running multiple")
	flag.StringVar(&opt.multi, "m", "", "mutiple config mode. Takes in a file with a line feed separated list of paths to config files (priority mode only) and outputs an abbreviated result")
	flag.BoolVar(&opt.minmax, "minmax", false, "track the min/max run seed and rerun those (single mode with debug only)")
	flag.BoolVar(&opt.calc, "calc", false, "run sim in calc mode")
	flag.BoolVar(&opt.ercalc, "ercalc", false, "run sim in ER calc mode")
	flag.BoolVar(&opt.substatOptim, "substatOptim", false, "optimize substats according to KQM standards. Set the out flag to output config with optimal substats inserted to a given file path")
	flag.BoolVar(&opt.verbose, "v", false, "Verbose output log (currently only for substat optimization)")
	flag.StringVar(&opt.out, "out", "", "output result to file? supply file path (otherwise empty string for disabled). default disabled")
	flag.StringVar(&opt.options, "options", "", `Additional options for substat optimization mode. Currently supports the following flags, set in a semi-colon delimited list (e.g. -options="total_liquid_substats=15;indiv_liquid_cap=8"):
- total_liquid_substats (default = 20): Total liquid substats available to be assigned across all substats
- indiv_liquid_cap (default = 10): Total liquid substats that can be assigned to a single substat
- fixed_substats_count (default = 2): Amount of fixed substats that are assigned to all substats
- sim_iter (default = 350): RECOMMENDED TO NOT TOUCH. Number of iterations used when optimizing. Only change (increase) this if you are working with a team with extremely high standard deviation (>25% of mean)
- tol_mean (default = 0.015): RECOMMENDED TO NOT TOUCH. Tolerance of changes in DPS mean used in ER optimization
- tol_sd (default = 0.33): RECOMMENDED TO NOT TOUCH. Tolerance of changes in DPS SD used in ER optimization`)

	// t := flag.Int("t", 1, "target multiplier")

	flag.Parse()

	// Keep js for backwards compatibility
	if opt.js != "" && opt.out == "" {
		opt.out = opt.js
	}

	if opt.multi != "" {
		content, err := ioutil.ReadFile(opt.multi)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		files := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
		// lines := strings.Split(string(content), `\n`)
		err = runMulti(files, opt.w, opt.i)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}

	if opt.substatOptim {
		runSubstatOptim(opt)
		return
	}

	runSingle(opt)
}

func runSingle(o opts) {

	src, err := ioutil.ReadFile(o.config)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//check for imports
	var data strings.Builder
	var re = regexp.MustCompile(`(?m)^import "(.+)"$`)

	rows := strings.Split(strings.ReplaceAll(string(src), "\r\n", "\n"), "\n")
	for _, row := range rows {
		if re.MatchString(row) {
			match := re.FindStringSubmatch(row)
			//read import
			p := path.Join(path.Dir(o.config), match[1])
			src, err = ioutil.ReadFile(p)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			data.WriteString(string(src))

		} else {
			data.WriteString(row)
			data.WriteString("\n")
		}
	}

	// fmt.Println(data.String())

	parser := parse.New("single", data.String())
	cfg, opts, err := parser.Parse()

	opts.ERCalcMode = o.ercalc

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if o.i > 0 {
		opts.Iteration = o.i
	}
	if o.w > 0 {
		opts.Workers = o.w
	}
	if o.debug {
		opts.Debug = o.debug
	}
	if o.seconds > 0 {
		opts.Duration = o.seconds
	}
	if o.detailed {
		opts.LogDetails = true
	}

	// log.Println(opts)
	// defer profile.Start(profile.ProfilePath("./")).Stop()

	var result gcsim.Result
	//if debug we're going to capture the logs
	if opts.Debug {
		r, w, err := os.Pipe()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		outC := make(chan string)
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		zap.RegisterSink("gsim", func(url *url.URL) (zap.Sink, error) {
			return w, nil
		})

		opts.DebugPaths = []string{"gsim://"}

		result, err = gcsim.Run(data.String(), cfg, opts, func(s *gcsim.Simulation) error {
			if !o.calc {
				return nil
			}
			var err error
			s.C.Queue, err = createQueue(cfg, s)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		w.Close()

		out := <-outC

		if o.debugHTML {
			chars := make([]string, len(cfg.Characters.Profile))
			for i, v := range cfg.Characters.Profile {
				chars[i] = v.Base.Key.String()
			}
			err = logtohtml.WriteString(out, o.debugPath, cfg.Characters.Initial.String(), chars)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		}

		result.Debug = out

	} else {
		result, err = gcsim.Run(data.String(), cfg, opts, func(s *gcsim.Simulation) error {
			if !o.calc {
				return nil
			}
			var err error
			s.C.Queue, err = createQueue(cfg, s)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	if o.print {
		fmt.Print(result.PrettyPrint())
	}

	if o.minmax && o.debug {

		minResult, err := runSeeded(data.String(), result.MinSeed, opts, o, "debugmin")
		if err != nil {
			log.Panic(err)
		}
		maxResult, err := runSeeded(data.String(), result.MaxSeed, opts, o, "debugmax")
		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("Min seed: %v | DPS: %v\n", result.MinSeed, minResult.DPS)
		fmt.Printf("Min seed: %v | DPS: %v\n", result.MaxSeed, maxResult.DPS)
	}

	if o.out != "" {
		//try creating file to write to
		result.Text = result.PrettyPrint()
		data, jsonErr := json.Marshal(result)
		if jsonErr != nil {
			log.Panic(jsonErr)
		}

		if o.gz {
			path := strings.TrimSuffix(o.out, filepath.Ext(o.out)) + ".json.gz"
			f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				log.Panic(err)
			}
			defer f.Close()
			zw := gzip.NewWriter(f)
			zw.Write(data)
			err = zw.Close()
			if err != nil {
				log.Panic(err)
			}
		} else {
			err := os.WriteFile(o.out, data, 0644)
			if err != nil {
				log.Panic(err)
			}

		}

	}
}

func runSeeded(data string, seed int64, opts core.RunOpt, o opts, file string) (gcsim.Stats, error) {
	r, w, err := os.Pipe()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	zap.RegisterSink(file, func(url *url.URL) (zap.Sink, error) {
		return w, nil
	})

	opts.DebugPaths = []string{fmt.Sprintf("%v://", file)}

	parser := parse.New("single", data)
	cfg, _, _ := parser.Parse()

	sim, err := gcsim.New(cfg, seed, opts, func(s *gcsim.Simulation) error {
		if !o.calc {
			return nil
		}
		var err error
		s.C.Queue, err = createQueue(cfg, s)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return gcsim.Stats{}, err
	}

	v, err := sim.Run()
	if err != nil {
		return gcsim.Stats{}, err
	}

	w.Close()

	out := <-outC

	if file != "" {

		chars := make([]string, len(cfg.Characters.Profile))
		for i, v := range cfg.Characters.Profile {
			chars[i] = v.Base.Key.String()
		}
		err = logtohtml.WriteString(out, fmt.Sprintf("./%v.html", file), cfg.Characters.Initial.String(), chars)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	return v, nil
}

func runMulti(files []string, w, i int) error {
	fmt.Print("Filename                                                     |      Mean|       Min|       Max|   Std Dev|   HP Mode|     Iters|\n")
	fmt.Print("--------------------------------------------------------------------------------------------------------------------------------\n")
	for _, f := range files {
		if f == "" || f[0] == '#' {
			continue
		}
		src, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		var data strings.Builder
		var re = regexp.MustCompile(`(?m)^import "(.+)"$`)

		rows := strings.Split(strings.ReplaceAll(string(src), "\r\n", "\n"), "\n")
		for _, row := range rows {
			if re.MatchString(row) {
				match := re.FindStringSubmatch(row)
				//read import
				p := path.Join(path.Dir(f), match[1])
				src, err = ioutil.ReadFile(p)
				if err != nil {
					log.Fatal(err)
				}

				data.WriteString(string(src))

			} else {
				data.WriteString(row)
				data.WriteString("\n")
			}
		}

		parser := parse.New("single", data.String())
		cfg, opts, err := parser.Parse()
		if err != nil {
			return err
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
		r, err := gcsim.Run(data.String(), cfg, opts)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println(r)
		fmt.Printf("%10.2f|%10.2f|%10.2f|%10.2f|%10.10v|%10d|\n", r.DPS.Mean, r.DPS.Min, r.DPS.Max, r.DPS.SD, r.IsDamageMode, r.Iterations)
	}
	return nil
}

func createQueue(cfg core.SimulationConfig, s *gcsim.Simulation) (core.QueueHandler, error) {
	for _, v := range cfg.Rotation {
		if _, ok := s.C.CharByName(v.SequenceChar); v.Type == core.ActionBlockTypeSequence && !ok {
			return nil, fmt.Errorf("invalid char in rotation %v; %v", v.SequenceChar, v)
		}
	}

	r := calcqueue.New(s.C)
	r.SetActionList(cfg.Rotation)

	return r, nil
}
