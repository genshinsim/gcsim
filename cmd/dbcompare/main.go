package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/genshinsim/gcsim/internal/simulator"
)

type DBData struct {
	Config    string  `json:"config"`
	DPS       float64 `json:"dps"`
	ViewerKey string  `json:"viewer_key"`
}

type Result struct {
	DPS struct {
		Mean float64 `json:"mean"`
		SD   float64 `json:"sd"`
	} `json:"dps"`
}

func main() {

	fPtr := flag.String("f", "", "check only the specified hash")
	flag.Parse()

	//fetch from: https://viewer.gcsim.workers.dev/gcsimdb
	var data []DBData
	getJson("https://viewer.gcsim.workers.dev/gcsimdb", &data)

	for _, v := range data {
		if *fPtr != "" && v.ViewerKey != *fPtr {
			continue
		}
		err := runAndCompare(v)
		if err != nil {
			fmt.Printf("error encountered running file: %v\n", v.ViewerKey)
			panic(err)
		}
	}

}

func runAndCompare(d DBData) error {
	//write config to file
	f := "./data/" + d.ViewerKey + ".txt"
	fj := "./data/" + d.ViewerKey + ".json"
	os.Remove(f)
	os.Remove(fj)
	os.WriteFile(f, []byte(d.Config), 0700)
	fmt.Printf("Comparing %v: ", d.ViewerKey)
	simopt := simulator.Options{
		ConfigPath:       f,
		ResultSaveToPath: fj,
	}
	_, err := simulator.Run(simopt)
	if err != nil {
		return err
	}
	//read the file and find dps
	resbyte, err := os.ReadFile(fj)
	if err != nil {
		return err
	}
	var res Result
	err = json.Unmarshal(resbyte, &res)
	if err != nil {
		return err
	}
	//compare
	fmt.Printf("original %.0f, new %.0f (sd: %.0f), diff %.0f (p of sd: %.2f)\n", d.DPS, res.DPS.Mean, res.DPS.SD, res.DPS.Mean-d.DPS, (res.DPS.Mean-d.DPS)/res.DPS.SD)

	return nil
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
