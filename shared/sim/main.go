package main

import "C"

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/parse"

	//characters
	_ "github.com/genshinsim/gsim/internal/characters/ayaka"
	_ "github.com/genshinsim/gsim/internal/characters/beidou"
	_ "github.com/genshinsim/gsim/internal/characters/bennett"
	_ "github.com/genshinsim/gsim/internal/characters/diona"
	_ "github.com/genshinsim/gsim/internal/characters/eula"
	_ "github.com/genshinsim/gsim/internal/characters/fischl"
	_ "github.com/genshinsim/gsim/internal/characters/ganyu"
	_ "github.com/genshinsim/gsim/internal/characters/kaeya"
	_ "github.com/genshinsim/gsim/internal/characters/ningguang"
	_ "github.com/genshinsim/gsim/internal/characters/noelle"
	_ "github.com/genshinsim/gsim/internal/characters/raiden"
	_ "github.com/genshinsim/gsim/internal/characters/sucrose"
	_ "github.com/genshinsim/gsim/internal/characters/xiangling"
	_ "github.com/genshinsim/gsim/internal/characters/xingqiu"
	_ "github.com/genshinsim/gsim/internal/characters/yoimiya"

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
	_ "github.com/genshinsim/gsim/internal/weapons/bow/hamayumi"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/prototype"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/rust"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/sharpshooter"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/skyward"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/stringless"
	_ "github.com/genshinsim/gsim/internal/weapons/bow/thundering"
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
	_ "github.com/genshinsim/gsim/internal/weapons/spear/grasscutter"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/homa"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/kitain"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/primordial"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/prototype"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/skyward"
	_ "github.com/genshinsim/gsim/internal/weapons/spear/vortex"

	_ "github.com/genshinsim/gsim/internal/weapons/sword/alley"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/amenoma"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/aquila"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/black"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/festering"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/freedom"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/harbinger"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/ironsting"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/lion"
	_ "github.com/genshinsim/gsim/internal/weapons/sword/mistsplitter"
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

func main() {}

//export Run
func Run(config string) *C.char {

	parser := parse.New("single", config)
	cfg, err := parser.Parse()
	if err != nil {
		return C.CString(errToString("error parsing config"))
	}

	var gostring string
	if cfg.RunOptions.Average {
		gostring = runAvg(cfg, config)
	} else {
		s, err := runSingle(cfg)
		if err != nil {
			return C.CString(errToString(err.Error()))
		}
		s.Mode = "single"

		b, _ := json.Marshal(s)

		gostring = string(b)
	}

	return C.CString(gostring)

}

func errToString(s string) string {
	var r struct {
		Err string `json:"err"`
	}
	r.Err = s

	b, _ := json.Marshal(r)

	return string(b)
}

func runSingle(cfg core.Config) (combat.SimStats, error) {

	if !cfg.RunOptions.DamageMode {
		cfg.RunOptions.FrameLimit = cfg.RunOptions.Duration * 60
	}

	s, err := combat.NewSim(cfg)
	if err != nil {
		return combat.SimStats{}, err
	}

	return s.Run()
}

func runAvg(cfg core.Config, source string) string {
	stats, err := runDetailedIter(cfg, source)
	if err != nil {
		return errToString(err.Error())
	}

	stats.Mode = "average"

	b, _ := json.Marshal(stats)

	return string(b)
}

