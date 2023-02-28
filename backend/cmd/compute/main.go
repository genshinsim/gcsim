package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"github.com/joho/godotenv"
)

var (
	sha1ver  string
	shareKey string
)

type opts struct {
	version bool
	max     int
	apikey  string
}

func init() {
	info, _ := debug.ReadBuildInfo()
	for _, bs := range info.Settings {
		if bs.Key == "vcs.revision" {
			sha1ver = bs.Value
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var opt opts
	flag.BoolVar(&opt.version, "version", false, "compute cli version")
	flag.IntVar(&opt.max, "max", -1, "max number of entries to compute")
	flag.StringVar(&opt.apikey, "key", "", "api key to used for authentication purposes")

	flag.Parse()

	if opt.version {
		fmt.Println(sha1ver)
		return
	}

	if shareKey == "" {
		shareKey = os.Getenv("GCSIM_SHARE_KEY")
	}

	if opt.apikey == "" {
		opt.apikey = os.Getenv("COMPUTE_API_KEY")
	}

	log.Println("compute version: " + sha1ver)
	// return
	//steps:
	// 1. ask backend server for compute work to do
	// 2. run sim
	// 3. post result to callback url
	log.Printf("start looking for work...")
	count := 0
	for opt.max == -1 || count < opt.max {
		err := processWork(opt.apikey)
		switch err {
		case nil:
		case errNoMoreWork:
			log.Println("no more work; all done!")
			return
		case errSimFailed:
			continue
		default:
			log.Fatalf("compute failed with err: %v", err)

		}
		count++
	}
	log.Println("all done")

}

var errNoMoreWork = errors.New("no more work")

var errSimFailed = errors.New("sim failed")

func processWork(key string) error {
	var w work
	start := time.Now()
	w, err := getWork(key)
	if err != nil {
		log.Fatalf("error getting work: %v", err)
	}
	//blank key means no more work
	if w.Id == "" {
		return errNoMoreWork
	}
	//compute work??
	log.Printf("got work %v; starting compute", w.Id)
	// compute result
	simcfg, err := simulator.Parse(w.Config)
	if err != nil {
		log.Printf("could not parse config for id %v: %v\n", w.Id, err)
		//TODO: we should post something here??
		return errSimFailed
	}
	simcfg.Settings.Iterations = w.Iterations
	simcfg.Settings.NumberOfWorkers = 30

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	result, err := simulator.RunWithConfig(w.Config, simcfg, simulator.Options{}, time.Now(), ctx)
	if err != nil {
		log.Printf("error running sim %v: %v\n", w.Id, err)
		return errSimFailed
	}

	t := model.ComputeWorkSource_name[int32(w.Source)]

	err = postResult(result, key, w.Id, t)
	if err != nil {
		log.Printf("error posting result: %v\n", err)
		return err
	}

	elapsed := time.Since(start)
	log.Printf("Work %v took %s", w.Id, elapsed)
	return nil
}

func postResult(result *model.SimulationResult, key, id, t string) error {
	hash, _ := result.Sign(shareKey)
	data, _ := result.MarshalJson()
	req, err := http.NewRequest("POST", "https://simimpact.app/api/db/compute/work", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("X-GCSIM-COMPUTE-API-KEY", key)
	req.Header.Add("X-GCSIM-SIGNED-HASH", hash)
	req.Header.Add("X-GCSIM-COMPUTE-SRC", t)
	req.Header.Add("X-GCSIM-COMPUTE-ID", id)
	_, err = http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func getJson(url, key string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-GCSIM-COMPUTE-API-KEY", key)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

type work struct {
	Id         string `json:"_id"`
	Config     string `json:"config"`
	Source     int    `json:"source"`
	Iterations int    `json:"iterations"`
}

func getWork(key string) (work, error) {
	const url = `https://simimpact.app/api/db/compute/work`
	var w work
	err := getJson(url, key, &w)
	if err != nil {
		return work{}, err
	}

	return w, nil
}
