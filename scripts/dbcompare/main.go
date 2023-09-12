package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/genshinsim/gcsim/pkg/simulator"
)

type dbGetOpt struct {
	Query   map[string]interface{} `json:"query,omitempty"`
	Sort    map[string]interface{} `json:"sort,omitempty"`
	Project map[string]interface{} `json:"project,omitempty"`
	Skip    int64                  `json:"skip,omitempty"`
	Limit   int64                  `json:"limit,omitempty"`
}

type dbData struct {
	Data []dbEntry `json:"data"`
}

type dbEntry struct {
	Id      string `json:"_id"`
	Config  string `json:"config"`
	Summary struct {
		MeanDpsPerTarget float64 `json:"mean_dps_per_target"`
	} `json:"summary"`
}

func main() {
	res, err := getDBEntries()
	if err != nil {
		panic(err)
	}
	// fmt.Println(len(res))
	fmt.Println("id,original,next,err")
	for _, v := range res {
		dps, err := runSim(v)
		fmt.Printf("%v,%v,%v,%v\n", v.Id, v.Summary.MeanDpsPerTarget, dps, err)
	}
}

func runSim(w dbEntry) (float64, error) {
	// start := time.Now()
	// compute work??
	// log.Printf("got work %v; starting compute", w.Id)
	// compute result
	simcfg, gcsl, err := simulator.Parse(w.Config)
	if err != nil {
		// log.Printf("could not parse config for id %v: %v\n", w.Id, err)
		//TODO: we should post something here??
		return 0, err
	}
	simcfg.Settings.Iterations = 1000
	simcfg.Settings.NumberOfWorkers = 30

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	result, err := simulator.RunWithConfig(w.Config, simcfg, gcsl, simulator.Options{}, time.Now(), ctx)
	if err != nil {
		// log.Printf("error running sim %v: %v\n", w.Id, err)
		return 0, err
	}
	avgPerTarget := *result.Statistics.TotalDamage.Mean / (float64(len(result.TargetDetails)) * *result.Statistics.Duration.Mean)

	// elapsed := time.Since(start)
	// fmt.Printf()
	return avgPerTarget, nil
}

func getDBEntries() ([]dbEntry, error) {
	q := dbGetOpt{
		Limit: 100,
	}

	var page int64
	var res []dbEntry

	for {
		q.Skip = 100 * page
		jsonStr, _ := json.Marshal(q)
		url := fmt.Sprintf("https://simpact.app/api/db?q=%v", url.QueryEscape(string(jsonStr)))

		var d dbData
		err := getJSON(url, &d)
		if err != nil {
			return nil, err
		}
		if len(d.Data) == 0 {
			break
		}
		res = append(res, d.Data...)

		page++
	}

	return res, nil
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target any) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
