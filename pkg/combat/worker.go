package combat

import (
	"errors"
	"fmt"
	"log"
	"math"
	"runtime"
	"sort"
	"strings"

	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/parse"
)

type Stats struct {
	IsDamageMode         bool                      `json:"is_damage_mode"`
	CharNames            []string                  `json:"char_names"`
	DamageByChar         []map[string]float64      `json:"damage_by_char"`
	CharActiveTime       []int                     `json:"char_active_time"`
	AbilUsageCountByChar []map[string]int          `json:"abil_usage_count_by_char"`
	ReactionsTriggered   map[core.ReactionType]int `json:"reactions_triggered"`
	Duration             int                       `json:"sim_duration"`
	//final result
	Damage float64 `json:"damage"`
	DPS    float64 `json:"dps"`
}

type AverageStats struct {
	IsDamageMode         bool                            `json:"is_damage_mode"`
	CharNames            []string                        `json:"char_names"`
	DamageByChar         []map[string]FloatResult        `json:"damage_by_char"`
	CharActiveTime       []IntResult                     `json:"char_active_time"`
	AbilUsageCountByChar []map[string]IntResult          `json:"abil_usage_count_by_char"`
	ReactionsTriggered   map[core.ReactionType]IntResult `json:"reactions_triggered"`
	Duration             FloatResult                     `json:"sim_duration"`
	//final result
	Damage     FloatResult `json:"damage"`
	DPS        FloatResult `json:"dps"`
	Iterations int
}

type IntResult struct {
	Min  int     `json:"min"`
	Max  int     `json:"max"`
	Mean float64 `json:"mean"`
	SD   float64 `json:"sd"`
}

type FloatResult struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Mean float64 `json:"mean"`
	SD   float64 `json:"sd"`
}

type workerResp struct {
	stats Stats
	err   error
}

func Run(src string, details bool, debug bool, cust ...func(*Simulation) error) (AverageStats, error) {

	//options mode=damage debug=true iteration=5000 duration=90 workers=24;
	var data []Stats

	parser := parse.New("single", string(src))
	cfg, err := parser.Parse()
	if err != nil {
		return AverageStats{}, err
	}

	charCount := len(cfg.Characters.Profile)

	if charCount > 4 {
		return AverageStats{}, errors.New("cannot have more than 4 characters in a team")
	}

	chars := make([]string, len(cfg.Characters.Profile))
	for i, v := range cfg.Characters.Profile {
		chars[i] = v.Base.Name
	}

	//set defaults if nothing specified
	count := cfg.RunOptions.Iteration
	if count == 0 {
		count = 1000
	}
	n := count

	if debug {
		count--
	}

	// fmt.Printf("running %v iterations\n", count)

	w := cfg.RunOptions.Workers
	if w <= 0 {
		w = runtime.NumCPU()
	}
	if w > n {
		w = n
	}

	resp := make(chan workerResp, count)
	req := make(chan bool)
	done := make(chan bool)

	for i := 0; i < w; i++ {
		go worker(src, details, resp, req, done, cust...)
	}

	go func() {
		var wip int
		for wip < n {
			req <- true
			wip++
		}
	}()

	defer close(done)

	for count > 0 {
		vv := <-resp
		if vv.err != nil {
			return AverageStats{}, vv.err
		}
		v := vv.stats
		// log.Println(v)
		count--
		data = append(data, v)
	}

	//if debug is true, run one more purely for debug do not add to stats
	if debug {
		s, err := NewSim(cfg, details, cust...)
		if err != nil {
			log.Fatal(err)
		}
		v, err := s.Run()
		if err != nil {
			log.Fatal(err)
		}
		// log.Println(v)
		data = append(data, v)
	}

	result := collectResult(data, cfg.RunOptions.DamageMode, chars, details)
	result.Iterations = n
	if !cfg.RunOptions.DamageMode {
		result.Duration.Mean = float64(cfg.RunOptions.Duration)
		result.Duration.Min = float64(cfg.RunOptions.Duration)
		result.Duration.Max = float64(cfg.RunOptions.Duration)
	}

	return result, nil
}

