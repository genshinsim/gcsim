package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/genshinsim/gsim/internal/logtohtml"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
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

func main() {

	var source []byte

	var err error

	debug := flag.String("d", "debug", "output level: debug, info, warn")
	seconds := flag.Int("s", 90, "how many seconds to run the sim for")
	cfgFile := flag.String("p", "config.txt", "which profile to use")
	f := flag.String("o", "debug.log", "detailed log file")
	hp := flag.Float64("hp", 0, "hp mode: how much hp to deal damage to")
	showCaller := flag.Bool("caller", false, "show caller in debug low")
	fixedRand := flag.Bool("noseed", false, "use 0 for rand seed always - guarantee same results every time; only in single mode")
	avgMode := flag.Bool("a", false, "run sim multiple times and calculate avg damage (smooth out randomness). default false. note that there is no debug log in this mode")
	w := flag.Int("w", 24, "number of workers to run when running multiple iterations; default 24")
	i := flag.Int("i", 5000, "number of iterations to run if we're running multiple")
	multi := flag.String("comp", "", "comparison mode")
	t := flag.Int("t", 1, "target multiplier")

	flag.Parse()

	source, err = ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case *multi != "":
		content, err := ioutil.ReadFile(*multi)
		if err != nil {
			log.Fatal(err)
		}
		files := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
		// lines := strings.Split(string(content), `\n`)
		runMulti(*i, *w, files, *hp, *seconds, *t)
	case *avgMode:
		runAvg(*i, *w, source, *hp, *seconds, *t)
	default:
		// defer profile.Start(profile.ProfilePath("./")).Stop()
		parser := parse.New("single", string(source))
		cfg, err := parser.Parse()
		if err != nil {
			log.Fatal(err)
		}
		cfg.LogConfig.LogLevel = *debug
		cfg.LogConfig.LogFile = *f
		cfg.LogConfig.LogShowCaller = *showCaller
		cfg.FixedRand = *fixedRand

		//make it all true for now
		os.Remove(*f)
		runSingle(cfg, *hp, *seconds)
	}

}

func runSingle(cfg def.Config, hp float64, dur int) {

	if hp > 0 {
		cfg.Mode.HPMode = true
		cfg.Mode.HP = hp
	} else {
		cfg.Mode.FrameLimit = dur * 60
		cfg.Mode.HP = 0
	}
	log.Println(cfg.Mode.FrameLimit)

	s, err := combat.NewSim(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var stats combat.SimStats
	var dmg float64

	start := time.Now()
	stats, err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
	dur = stats.SimDuration / 60
	dmg = stats.Damage

	elapsed := time.Since(start)
	fmt.Println("------------------------------------------")
	for i, t := range stats.DamageByChar {
		fmt.Printf("%v contributed the following dps:\n", stats.CharNames[i])
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var total float64
		for _, k := range keys {
			v := t[k]
			fmt.Printf("\t%v (%.2f%% of total): %.2f dps \n", k, 100*v/dmg, v*60/float64(stats.SimDuration))
			total += v
		}

		fmt.Printf("%v total dps: %.2f; total percentage: %.0f%%\n", stats.CharNames[i], total/float64(stats.SimDuration), 100*total/dmg)
	}
	fmt.Println("------------------------------------------")
	for i, t := range stats.AbilUsageCountByChar {
		fmt.Printf("%v used the following abilities:\n", stats.CharNames[i])
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t[k]
			fmt.Printf("\t%v: %v\n", k, v)
		}
	}
	fmt.Println("------------------------------------------")
	for i, v := range stats.CharActiveTime {
		fmt.Printf("%v active for %v (%v seconds - %.0f%%)\n", stats.CharNames[i], v, v/60, 100*float64(v)/float64(dur*60))
	}
	fmt.Println("------------------------------------------")
	rk := make([]def.ReactionType, 0, len(stats.ReactionsTriggered))
	for k := range stats.ReactionsTriggered {
		rk = append(rk, k)
	}
	for _, k := range rk {
		v := stats.ReactionsTriggered[k]
		fmt.Printf("%v: %v\n", k, v)
	}
	fmt.Println("------------------------------------------")
	fmt.Printf("Running profile %v, total damage dealt: %.2f over %v seconds. DPS = %.2f. Sim took %s\n", cfg.Label, stats.Damage, dur, stats.DPS, elapsed)

	if cfg.LogConfig.LogFile != "" {
		//make excel log

		err = logtohtml.Write(cfg.LogConfig.LogFile, "./debug.html", cfg.Characters.Initial, stats.CharNames)
		if err != nil {
			log.Fatal(err)
		}

		// //read the log file
		// file, err := os.Open(cfg.LogConfig.LogFile)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer file.Close()
		// var logs strings.Builder
		// scanner := bufio.NewScanner(file)
		// for scanner.Scan() {
		// 	logs.WriteString(scanner.Text())
		// 	logs.WriteString("\n")
		// }

		// if err := scanner.Err(); err != nil {
		// 	log.Fatal(err)
		// }

		// file.Close()

		// t, err := template.ParseFiles("./log.tmpl")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// var tConfig struct {
		// 	Active string
		// 	Team   string
		// 	Log    string
		// }
		// tConfig.Active = cfg.Characters.Initial
		// tConfig.Team = fmt.Sprint(strings.Join(strings.Split(fmt.Sprintf("%+q", stats.CharNames), " "), ", "))
		// tConfig.Log = logs.String()

		// os.Remove("./debug.html")
		// f, err := os.Create("debug.html")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer f.Close()
		// err = t.Execute(f, tConfig)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// os.Remove(cfg.LogConfig.LogFile)
	}

}

