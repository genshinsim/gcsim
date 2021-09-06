package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

type runConfig struct {
	Options core.RunOpt `json:"options"`
	Config  string      `json:"config"`
}

func (s *Server) handleRun(ctx context.Context, r wsRequest) {

	s.Log.Debugw("handleRun: request to run received")

	var cfg runConfig
	err := json.Unmarshal([]byte(r.Payload), &cfg)

	if err != nil {
		s.Log.Debugw("handleRun: invalid request payload", "payload", r.Payload)
		handleErr(r, http.StatusBadRequest, "bad request payload")
		return
	}

	if cfg.Options.Debug {
		s.Log.Debugw("handleRun: running with debug")
		s.runDebug(cfg, r)
	} else {
		s.Log.Debugw("handleRun: running without debug")
		s.run(cfg, r)
	}

}

func (s *Server) runDebug(cfg runConfig, r wsRequest) {

	now := time.Now()
	logfile := fmt.Sprintf("./%v.txt", now.Format("2006-01-02-15-04-05"))

	cfg.Options.DebugPaths = []string{logfile}

	result, err := combat.Run(cfg.Config, cfg.Options)
	if err != nil {
		handleErr(r, http.StatusBadRequest, err.Error())
		return
	}

	file, err := os.Open(logfile)
	if err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}
	defer file.Close()
	var log strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log.WriteString(scanner.Text())
		log.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}

	err = file.Close()

	if err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}

	s.Log.Debugw("run complete", "result", result)

	result.Text = result.PrettyPrint()
	result.Debug = log.String()

	data, _ := json.Marshal(result)
	e := wsResponse{
		ID:      r.ID,
		Status:  http.StatusOK,
		Payload: string(data),
	}
	msg, _ := json.Marshal(e)
	r.client.send <- msg

	os.Remove(logfile)
}

func (s *Server) run(cfg runConfig, r wsRequest) {

	result, err := combat.Run(cfg.Config, cfg.Options)
	if err != nil {
		handleErr(r, http.StatusBadRequest, err.Error())
		return
	}

	result.Text = result.PrettyPrint()

	s.Log.Debugw("run complete", "result", result)

	data, _ := json.Marshal(result)
	e := wsResponse{
		ID:      r.ID,
		Status:  http.StatusOK,
		Payload: string(data),
	}
	msg, _ := json.Marshal(e)
	r.client.send <- msg
}

// func (s *Server) runSingle(cfg runConfig, r wsRequest) {
// 	parser := parse.New("single", cfg.Config)
// 	prof, err := parser.Parse()
// 	if err != nil {
// 		handleErr(r, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	now := time.Now()

// 	prof.LogConfig.LogLevel = cfg.LogLvl
// 	prof.LogConfig.LogFile = logfile
// 	prof.FixedRand = cfg.NoSeed

// 	if cfg.HP > 0 {
// 		prof.RunOptions.DamageMode = true
// 		prof.RunOptions.HP = cfg.HP
// 	} else {
// 		prof.RunOptions.FrameLimit = cfg.Seconds * 60
// 		prof.RunOptions.HP = 0
// 	}

// 	sim, err := combat.NewSim(prof)
// 	if err != nil {
// 		handleErr(r, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var stats combat.Stats

// 	stats, err = sim.Run()
// 	if err != nil {
// 		handleErr(r, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	var dur float64 = float64(cfg.Seconds)
// 	if cfg.HP > 0 {
// 		dur = float64(stats.Duration) / 60
// 	}

// 	elapsed := time.Since(now)

// 	//read the log file
// 	file, err := os.Open(logfile)
// 	if err != nil {
// 		handleErr(r, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	defer file.Close()
// 	var log strings.Builder
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		log.WriteString(scanner.Text())
// 		log.WriteString("\n")
// 	}

// 	if err := scanner.Err(); err != nil {
// 		handleErr(r, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	file.Close()

// 	var sb strings.Builder
// 	sb.WriteString("------------------------------------------\n")
// 	for i, t := range stats.DamageByChar {
// 		sb.WriteString(fmt.Sprintf("%v contributed the following dps:\n", stats.CharNames[i]))
// 		keys := make([]string, 0, len(t))
// 		for k := range t {
// 			keys = append(keys, k)
// 		}
// 		sort.Strings(keys)
// 		var total float64
// 		for _, k := range keys {
// 			v := t[k]
// 			sb.WriteString(fmt.Sprintf("\t%v: %.2f (%.2f%%; total = %.0f)\n", k, v/float64(dur), 100*v/stats.Damage, v))
// 			total += v
// 		}