type result struct {
	Mean float64 `json:"mean"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	SD   float64 `json:"sd"`
}

type resulti struct {
	Min  int     `json:"min"`
	Max  int     `json:"max"`
	Mean float64 `json:"mean"`
}

type workerResp struct {
	stats combat.SimStats
	err   error
}

type sum struct {
	Mode                 string                        `json:"mode"`
	AvgDuration          float64                       `json:"duration"`
	DPS                  result                        `json:"dps"`
	DamageByChar         []map[string]result           `json:"damage_by_char"`
	CharActiveTime       []resulti                     `json:"char_active_time"`
	AbilUsageCountByChar []map[string]resulti          `json:"abil_usage_count_by_char"`
	ReactionsTriggered   map[core.ReactionType]resulti `json:"reactions_triggered"`
	CharNames            []string                      `json:"char_names"`
}

func runDetailedIter(cfg core.Config, source string) (sum, error) {
	// var progress float64
	var data []combat.SimStats
	var summary sum

	charCount := len(cfg.Characters.Profile)

	if charCount > 4 {
		return summary, errors.New("cannot have more than 4 characters in a team")
	}

	summary.DPS.Min = math.MaxFloat64
	summary.DPS.Max = -1
	summary.ReactionsTriggered = make(map[core.ReactionType]resulti)
	summary.CharNames = make([]string, charCount)
	summary.AbilUsageCountByChar = make([]map[string]resulti, charCount)
	summary.CharActiveTime = make([]resulti, charCount)
	summary.DamageByChar = make([]map[string]result, charCount)

	for i := range summary.CharNames {
		summary.CharNames[i] = cfg.Characters.Profile[i].Base.Name
		summary.CharActiveTime[i].Min = math.MaxInt64
		summary.CharActiveTime[i].Max = -1
		summary.AbilUsageCountByChar[i] = make(map[string]resulti)
		summary.DamageByChar[i] = make(map[string]result)
	}

	count := cfg.RunOptions.Iteration
	n := cfg.RunOptions.Iteration

	w := cfg.RunOptions.Workers
	if w <= 0 {
		w = 10
	}

	resp := make(chan workerResp, count)
	req := make(chan bool)
	done := make(chan bool)

	for i := 0; i < w; i++ {
		go detailedWorker(source, resp, req, done)
	}

	go func() {
		var wip int
		for wip < n {
			req <- true
			wip++
		}
	}()

	dd := float64(cfg.RunOptions.Duration)
	if !cfg.RunOptions.DamageMode {
		summary.AvgDuration = dd
	}

	defer close(done)

	for count > 0 {
		vv := <-resp
		if vv.err != nil {
			return summary, vv.err
		}
		v := vv.stats
		count--
		data = append(data, v)

		// log.Println(v)
		//print out progress
		// log.Printf("done %v\n", n-count)

		if cfg.RunOptions.DamageMode {
			dd = float64(v.SimDuration) / 60.0
			summary.AvgDuration += dd / float64(n)
		}

		//dps
		if summary.DPS.Min > v.DPS {
			summary.DPS.Min = v.DPS
		}
		if summary.DPS.Max < v.DPS {
			summary.DPS.Max = v.DPS
		}
		summary.DPS.Mean += v.DPS / float64(n)

		//char active time
		for i, amt := range v.CharActiveTime {

			if summary.CharActiveTime[i].Min > amt {
				summary.CharActiveTime[i].Min = amt
			}
			if summary.CharActiveTime[i].Max < amt {
				summary.CharActiveTime[i].Max = amt
			}
			summary.CharActiveTime[i].Mean += float64(amt) / float64(n)

		}

		//dmg by char
		for i, abil := range v.DamageByChar {
			for k, amt := range abil {
				x, ok := summary.DamageByChar[i][k]
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

				summary.DamageByChar[i][k] = x
			}
		}

		//abil use
		for c, abil := range v.AbilUsageCountByChar {
			for k, amt := range abil {
				x, ok := summary.AbilUsageCountByChar[c][k]
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

				summary.AbilUsageCountByChar[c][k] = x
			}
		}

		//reactions
		for c, amt := range v.ReactionsTriggered {
			x, ok := summary.ReactionsTriggered[c]
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

			summary.ReactionsTriggered[c] = x
		}
	}

	//calculate std dev

	for _, v := range data {
		summary.DPS.SD += (v.DPS - summary.DPS.Mean) * (v.DPS - summary.DPS.Mean)
	}

	summary.DPS.SD = math.Sqrt(summary.DPS.SD / float64(n))

	//calculate variances

	return summary, nil
}

func detailedWorker(src string, resp chan workerResp, req chan bool, done chan bool) {

	for {
		select {
		case <-req:
			parser := parse.New("single", src)
			cfg, _ := parser.Parse()

			cfg.LogConfig.LogLevel = "error"
			cfg.LogConfig.LogFile = ""
			cfg.LogConfig.LogShowCaller = false

			if !cfg.RunOptions.DamageMode {
				cfg.RunOptions.FrameLimit = cfg.RunOptions.Duration * 60
			}

			s, err := combat.NewSim(cfg)
			if err != nil {
				resp <- workerResp{
					err: err,
				}
				return
			}

			// log.Println("starting new job")

			stat, err := s.Run()

			// log.Println("job done")

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
