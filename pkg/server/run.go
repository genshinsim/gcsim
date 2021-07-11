package server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"github.com/genshinsim/gsim/pkg/parse"

	//characters
	_ "github.com/genshinsim/gsim/internal/characters/beidou"
	_ "github.com/genshinsim/gsim/internal/characters/bennett"
	_ "github.com/genshinsim/gsim/internal/characters/fischl"
	_ "github.com/genshinsim/gsim/internal/characters/ganyu"
	_ "github.com/genshinsim/gsim/internal/characters/kaeya"
	_ "github.com/genshinsim/gsim/internal/characters/ningguang"
	_ "github.com/genshinsim/gsim/internal/characters/noelle"
	_ "github.com/genshinsim/gsim/internal/characters/sucrose"
	_ "github.com/genshinsim/gsim/internal/characters/xiangling"
	_ "github.com/genshinsim/gsim/internal/characters/xingqiu"

	//weapons
	_ "github.com/genshinsim/gsim/internal/weapons/common/blackcliff"
	_ "github.com/genshinsim/gsim/internal/weapons/common/favonius"
	_ "github.com/genshinsim/gsim/internal/weapons/common/generic"
	_ "github.com/genshinsim/gsim/internal/weapons/common/lithic"
	_ "github.com/genshinsim/gsim/internal/weapons/common/royal"
	_ "github.com/genshinsim/gsim/internal/weapons/common/sacrificial"

	_ "github.com/genshinsim/gsim/internal/weapons/bow/alley"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/amos"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/elegy"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/prototype"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/rust"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/sharpshooter"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/skyward"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/stringless"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/viridescent"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/windblume"

	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/dodoco"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/frostbearer"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/magicguide"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/mappa"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/memory"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/perception"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/prayer"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/prototype"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/skyward"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/solar"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/thrilling"
	_ "github.com/genshinsim/gsim/internal/weapons/catalyst/widsith"

	_ "github.com/genshinsim/gsim/internal/weapons/claymore/bell"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/pines"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/prototype"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/rainslasher"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/skyrider"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/skyward"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/spine"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/starsilver"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/unforged"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/whiteblind"
	_ "github.com/genshinsim/gsim/internal/weapons/claymore/wolf"

	_ "github.com/genshinsim/gsim/internal/weapons/spear/crescent"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/deathmatch"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/dragonbane"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/dragonspine"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/homa"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/primordial"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/prototype"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/skyward"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/vortex"

	_ "github.com/genshinsim/gsim/internal/weapons/sword/alley"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/aquila"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/black"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/festering"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/freedom"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/harbinger"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/ironsting"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/lion"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/primordial"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/prototype"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/skyrider"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/skyward"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/summit"

	//artifacts
	_ "github.com/genshinsim/gsim/internal/artifacts/archaic"
	_ "github.com/genshinsim/gsim/internal/artifacts/blizzard"
	_ "github.com/genshinsim/gsim/internal/artifacts/bloodstained"
	_ "github.com/genshinsim/gsim/internal/artifacts/bolide"
	_ "github.com/genshinsim/gsim/internal/artifacts/crimson"
	_ "github.com/genshinsim/gsim/internal/artifacts/gladiator"
	_ "github.com/genshinsim/gsim/internal/artifacts/heartofdepth"
	_ "github.com/genshinsim/gsim/internal/artifacts/lavawalker"
	_ "github.com/genshinsim/gsim/internal/artifacts/maiden"
	_ "github.com/genshinsim/gsim/internal/artifacts/noblesse"
	_ "github.com/genshinsim/gsim/internal/artifacts/paleflame"
	_ "github.com/genshinsim/gsim/internal/artifacts/reminiscence"
	_ "github.com/genshinsim/gsim/internal/artifacts/seal"
	_ "github.com/genshinsim/gsim/internal/artifacts/tenacity"
	_ "github.com/genshinsim/gsim/internal/artifacts/thunderingfury"
	_ "github.com/genshinsim/gsim/internal/artifacts/viridescent"
	_ "github.com/genshinsim/gsim/internal/artifacts/wanderer"
)