func collectResult(data []Stats, mode bool, chars []string, detailed bool) (result AverageStats) {

	charCount := len(chars)
	result.DPS.Min = math.MaxFloat64
	result.DPS.Max = -1
	if detailed {
		result.ReactionsTriggered = make(map[core.ReactionType]IntResult)
		result.CharNames = make([]string, charCount)
		result.AbilUsageCountByChar = make([]map[string]IntResult, charCount)
		result.CharActiveTime = make([]IntResult, charCount)
		result.DamageByChar = make([]map[string]FloatResult, charCount)

		for i := range result.CharNames {
			result.CharNames[i] = chars[i]
			result.CharActiveTime[i].Min = math.MaxInt64
			result.CharActiveTime[i].Max = -1
			result.AbilUsageCountByChar[i] = make(map[string]IntResult)
			result.DamageByChar[i] = make(map[string]FloatResult)
		}
	}

	n := len(data)

	// var dd float64

	for _, v := range data {
		dd := float64(v.Duration) / 60 //sim reports in frames
		result.Duration.Mean += dd / float64(n)
		if dd > result.Duration.Max {
			result.Duration.Max = dd
		}
		if dd < result.Duration.Mean {
			result.Duration.Max = dd
		}

		//dmg
		if v.Damage < result.Damage.Min {
			result.DPS.Min = v.Damage
		}
		if v.Damage > result.Damage.Max {
			result.Damage.Max = v.Damage
		}
		result.Damage.Mean += v.Damage / float64(n)

		//dps
		if result.DPS.Min > v.DPS {
			result.DPS.Min = v.DPS
		}
		if result.DPS.Max < v.DPS {
			result.DPS.Max = v.DPS
		}
		result.DPS.Mean += v.DPS / float64(n)

		if !detailed {
			continue
		}

		//char active time
		for i, amt := range v.CharActiveTime {

			if result.CharActiveTime[i].Min > amt {
				result.CharActiveTime[i].Min = amt
			}
			if result.CharActiveTime[i].Max < amt {
				result.CharActiveTime[i].Max = amt
			}
			result.CharActiveTime[i].Mean += float64(amt) / float64(n)
		}

		//dmg by char
		for i, abil := range v.DamageByChar {
			for k, amt := range abil {
				x, ok := result.DamageByChar[i][k]
				if !ok {
					x.Min = math.MaxFloat64
					x.Max = -1
				}
				// log.Printf("dmg amount: %v\n", amt)
				amt = amt / float64(dd)
				if x.Min > amt {
					x.Min = amt
				}
				if x.Max < amt {
					x.Max = amt
				}
				x.Mean += amt / float64(n)

				result.DamageByChar[i][k] = x
			}
		}

		//abil use
		for c, abil := range v.AbilUsageCountByChar {
			for k, amt := range abil {
				x, ok := result.AbilUsageCountByChar[c][k]
				if !ok {
					x.Min = math.MaxInt64
					x.Max = -1
				}
				if x.Min > amt {
					x.Min = amt
				}
				if x.Max < amt {
					x.Max = amt
				}
				x.Mean += float64(amt) / float64(n)

				result.AbilUsageCountByChar[c][k] = x
			}
		}

		//reactions
		for c, amt := range v.ReactionsTriggered {
			x, ok := result.ReactionsTriggered[c]
			if !ok {
				x.Min = math.MaxInt64
				x.Max = -1
			}
			if x.Min > amt {
				x.Min = amt
			}
			if x.Max < amt {
				x.Max = amt
			}
			x.Mean += float64(amt) / float64(n)

			result.ReactionsTriggered[c] = x
		}
	}

	for _, v := range data {
		result.DPS.SD += (v.DPS - result.DPS.Mean) * (v.DPS - result.DPS.Mean)
		if mode {
			result.Duration.SD += (float64(v.Duration) - result.Duration.Mean) * (float64(v.Duration) - result.Duration.Mean)
		}
	}

	result.DPS.SD = math.Sqrt(result.DPS.SD / float64(n))

	return
}

func worker(src string, details bool, resp chan workerResp, req chan bool, done chan bool, cust ...func(*Simulation) error) {

	for {
		select {
		case <-req:
			parser := parse.New("single", src)
			cfg, _ := parser.Parse()
			cfg.RunOptions.Debug = false

			s, err := NewSim(cfg, details, cust...)
			if err != nil {
				resp <- workerResp{
					err: err,
				}
				return
			}

			stat, err := s.Run()

			if err != nil {
				resp <- workerResp{
					err: err,
				}
				return
			}

			resp <- workerResp{
				stats: stat,
				err:   nil,
			}
		case <-done:
			return
		}
	}
}

func (stats *AverageStats) PrettyPrint() string {

	var sb strings.Builder

	for i, t := range stats.DamageByChar {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
		}
		sb.WriteString(fmt.Sprintf("%v contributed the following dps:\n", stats.CharNames[i]))
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var total float64
		for _, k := range keys {
			v := t[k]
			sb.WriteString(fmt.Sprintf("\t%v (%.2f%% of total): avg %.2f [min: %.2f | max: %.2f] \n", k, 100*v.Mean/stats.DPS.Mean, v.Mean, v.Min, v.Max))
			total += v.Mean
		}

		sb.WriteString(fmt.Sprintf("%v total avg dps: %.2f; total percentage: %.0f%%\n", stats.CharNames[i], total, 100*total/stats.DPS.Mean))
	}
	for i, t := range stats.AbilUsageCountByChar {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Character ability usage:\n")
		}
		sb.WriteString(fmt.Sprintf("%v used the following abilities:\n", stats.CharNames[i]))
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t[k]
			sb.WriteString(fmt.Sprintf("\t%v: avg %.2f [min: %v | max: %v]\n", k, v.Mean, v.Min, v.Max))
		}
	}
	for i, v := range stats.CharActiveTime {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Character field time:\n")
		}
		sb.WriteString(fmt.Sprintf("%v on average active for %.0f%% [min: %.0f%% | max: %.0f%%]\n", stats.CharNames[i], 100*v.Mean/(stats.Duration.Mean*60), float64(100*v.Min)/(stats.Duration.Mean*60), float64(100*v.Max)/(stats.Duration.Mean*60)))
	}
	rk := make([]core.ReactionType, 0, len(stats.ReactionsTriggered))
	for k := range stats.ReactionsTriggered {
		rk = append(rk, k)
	}
	for i, k := range rk {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Reactions:\n")
		}
		v := stats.ReactionsTriggered[k]
		sb.WriteString(fmt.Sprintf("\t%v: avg %.2f [min: %v max: %v]\n", k, v.Mean, v.Min, v.Max))
	}
	sb.WriteString(fmt.Sprintf("Average %.2f damage over %.2f seconds, resulting in %.0f dps (min: %.2f max: %.2f std: %.2f) \n", stats.Damage.Mean, stats.Duration.Mean, stats.DPS.Mean, stats.DPS.Min, stats.DPS.Max, stats.DPS.SD))

	return sb.String()
}
