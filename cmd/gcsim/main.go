package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/genshinsim/gcsim/pkg/optimization"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/sample"
	"github.com/genshinsim/gcsim/pkg/simulator"
)

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

type opts struct {
	config       string
	out          string //file result name
	sample       string //file sample name
	gz           bool
	serve        bool
	nobrowser    bool
	norun        bool
	keepserving  bool
	substatOptim bool
	verbose      bool
	options      string
}

// command line tool; following options are available:
func main() {

	var opt opts
	var version bool
	flag.BoolVar(&version, "version", false, "check gcsim version (git hash)")
	flag.StringVar(&opt.config, "c", "config.txt", "which profile to use; default config.txt")
	flag.StringVar(&opt.out, "out", "", "output result to file? supply file path (otherwise empty string for disabled). default disabled")
	flag.StringVar(&opt.sample, "sample", "", "create sample result. supply file path (otherwise empty string for disabled). default disabled")
	flag.BoolVar(&opt.gz, "gz", false, "gzip json results; require out flag")
	flag.BoolVar(&opt.serve, "s", false, "serve results to viewer (local). default false")
	flag.BoolVar(&opt.norun, "nr", false, "disable running the simulation (useful if you only want to generate a sample")
	flag.BoolVar(&opt.nobrowser, "nb", false, "disable opening default browser")
	flag.BoolVar(&opt.keepserving, "ks", false, "keep serving same results without terminating web server")
	flag.BoolVar(&opt.substatOptim, "substatOptim", false, "optimize substats according to KQM standards. Set the out flag to output config with optimal substats inserted to a given file path")
	flag.BoolVar(&opt.verbose, "v", false, "Verbose output log (currently only for substat optimization)")
	flag.StringVar(&opt.options, "options", "", `Additional options for substat optimization mode. Currently supports the following flags, set in a semi-colon delimited list (e.g. -options="total_liquid_substats=15;indiv_liquid_cap=8"):
- total_liquid_substats (default = 20): Total liquid substats available to be assigned across all substats
- indiv_liquid_cap (default = 10): Total liquid substats that can be assigned to a single substat
- fixed_substats_count (default = 2): Amount of fixed substats that are assigned to all substats
- sim_iter (default = 350): RECOMMENDED TO NOT TOUCH. Number of iterations used when optimizing. Only change (increase) this if you are working with a team with extremely high standard deviation (>25% of mean)
- tol_mean (default = 0.015): RECOMMENDED TO NOT TOUCH. Tolerance of changes in DPS mean used in ER optimization
- tol_sd (default = 0.33): RECOMMENDED TO NOT TOUCH. Tolerance of changes in DPS SD used in ER optimization`)

	flag.Parse()

	if version {
		fmt.Println(sha1ver)
		return
	}

	if opt.serve {
		os.Remove("serve_data.json.gz") //not really needed since we truncate anyways
		opt.out = "serve_data.json"
		opt.gz = true
	}

	simopt := simulator.Options{
		ConfigPath:       opt.config,
		ResultSaveToPath: opt.out,
		GZIPResult:       opt.gz,
		Version:          sha1ver,
		BuildDate:        buildTime,
	}

	if opt.substatOptim {
		// TODO: Eventually will want to handle verbose/options in some other way.
		// Ideally once documentation is standardized, can move options to a config file, and verbose can also be moved into options or something
		optimization.RunSubstatOptim(simopt, opt.verbose, opt.options)
		return
	}

	// TODO: should perform the config parsing here and then share the parsed results between run & sample

	var res result.Summary
	if !opt.norun {
		res, err := simulator.Run(simopt)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(res.PrettyPrint())
	}

	if opt.sample != "" {
		var seed uint64
		if opt.norun {
			seed = uint64(simulator.CryptoRandSeed())
		} else {
			seed, _ = strconv.ParseUint(res.SampleSeed, 10, 64)
		}

		cfg, err := simulator.ReadConfig(opt.config)
		if err != nil {
			log.Println(err)
			return
		}

		sample, err := sample.GenerateSampleWithSeed(cfg, seed)
		if err != nil {
			log.Println(err)
			return
		}
		sample.Save(opt.sample, opt.gz)
		fmt.Printf("Generated sample with seed: %v\n", seed)
	}

	if opt.serve && !opt.norun {
		fmt.Println("Serving result to HTTP...")
		//start server to listen for token
		serverDone := &sync.WaitGroup{}
		serverDone.Add(1)
		serveLocal(serverDone, "./serve_data.json.gz", opt.keepserving)
		url := "https://gcsim.app/viewer/local"
		if !opt.nobrowser {
			err := open(url)
			if err != nil {
				//try "xdg-open-wsl"
				err = openWSL(url)
				if err != nil {
					fmt.Printf("Error opening default browser... please visit: %v\n", url)
				}
			}
		}
		serverDone.Wait()
	}
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

func openWSL(url string) error {
	cmd := "powershell.exe"
	args := []string{"/c", "start", url}
	return exec.Command(cmd, args...).Start()
}

var ctxShutdown, cancel = context.WithCancel(context.Background())

var quit = make(chan bool, 1)

type viewerData struct {
	Data        string `json:"data"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func serveLocal(wg *sync.WaitGroup, path string, keepserving bool) {
	srv := &http.Server{Addr: "127.0.0.1:8381"}
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctxShutdown.Done():
			fmt.Println("HTTP server shuting down ...")
			return
		default:
		}
		//check CORS
		switch r.Method {
		case "OPTIONS":
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Access-Control-Allow-Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
			w.WriteHeader(http.StatusNoContent)
			fmt.Println("OPTIONS request received, responding")
			return
		case "GET":
		default:
			fmt.Printf("Invalid request method: %v\n", r.Method)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		//read the gz file
		gzData, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("error reading gz data: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b64string := base64.StdEncoding.EncodeToString(gzData)

		x := viewerData{
			Data:        b64string,
			Author:      "none",
			Description: "none",
		}

		jsonData, err := json.Marshal(x)
		if err != nil {
			fmt.Printf("error marshal json: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		if !keepserving {
			// // Shut down server here
			cancel() // to say sorry, above.

			close(quit)
		}

	})

	go gracefullShutdown(srv)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		fmt.Println("HTTP server closed.")
	}()
}

func gracefullShutdown(server *http.Server) {
	<-quit
	fmt.Println("Server is shutting down...")

	ctx, c := context.WithTimeout(context.Background(), 30*time.Second)
	defer c()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}
