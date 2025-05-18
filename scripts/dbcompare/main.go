package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
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

type opts struct {
	mustHaveChar string
	iters        int
	workers      int
}

func main() {
	var opt opts
	flag.StringVar(&opt.mustHaveChar, "must", "", "comma separated character names that must be present, otherwise skip. default blank")
	flag.IntVar(&opt.iters, "iters", 1000, "iterations to run each comparison")
	flag.IntVar(&opt.workers, "workers", 30, "number of workers to use")
	flag.Parse()
	res, err := getDBEntries()
	if err != nil {
		panic(err)
	}
	filters := strings.Split(opt.mustHaveChar, ",")
	fmt.Println("id,original,next,diff,abs_diff,per_diff,err")
	for _, v := range res {
		simcfg, gcsl, err := simulator.Parse(v.Config)
		if err != nil {
			fmt.Printf(",,,,,,,%v\n", err)
			continue
		}
		simcfg.Settings.Iterations = opt.iters
		simcfg.Settings.NumberOfWorkers = opt.workers

		var team string
		for i := range simcfg.Characters {
			team += "+" + simcfg.Characters[i].Base.Key.String()
		}
		team = strings.TrimPrefix(team, "+")
		// check for ignores
		if len(filters) > 0 {
			// teams must contain all filters
			allFound := true
			for _, v := range filters {
				if !strings.Contains(team, v) {
					allFound = false
					break
				}
			}
			if !allFound {
				continue
			}
		}
		dps, err := runSim(simcfg, gcsl, v.Config)
		diff := dps - v.Summary.MeanDpsPerTarget
		absDiff := math.Abs(diff)
		percentDiff := absDiff / v.Summary.MeanDpsPerTarget
		fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v\n", team, v.Id, v.Summary.MeanDpsPerTarget, dps, diff, absDiff, percentDiff, err)
	}
}

func runSim(simcfg *info.ActionList, gcsl ast.Node, cfg string) (float64, error) {
	// start := time.Now()
	// compute work??
	// log.Printf("got work %v; starting compute", w.Id)
	// compute result
	simcfg.Settings.CollectStats = []string{"overview"}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	result, err := simulator.RunWithConfig(ctx, cfg, simcfg, gcsl, simulator.Options{}, time.Now())
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
		jsonStr, err := json.Marshal(q)
		if err != nil {
			return nil, err
		}
		url := fmt.Sprintf("https://simpact.app/api/db?q=%v", url.QueryEscape(string(jsonStr)))

		var d dbData
		err = getJSON(url, &d)
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