type runConfig struct {
	LogLvl    string   `json:"log"`
	Seconds   int      `json:"seconds"`
	Config    string   `json:"config"`
	HP        float64  `json:"hp"`
	AvgMode   bool     `json:"avg_mode"`
	Iter      int      `json:"iter"`
	NoSeed    bool     `json:"noseed"`
	LogEvents []string `json:"logs"`
}

func (s *Server) handleRun(ctx context.Context, r wsRequest) {

	var cfg runConfig
	err := json.Unmarshal([]byte(r.Payload), &cfg)

	if err != nil {
		s.Log.Debugw("handleRun: invalid request payload", "payload", r.Payload)
		handleErr(r, http.StatusBadRequest, "bad request payload")
		return
	}

	if cfg.AvgMode {
		s.runAvg(cfg, r)
	} else {
		s.runSingle(cfg, r)
	}

}

func (s *Server) runSingle(cfg runConfig, r wsRequest) {
	parser := parse.New("single", cfg.Config)
	prof, err := parser.Parse()
	if err != nil {
		handleErr(r, http.StatusBadRequest, err.Error())
		return
	}
	now := time.Now()
	logfile := fmt.Sprintf("./%v.txt", now.Format("2006-01-02-15-04-05"))
	prof.LogConfig.LogLevel = cfg.LogLvl
	prof.LogConfig.LogFile = logfile
	prof.FixedRand = cfg.NoSeed

	if cfg.HP > 0 {
		prof.Mode.HPMode = true
		prof.Mode.HP = cfg.HP
	} else {
		prof.Mode.FrameLimit = cfg.Seconds * 60
		prof.Mode.HP = 0
	}

	sim, err := combat.NewSim(prof)
	if err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}

	var stats combat.SimStats

	stats, err = sim.Run()
	if err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}
	var dur float64 = float64(cfg.Seconds)
	if cfg.HP > 0 {
		dur = float64(stats.SimDuration) / 60
	}

	elapsed := time.Since(now)

	//read the log file
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

	file.Close()

	var sb strings.Builder
	sb.WriteString("------------------------------------------\n")
	for i, t := range stats.DamageByChar {
		sb.WriteString(fmt.Sprintf("%v contributed the following dps:\n", stats.CharNames[i]))
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var total float64
		for _, k := range keys {
			v := t[k]
			sb.WriteString(fmt.Sprintf("\t%v: %.2f (%.2f%%; total = %.0f)\n", k, v/float64(dur), 100*v/stats.Damage, v))
			total += v
		}

		sb.WriteString(fmt.Sprintf("%v total dps: %.2f (dmg: %.2f); total percentage: %.0f%%\n", stats.CharNames[i], total/float64(dur), total, 100*total/stats.Damage))

	}
	sb.WriteString("------------------------------------------\n")
	for char, t := range stats.AbilUsageCountByChar {
		sb.WriteString(fmt.Sprintf("%v used the following abilities:\n", stats.CharNames[char]))
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t[k]
			sb.WriteString(fmt.Sprintf("\t%v: %v\n", k, v))
		}
	}
	sb.WriteString("------------------------------------------\n")
	for i, v := range stats.CharActiveTime {
		sb.WriteString(fmt.Sprintf("%v active for %v (%v seconds - %.0f%%)\n", stats.CharNames[i], v, v/60, 100*float64(v)/float64(dur*60)))
	}
	sb.WriteString("------------------------------------------\n")
	rk := make([]def.ReactionType, 0, len(stats.ReactionsTriggered))
	for k := range stats.ReactionsTriggered {
		rk = append(rk, k)
	}
	for _, k := range rk {
		v := stats.ReactionsTriggered[k]
		sb.WriteString(fmt.Sprintf("%v: %v\n", k, v))
	}
	sb.WriteString("------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Running profile %v, total damage dealt: %.2f over %v seconds. DPS = %.2f. Sim took %s\n", prof.Label, stats.Damage, dur, stats.DPS, elapsed))

	var result struct {
		Names   []string        `json:"names"`
		Summary string          `json:"summary"`
		Log     string          `json:"log"`
		Details combat.SimStats `json:"details"`
	}
	result.Summary = sb.String()
	result.Log = log.String()
	result.Names = stats.CharNames
	result.Details = stats

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