// 		sb.WriteString(fmt.Sprintf("%v total dps: %.2f (dmg: %.2f); total percentage: %.0f%%\n", stats.CharNames[i], total/float64(dur), total, 100*total/stats.Damage))

// 	}
// 	sb.WriteString("------------------------------------------\n")
// 	for char, t := range stats.AbilUsageCountByChar {
// 		sb.WriteString(fmt.Sprintf("%v used the following abilities:\n", stats.CharNames[char]))
// 		keys := make([]string, 0, len(t))
// 		for k := range t {
// 			keys = append(keys, k)
// 		}
// 		sort.Strings(keys)
// 		for _, k := range keys {
// 			v := t[k]
// 			sb.WriteString(fmt.Sprintf("\t%v: %v\n", k, v))
// 		}
// 	}
// 	sb.WriteString("------------------------------------------\n")
// 	for i, v := range stats.CharActiveTime {
// 		sb.WriteString(fmt.Sprintf("%v active for %v (%v seconds - %.0f%%)\n", stats.CharNames[i], v, v/60, 100*float64(v)/float64(dur*60)))
// 	}
// 	sb.WriteString("------------------------------------------\n")
// 	rk := make([]core.ReactionType, 0, len(stats.ReactionsTriggered))
// 	for k := range stats.ReactionsTriggered {
// 		rk = append(rk, k)
// 	}
// 	for _, k := range rk {
// 		v := stats.ReactionsTriggered[k]
// 		sb.WriteString(fmt.Sprintf("%v: %v\n", k, v))
// 	}
// 	sb.WriteString("------------------------------------------\n")
// 	sb.WriteString(fmt.Sprintf("Running profile %v, total damage dealt: %.2f over %v seconds. DPS = %.2f. Sim took %s\n", prof.Label, stats.Damage, dur, stats.DPS, elapsed))

// 	var result struct {
// 		Names   []string     `json:"names"`
// 		Summary string       `json:"summary"`
// 		Log     string       `json:"log"`
// 		Details combat.Stats `json:"details"`
// 	}
// 	result.Summary = sb.String()
// 	result.Log = log.String()
// 	result.Names = stats.CharNames
// 	result.Details = stats

// 	data, _ := json.Marshal(result)
// 	e := wsResponse{
// 		ID:      r.ID,
// 		Status:  http.StatusOK,
// 		Payload: string(data),
// 	}
// 	msg, _ := json.Marshal(e)
// 	r.client.send <- msg

// 	os.Remove(logfile)
// }

// func (s *Server) runAvg(cfg runConfig, r wsRequest) {
// 	w := runtime.NumCPU() * 3
// 	start := time.Now()
// 	stats, err := runDetailedIter(cfg.Iter, w, cfg.Config, cfg.HP, cfg.Seconds)
// 	if err != nil {
// 		handleErr(r, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	stats.Iter = cfg.Iter
// 	elapsed := time.Since(start)

// 	var sb strings.Builder

// 	sb.WriteString("------------------------------------------\n")
// 	for i, t := range stats.DamageByChar {
// 		sb.WriteString(fmt.Sprintf("%v contributed the following dps:\n", stats.CharNames[i]))
// 		keys := make([]string, 0, len(t))
// 		for k := range t {
// 			keys = append(keys, k)
// 		}
// 		sort.Strings(keys)
// 		var total float64
// 		for _, k := range keys {
// 			v := t[k]
// 			sb.WriteString(fmt.Sprintf("\t%v (%.2f%% of total): avg %.2f [min: %.2f | max: %.2f] \n", k, 100*v.Mean/stats.DPS.Mean, v.Mean, v.Min, v.Max))
// 			total += v.Mean
// 		}