func runAvg(n, w int, src []byte, hp float64, dur int, t int) {
	start := time.Now()
	stats := runDetailedIter(n, w, src, hp, dur)
	elapsed := time.Since(start)
	fmt.Println("------------------------------------------")
	for i, t := range stats.DamageByChar {
		fmt.Printf("%v contributed the following dps:\n", stats.CharNames[i])
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var total float64
		for _, k := range keys {
			v := t[k]
			fmt.Printf("\t%v (%.2f%% of total): avg %.2f [min: %.2f | max: %.2f] \n", k, 100*v.mean/stats.dps.mean, v.mean, v.min, v.max)
			total += v.mean
		}

		fmt.Printf("%v total avg dps: %.2f; total percentage: %.0f%%\n", stats.CharNames[i], total, 100*total/stats.dps.mean)
	}
	fmt.Println("------------------------------------------")
	for i, t := range stats.AbilUsageCountByChar {
		fmt.Printf("%v used the following abilities:\n", stats.CharNames[i])
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t[k]
			fmt.Printf("\t%v: avg %.2f [min: %v | max: %v]\n", k, v.mean, v.min, v.max)
		}
	}
	fmt.Println("------------------------------------------")
	for i, v := range stats.CharActiveTime {
		fmt.Printf("%v on average active for %.0f%% [min: %.0f%% | max: %.0f%%]\n", stats.CharNames[i], 100*v.mean/(stats.avgdur*60), float64(100*v.min)/(stats.avgdur*60), float64(100*v.max)/(stats.avgdur*60))
	}
	fmt.Println("------------------------------------------")
	rk := make([]def.ReactionType, 0, len(stats.ReactionsTriggered))
	for k := range stats.ReactionsTriggered {
		rk = append(rk, k)
	}
	for _, k := range rk {
		v := stats.ReactionsTriggered[k]
		fmt.Printf("\t%v: avg %.2f [min: %v max: %v]\n", k, v.mean, v.min, v.max)
	}
	fmt.Printf("Simulation done in %s; %v iterations; average %.0f dps over %v seconds (min: %.2f max: %.2f std: %.2f) \n", elapsed, n, stats.dps.mean, stats.avgdur, stats.dps.min, stats.dps.max, stats.dps.sd)
}

type result struct {
	mean float64
	min  float64
	max  float64
	sd   float64
}

func runDetailedIter(n, w int, src []byte, hp float64, dur int) sum {
	// var progress float64
	var data []combat.SimStats
	var summary sum

	parser := parse.New("single", string(src))
	cfg, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	charCount := len(cfg.Characters.Profile)

	if charCount > 4 {
		panic(errors.New("cannot have more than 4 characters in a team"))
	}

	summary.dps.min = math.MaxFloat64
	summary.dps.max = -1
	summary.ReactionsTriggered = make(map[def.ReactionType]resulti)
	summary.CharNames = make([]string, charCount)
	summary.AbilUsageCountByChar = make([]map[string]resulti, charCount)
	summary.CharActiveTime = make([]resulti, charCount)
	summary.DamageByChar = make([]map[string]result, charCount)

	for i := range summary.CharNames {
		summary.CharNames[i] = cfg.Characters.Profile[i].Base.Name
		summary.CharActiveTime[i].min = math.MaxInt64
		summary.CharActiveTime[i].max = -1
		summary.AbilUsageCountByChar[i] = make(map[string]resulti)
		summary.DamageByChar[i] = make(map[string]result)
	}

	count := n

	resp := make(chan combat.SimStats, n)
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
		summary.avgdur = dd
	}

	for count > 0 {
		v := <-resp
		count--
		data = append(data, v)

		// log.Println(v)
		//print out progress
		// log.Printf("done %v\n", n-count)

		if hp > 0 {
			dd = float64(v.SimDuration) / 60.0
			summary.avgdur += dd / float64(n)
		}

		//dps
		if summary.dps.min > v.DPS {
			summary.dps.min = v.DPS
		}
		if summary.dps.max < v.DPS {
			summary.dps.max = v.DPS
		}
		summary.dps.mean += v.DPS / float64(n)

		//char active time
		for i, amt := range v.CharActiveTime {

			if summary.CharActiveTime[i].min > amt {
				summary.CharActiveTime[i].min = amt
			}
			if summary.CharActiveTime[i].max < amt {
				summary.CharActiveTime[i].max = amt
			}
			summary.CharActiveTime[i].mean += float64(amt) / float64(n)

		}

		//dmg by char
		for i, abil := range v.DamageByChar {
			for k, amt := range abil {
				x, ok := summary.DamageByChar[i][k]
				if !ok {
					x.min = math.MaxFloat64
					x.max = -1
				}
				amt = amt / float64(dd)
				if x.min > amt {
					x.min = amt
				}
				if x.max < amt {
					x.max = amt
				}
				x.mean += amt / float64(n)

				summary.DamageByChar[i][k] = x
			}
		}

		//abil use
		for c, abil := range v.AbilUsageCountByChar {
			for k, amt := range abil {
				x, ok := summary.AbilUsageCountByChar[c][k]
				if !ok {
					x.min = math.MaxInt64
					x.max = -1
				}
				if x.min > amt {
					x.min = amt
				}
				if x.max < amt {
					x.max = amt
				}
				x.mean += float64(amt) / float64(n)

				summary.AbilUsageCountByChar[c][k] = x
			}
		}

		//reactions
		for c, amt := range v.ReactionsTriggered {
			x, ok := summary.ReactionsTriggered[c]
			if !ok {
				x.min = math.MaxInt64
				x.max = -1
			}
			if x.min > amt {
				x.min = amt
			}
			if x.max < amt {
				x.max = amt
			}
			x.mean += float64(amt) / float64(n)

			summary.ReactionsTriggered[c] = x
		}
	}

	close(done)

	//calculate std dev

	for _, v := range data {
		summary.dps.sd += (v.DPS - summary.dps.mean) * (v.DPS - summary.dps.mean)
	}

	summary.dps.sd = math.Sqrt(summary.dps.sd / float64(n))

	//calculate variances

	return summary
}