func (s *Server) runAvg(cfg runConfig, r wsRequest) {
	w := runtime.NumCPU() * 3
	start := time.Now()
	stats, err := runDetailedIter(cfg.Iter, w, cfg.Config, cfg.HP, cfg.Seconds, r)
	if err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}
	stats.Iter = cfg.Iter
	elapsed := time.Since(start)

	var sb strings.Builder

	sb.WriteString("------------------------------------------\n")
	for i, t := range stats.DamageByChar {
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
	sb.WriteString("------------------------------------------\n")
	for i, t := range stats.AbilUsageCountByChar {
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
	sb.WriteString("------------------------------------------\n")
	for i, v := range stats.CharActiveTime {
		sb.WriteString(fmt.Sprintf("%v on average active for %.0f%% [min: %.0f%% | max: %.0f%%]\n", stats.CharNames[i], 100*v.Mean/(stats.AvgDuration*60), float64(100*v.Min)/(stats.AvgDuration*60), float64(100*v.Max)/(stats.AvgDuration*60)))
	}
	sb.WriteString("------------------------------------------\n")
	rk := make([]def.ReactionType, 0, len(stats.ReactionsTriggered))
	for k := range stats.ReactionsTriggered {
		rk = append(rk, k)
	}
	for _, k := range rk {
		v := stats.ReactionsTriggered[k]
		sb.WriteString(fmt.Sprintf("\t%v: avg %.2f [min: %v max: %v]\n", k, v.Mean, v.Min, v.Max))
	}
	sb.WriteString(fmt.Sprintf("Simulation done in %s; %v iterations; average %.0f dps over %v seconds (min: %.2f max: %.2f std: %.2f) \n", elapsed, cfg.Iter, stats.DPS.Mean, stats.AvgDuration, stats.DPS.Min, stats.DPS.Max, stats.DPS.SD))

	var resp struct {
		Summary string  `json:"summary"`
		Details Summary `json:"details"`
	}
	resp.Summary = sb.String()
	resp.Details = stats

	data, _ := json.Marshal(resp)
	e := wsResponse{
		ID:      r.ID,
		Status:  http.StatusOK,
		Payload: string(data),
	}
	msg, _ := json.Marshal(e)
	r.client.send <- msg
}

func runDetailedIter(n, w int, src string, hp float64, dur int, r wsRequest) (Summary, error) {
	// var progress float64
	var data []combat.SimStats
	var s Summary

	//parse the config just so we know how many chars there are; and also abort early if needed
	parser := parse.New("single", src)
	cfg, err := parser.Parse()
	if err != nil {
		return s, err
	}

	charCount := len(cfg.Characters.Profile)

	if charCount > 4 {
		return s, errors.New("cannot have more than 4 characters in a team")
	}

	s.DPS.Min = math.MaxFloat64
	s.DPS.Max = -1
	s.ReactionsTriggered = make(map[def.ReactionType]ResultInt)
	s.CharNames = make([]string, charCount)
	s.AbilUsageCountByChar = make([]map[string]ResultInt, charCount)
	s.CharActiveTime = make([]ResultInt, charCount)
	s.DamageByChar = make([]map[string]ResultFloat, charCount)

	for i := range s.CharNames {
		s.CharNames[i] = cfg.Characters.Profile[i].Base.Name
		s.CharActiveTime[i].Min = math.MaxInt64
		s.CharActiveTime[i].Max = -1
		s.AbilUsageCountByChar[i] = make(map[string]ResultInt)
		s.DamageByChar[i] = make(map[string]ResultFloat)
	}

	count := n

	resp := make(chan workerResp, n)
	req := make(chan bool)
	done := make(chan bool)

	for i := 0; i < w; i++ {
		go detailedWorker(src, hp, dur, resp, req, done)
	}

	go func() {
		var wip int
		for wip < n {
			req <- true
			wip++
		}
	}()

	dd := float64(dur)
	if hp == 0 {
		s.AvgDuration = dd
	}

	//initialize some vars in summary

	for count > 0 {
		z := <-resp
		//safely abort
		if z.err != nil {
			close(done)
			return s, z.err
		}
		v := z.stats
		count--
		data = append(data, v)

		// log.Println(v)

		if hp > 0 {
			dd = float64(v.SimDuration) / 60.0
			s.AvgDuration += dd / float64(n)
		}

		//dps
		if s.DPS.Min > v.DPS {
			s.DPS.Min = v.DPS
		}
		if s.DPS.Max < v.DPS {
			s.DPS.Max = v.DPS
		}
		s.DPS.Mean += v.DPS / float64(n)

		//char active time
		for i, amt := range v.CharActiveTime {

			if s.CharActiveTime[i].Min > amt {
				s.CharActiveTime[i].Min = amt
			}
			if s.CharActiveTime[i].Max < amt {
				s.CharActiveTime[i].Max = amt
			}
			s.CharActiveTime[i].Mean += float64(amt) / float64(n)

		}

		//dmg by char
		for i, abil := range v.DamageByChar {
			for k, amt := range abil {
				x, ok := s.DamageByChar[i][k]
				if !ok {
					x.Min = math.MaxFloat64
					x.Max = -1
				}
				amt = amt / float64(dd)
				if x.Min > amt {
					x.Min = amt
				}
				if x.Max < amt {
					x.Max = amt
				}
				x.Mean += amt / float64(n)

				s.DamageByChar[i][k] = x
			}
		}

		//abil use
		for c, abil := range v.AbilUsageCountByChar {
			for k, amt := range abil {
				x, ok := s.AbilUsageCountByChar[c][k]
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

				s.AbilUsageCountByChar[c][k] = x
			}
		}

		//reactions
		for c, amt := range v.ReactionsTriggered {
			x, ok := s.ReactionsTriggered[c]
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

			s.ReactionsTriggered[c] = x
		}

	}

	close(done)

	//calculate std dev

	for _, v := range data {
		s.DPS.SD += (v.DPS - s.DPS.Mean) * (v.DPS - s.DPS.Mean)
	}

	s.DPS.SD = math.Sqrt(s.DPS.SD / float64(n))

	//calculate variances

	return s, nil
}

func detailedWorker(src string, hp float64, dur int, resp chan workerResp, req chan bool, done chan bool) {

	for {
		select {
		case <-req:
			parser := parse.New("single", src)
			cfg, err := parser.Parse()
			if err != nil {
				resp <- workerResp{stats: combat.SimStats{}, err: err}
				return
			}
			cfg.LogConfig.LogLevel = "error"
			cfg.LogConfig.LogFile = ""
			cfg.LogConfig.LogShowCaller = false

			if hp > 0 {
				cfg.Mode.HPMode = true
				cfg.Mode.HP = hp
			} else {
				cfg.Mode.FrameLimit = dur * 60
				cfg.Mode.HP = 0
			}

			s, err := combat.NewSim(cfg)
			if err != nil {
				resp <- workerResp{stats: combat.SimStats{}, err: err}
				return
			}

			stat, err := s.Run()
			if err != nil {
				resp <- workerResp{stats: combat.SimStats{}, err: err}
				return
			}

			resp <- workerResp{stats: stat, err: nil}

		case <-done:
			return
		}
	}
}

type workerResp struct {
	stats combat.SimStats
	err   error
}

type Summary struct {
	Iter                 int                            `json:"iter"`
	AvgDuration          float64                        `json:"avg_duration"`
	DPS                  ResultFloat                    `json:"dps"`
	DamageByChar         []map[string]ResultFloat       `json:"damage_by_char"`
	CharActiveTime       []ResultInt                    `json:"char_active_time"`
	AbilUsageCountByChar []map[string]ResultInt         `json:"abil_usage_count_by_char"`
	ReactionsTriggered   map[def.ReactionType]ResultInt `json:"reactions_triggered"`
	CharNames            []string                       `json:"char_names"`
}

type ResultFloat struct {
	Mean float64 `json:"mean"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	SD   float64 `json:"sd"`
}

type ResultInt struct {
	Min  int     `json:"min"`
	Max  int     `json:"max"`
	Mean float64 `json:"mean"`
}