// 		sb.WriteString(fmt.Sprintf("%v total avg dps: %.2f; total percentage: %.0f%%\n", stats.CharNames[i], total, 100*total/stats.DPS.Mean))
// 	}
// 	sb.WriteString("------------------------------------------\n")
// 	for i, t := range stats.AbilUsageCountByChar {
// 		sb.WriteString(fmt.Sprintf("%v used the following abilities:\n", stats.CharNames[i]))
// 		keys := make([]string, 0, len(t))
// 		for k := range t {
// 			keys = append(keys, k)
// 		}
// 		sort.Strings(keys)
// 		for _, k := range keys {
// 			v := t[k]
// 			sb.WriteString(fmt.Sprintf("\t%v: avg %.2f [min: %v | max: %v]\n", k, v.Mean, v.Min, v.Max))
// 		}
// 	}
// 	sb.WriteString("------------------------------------------\n")
// 	for i, v := range stats.CharActiveTime {
// 		sb.WriteString(fmt.Sprintf("%v on average active for %.0f%% [min: %.0f%% | max: %.0f%%]\n", stats.CharNames[i], 100*v.Mean/(stats.AvgDuration*60), float64(100*v.Min)/(stats.AvgDuration*60), float64(100*v.Max)/(stats.AvgDuration*60)))
// 	}
// 	sb.WriteString("------------------------------------------\n")
// 	rk := make([]core.ReactionType, 0, len(stats.ReactionsTriggered))
// 	for k := range stats.ReactionsTriggered {
// 		rk = append(rk, k)
// 	}
// 	for _, k := range rk {
// 		v := stats.ReactionsTriggered[k]
// 		sb.WriteString(fmt.Sprintf("\t%v: avg %.2f [min: %v max: %v]\n", k, v.Mean, v.Min, v.Max))
// 	}
// 	sb.WriteString(fmt.Sprintf("Simulation done in %s; %v iterations; average %.0f dps over %v seconds (min: %.2f max: %.2f std: %.2f) \n", elapsed, cfg.Iter, stats.DPS.Mean, stats.AvgDuration, stats.DPS.Min, stats.DPS.Max, stats.DPS.SD))

// 	var resp struct {
// 		Summary string  `json:"summary"`
// 		Details Summary `json:"details"`
// 	}
// 	resp.Summary = sb.String()
// 	resp.Details = stats

// 	data, _ := json.Marshal(resp)
// 	e := wsResponse{
// 		ID:      r.ID,
// 		Status:  http.StatusOK,
// 		Payload: string(data),
// 	}
// 	msg, _ := json.Marshal(e)
// 	r.client.send <- msg
// }

// func runDetailedIter(n, w int, src string, hp float64, dur int) (Summary, error) {
// 	// var progress float64
// 	var data []combat.Stats
// 	var s Summary

// 	//parse the config just so we know how many chars there are; and also abort early if needed
// 	parser := parse.New("single", src)
// 	cfg, err := parser.Parse()
// 	if err != nil {
// 		return s, err
// 	}

// 	charCount := len(cfg.Characters.Profile)

// 	if charCount > 4 {
// 		return s, errors.New("cannot have more than 4 characters in a team")
// 	}

// 	s.DPS.Min = math.MaxFloat64
// 	s.DPS.Max = -1
// 	s.ReactionsTriggered = make(map[core.ReactionType]ResultInt)
// 	s.CharNames = make([]string, charCount)
// 	s.AbilUsageCountByChar = make([]map[string]ResultInt, charCount)
// 	s.CharActiveTime = make([]ResultInt, charCount)
// 	s.DamageByChar = make([]map[string]ResultFloat, charCount)

// 	for i := range s.CharNames {
// 		s.CharNames[i] = cfg.Characters.Profile[i].Base.Name
// 		s.CharActiveTime[i].Min = math.MaxInt64
// 		s.CharActiveTime[i].Max = -1
// 		s.AbilUsageCountByChar[i] = make(map[string]ResultInt)
// 		s.DamageByChar[i] = make(map[string]ResultFloat)
// 	}

// 	count := n

// 	resp := make(chan workerResp, n)
// 	req := make(chan bool)
// 	done := make(chan bool)

// 	for i := 0; i < w; i++ {
// 		go detailedWorker(src, hp, dur, resp, req, done)
// 	}

// 	go func() {
// 		var wip int
// 		for wip < n {
// 			req <- true
// 			wip++
// 		}
// 	}()

// 	dd := float64(dur)
// 	if hp == 0 {
// 		s.AvgDuration = dd
// 	}

// 	//initialize some vars in summary

// 	for count > 0 {
// 		z := <-resp
// 		//safely abort
// 		if z.err != nil {
// 			close(done)
// 			return s, z.err
// 		}
// 		v := z.stats
// 		count--
// 		data = append(data, v)

// 		// log.Println(v)

// 		if hp > 0 {
// 			dd = float64(v.Duration) / 60.0
// 			s.AvgDuration += dd / float64(n)
// 		}

// 		//dps
// 		if s.DPS.Min > v.DPS {
// 			s.DPS.Min = v.DPS
// 		}
// 		if s.DPS.Max < v.DPS {
// 			s.DPS.Max = v.DPS
// 		}
// 		s.DPS.Mean += v.DPS / float64(n)

// 		//char active time
// 		for i, amt := range v.CharActiveTime {