func detailedWorker(src []byte, hp float64, dur int, resp chan combat.SimStats, req chan bool, done chan bool) {

	for {
		select {
		case <-req:
			parser := parse.New("single", string(src))
			cfg, err := parser.Parse()
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}

			// log.Println("starting new job")

			stat, err := s.Run()

			// log.Println("job done")

			if err != nil {
				log.Fatal(err)
			}

			resp <- stat
		case <-done:
			return
		}
	}
}

type sum struct {
	avgdur               float64
	dps                  result
	DamageByChar         []map[string]result
	CharActiveTime       []resulti
	AbilUsageCountByChar []map[string]resulti
	ReactionsTriggered   map[def.ReactionType]resulti
	CharNames            []string
}

type resulti struct {
	min  int
	max  int
	mean float64
}

func runMulti(n, w int, files []string, hp float64, dur int, t int) {
	fmt.Printf("Simulating %v seconds of combat over %v iterations\n", dur, n)
	start := time.Now()
	fmt.Print("Filename                                 |      Mean|       Min|       Max|   Std Dev|\n")
	fmt.Print("--------------------------------------------------------------------------------------\n")
	for _, f := range files {
		if f == "" || f[0] == '#' {
			continue
		}
		source, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%40.40v |", f)
		r := runIter(n, w, source, hp, dur, t)
		fmt.Printf("%10.2f|%10.2f|%10.2f|%10.2f|\n", r.mean, r.min, r.max, r.sd)
	}
	elapsed := time.Since(start)
	fmt.Printf("Completed in %s\n", elapsed)
}

func runIter(n, w int, src []byte, hp float64, dur int, t int) result {
	// var progress float64
	var sum, ss, min, max float64
	var data []float64
	min = math.MaxFloat64
	max = -1

	count := n

	resp := make(chan float64, n)
	req := make(chan bool)
	done := make(chan bool)

	for i := 0; i < w; i++ {
		go worker(src, hp, dur, t, resp, req, done)
	}

	go func() {
		var wip int
		for wip < n {
			req <- true
			wip++
		}
	}()

	for count > 0 {
		val := <-resp
		count--
		data = append(data, val)
		sum += val
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}

	}

	close(done)

	mean := sum / float64(n)

	for _, v := range data {
		ss += (v - mean) * (v - mean)
	}

	sd := math.Sqrt(ss / float64(n))

	return result{
		mean: mean,
		min:  min,
		max:  max,
		sd:   sd,
	}
}

func worker(src []byte, hp float64, dur int, t int, resp chan float64, req chan bool, done chan bool) {

	for {
		select {
		case <-req:
			parser := parse.New("single", string(src))
			cfg, err := parser.Parse()
			if err != nil {
				log.Fatal(err)
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

			var additional []def.EnemyProfile

			for i := 1; i < t; i++ {
				for _, e := range cfg.Targets {
					additional = append(additional, def.CloneEnemy(e))
				}
			}

			cfg.Targets = append(cfg.Targets, additional...)

			s, err := combat.NewSim(cfg)
			if err != nil {
				log.Fatal(err)
			}

			stat, err := s.Run()
			if err != nil {
				log.Fatal(err)
			}

			resp <- stat.DPS

		case <-done:
			return
		}
	}
}