// 			if s.CharActiveTime[i].Min > amt {
// 				s.CharActiveTime[i].Min = amt
// 			}
// 			if s.CharActiveTime[i].Max < amt {
// 				s.CharActiveTime[i].Max = amt
// 			}
// 			s.CharActiveTime[i].Mean += float64(amt) / float64(n)

// 		}

// 		//dmg by char
// 		for i, abil := range v.DamageByChar {
// 			for k, amt := range abil {
// 				x, ok := s.DamageByChar[i][k]
// 				if !ok {
// 					x.Min = math.MaxFloat64
// 					x.Max = -1
// 				}
// 				amt = amt / float64(dd)
// 				if x.Min > amt {
// 					x.Min = amt
// 				}
// 				if x.Max < amt {
// 					x.Max = amt
// 				}
// 				x.Mean += amt / float64(n)

// 				s.DamageByChar[i][k] = x
// 			}
// 		}

// 		//abil use
// 		for c, abil := range v.AbilUsageCountByChar {
// 			for k, amt := range abil {
// 				x, ok := s.AbilUsageCountByChar[c][k]
// 				if !ok {
// 					x.Min = math.MaxInt64
// 					x.Max = -1
// 				}
// 				if x.Min > amt {
// 					x.Min = amt
// 				}
// 				if x.Max < amt {
// 					x.Max = amt
// 				}
// 				x.Mean += float64(amt) / float64(n)

// 				s.AbilUsageCountByChar[c][k] = x
// 			}
// 		}

// 		//reactions
// 		for c, amt := range v.ReactionsTriggered {
// 			x, ok := s.ReactionsTriggered[c]
// 			if !ok {
// 				x.Min = math.MaxInt64
// 				x.Max = -1
// 			}
// 			if x.Min > amt {
// 				x.Min = amt
// 			}
// 			if x.Max < amt {
// 				x.Max = amt
// 			}
// 			x.Mean += float64(amt) / float64(n)

// 			s.ReactionsTriggered[c] = x
// 		}

// 	}

// 	close(done)

// 	//calculate std dev

// 	for _, v := range data {
// 		s.DPS.SD += (v.DPS - s.DPS.Mean) * (v.DPS - s.DPS.Mean)
// 	}

// 	s.DPS.SD = math.Sqrt(s.DPS.SD / float64(n))

// 	//calculate variances

// 	return s, nil
// }

// func detailedWorker(src string, hp float64, dur int, resp chan workerResp, req chan bool, done chan bool) {

// 	for {
// 		select {
// 		case <-req:
// 			parser := parse.New("single", src)
// 			cfg, err := parser.Parse()
// 			if err != nil {
// 				resp <- workerResp{stats: combat.Stats{}, err: err}
// 				return
// 			}
// 			cfg.LogConfig.LogLevel = "error"
// 			cfg.LogConfig.LogFile = ""
// 			cfg.LogConfig.LogShowCaller = false

// 			if hp > 0 {
// 				cfg.RunOptions.DamageMode = true
// 				cfg.RunOptions.HP = hp
// 			} else {
// 				cfg.RunOptions.FrameLimit = dur * 60
// 				cfg.RunOptions.HP = 0
// 			}

// 			s, err := combat.NewSim(cfg)
// 			if err != nil {
// 				resp <- workerResp{stats: combat.Stats{}, err: err}
// 				return
// 			}

// 			stat, err := s.Run()
// 			if err != nil {
// 				resp <- workerResp{stats: combat.Stats{}, err: err}
// 				return
// 			}

// 			resp <- workerResp{stats: stat, err: nil}

// 		case <-done:
// 			return
// 		}
// 	}
// }

// type workerResp struct {
// 	stats combat.Stats
// 	err   error
// }

// type Summary struct {
// 	Iter                 int                             `json:"iter"`
// 	AvgDuration          float64                         `json:"avg_duration"`
// 	DPS                  ResultFloat                     `json:"dps"`
// 	DamageByChar         []map[string]ResultFloat        `json:"damage_by_char"`
// 	CharActiveTime       []ResultInt                     `json:"char_active_time"`
// 	AbilUsageCountByChar []map[string]ResultInt          `json:"abil_usage_count_by_char"`
// 	ReactionsTriggered   map[core.ReactionType]ResultInt `json:"reactions_triggered"`
// 	CharNames            []string                        `json:"char_names"`
// }

// type ResultFloat struct {
// 	Mean float64 `json:"mean"`
// 	Min  float64 `json:"min"`
// 	Max  float64 `json:"max"`
// 	SD   float64 `json:"sd"`
// }

// type ResultInt struct {
// 	Min  int     `json:"min"`
// 	Max  int     `json:"max"`
// 	Mean float64 `json:"mean"`
// }
