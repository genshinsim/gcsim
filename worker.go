//
package gcsim

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/montanaflynn/stats"

	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
)

type Stats struct {
	IsDamageMode          bool                      `json:"is_damage_mode"`
	CharNames             []string                  `json:"char_names"`
	DamageByChar          []map[string]float64      `json:"damage_by_char"`
	DamageInstancesByChar []map[string]int          `json:"damage_instances_by_char"`
	DamageByCharByTargets [][]float64               `json:"damage_by_char_by_targets"`
	DamageDetailByTime    map[DamageDetails]float64 `json:"damage_detail_by_time"`
	CharActiveTime        []int                     `json:"char_active_time"`
	AbilUsageCountByChar  []map[string]int          `json:"abil_usage_count_by_char"`
	ParticleCount         map[string]int            `json:"particle_count"`
	ReactionsTriggered    map[core.ReactionType]int `json:"reactions_triggered"`
	Duration              int                       `json:"sim_duration"`
	ElementUptime         []map[core.EleType]int    `json:"ele_uptime"`
	EnergyWhenBurst       [][]float64               `json:"energy_when_burst"`
	//final result
	Damage float64 `json:"damage"`
	DPS    float64 `json:"dps"`
	//for tracking min/max run
	seed int64
}

// Used to track all damage instances for use in damage over time charts
// Currently only the FrameBucket results are exported, but per Char/Target views are likely also going to be useful
type DamageDetails struct {
	FrameBucket int
	Char        int
	Target      int
}

type Result struct {
	IsDamageMode          bool                            `json:"is_damage_mode"`
	ActiveChar            string                          `json:"active_char"`
	CharNames             []string                        `json:"char_names"`
	DamageByChar          []map[string]FloatResult        `json:"damage_by_char"`
	DamageInstancesByChar []map[string]IntResult          `json:"damage_instances_by_char"`
	DamageByCharByTargets []map[int]FloatResult           `json:"damage_by_char_by_targets"`
	CharActiveTime        []IntResult                     `json:"char_active_time"`
	AbilUsageCountByChar  []map[string]IntResult          `json:"abil_usage_count_by_char"`
	ParticleCount         map[string]IntResult            `json:"particle_count"`
	ReactionsTriggered    map[core.ReactionType]IntResult `json:"reactions_triggered"`
	Duration              FloatResult                     `json:"sim_duration"`
	ElementUptime         []map[core.EleType]IntResult    `json:"ele_uptime"`
	RequiredER            []float64                       `json:"required_er"`
	//final result
	Damage         FloatResult            `json:"damage"`
	DPS            FloatResult            `json:"dps"`
	DPSByTarget    map[int]FloatResult    `json:"dps_by_target"`
	DamageOverTime map[string]FloatResult `json:"damage_over_time"`
	Iterations     int                    `json:"iter"`
	Text           string                 `json:"text"`
	Debug          string                 `json:"debug"`
	Runtime        time.Duration          `json:"runtime"`
	//for tracking min/max run
	MinSeed int64 `json:"-"`
	MaxSeed int64 `json:"-"`
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

func Run(src string, opt core.RunOpt, cust ...func(*Simulation) error) (Result, error) {
	start := time.Now()

	//options mode=damage debug=true iteration=5000 duration=90 workers=24;
	var data []Stats

	parser := parse.New("single", string(src))
	cfg, _, err := parser.Parse()
	if err != nil {
		return Result{}, err
	}

	charCount := len(cfg.Characters.Profile)

	if charCount > 4 {
		return Result{}, errors.New("cannot have more than 4 characters in a team")
	}

	chars := make([]string, len(cfg.Characters.Profile))
	for i, v := range cfg.Characters.Profile {
		chars[i] = v.Base.Key.String()
	}

	//set defaults if nothing specified
	count := opt.Iteration
	if count == 0 {
		count = 500
	}
	n := count

	if opt.Debug {
		count--
	}

	// fmt.Printf("running %v iterations\n", count)

	w := opt.Workers
	if w <= 0 {
		w = runtime.NumCPU()
	}
	if w > n {
		w = n
	}

	resp := make(chan workerResp, count)
	req := make(chan int64)
	done := make(chan bool)

	var b [8]byte

	for i := 0; i < w; i++ {
		go worker(src, opt, resp, req, done, cust...)
	}

	go func() {
		var wip int
		for wip < n {
			_, err = crypto_rand.Read(b[:])
			if err != nil {
				log.Panic("cannot seed math/rand package with cryptographically secure random number generator")
			}
			seed := int64(binary.LittleEndian.Uint64(b[:]))
			req <- seed
			wip++
		}
	}()

	defer close(done)

	for count > 0 {
		vv := <-resp
		if vv.err != nil {
			return Result{}, vv.err
		}
		v := vv.stats
		// log.Println(v)
		count--
		data = append(data, v)
	}

	//if debug is true, run one more purely for debug do not add to stats
	if opt.Debug {
		_, err = crypto_rand.Read(b[:])
		if err != nil {
			log.Panic("cannot seed math/rand package with cryptographically secure random number generator")
		}
		_, err = crypto_rand.Read(b[:])
		if err != nil {
			log.Panic("cannot seed math/rand package with cryptographically secure random number generator")
		}
		seed := int64(binary.LittleEndian.Uint64(b[:]))
		s, err := NewSim(cfg, seed, opt, cust...)
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

	result := CollectResult(data, cfg.DamageMode, chars, opt.LogDetails, opt.ERCalcMode)
	result.Iterations = n
	result.ActiveChar = cfg.Characters.Initial.String()
	if !cfg.DamageMode {
		result.Duration.Mean = float64(opt.Duration)
		result.Duration.Min = float64(opt.Duration)
		result.Duration.Max = float64(opt.Duration)
	}
	result.Runtime = time.Since(start)

	return result, nil
}

func CollectResult(data []Stats, mode bool, chars []string, detailed bool, erCalcMode bool) (result Result) {

	n := len(data)

	// TODO: Kind of brittle - maybe track something separate for this?
	targetCount := len(data[0].DamageByCharByTargets[0])
	charCount := len(chars)
	result.DPS.Min = math.MaxFloat64
	result.DPS.Max = -1
	result.DPSByTarget = make(map[int]FloatResult, targetCount)
	result.DamageOverTime = make(map[string]FloatResult)
	if detailed {
		result.ReactionsTriggered = make(map[core.ReactionType]IntResult)
		result.CharNames = make([]string, charCount)
		result.ParticleCount = make(map[string]IntResult)
		result.AbilUsageCountByChar = make([]map[string]IntResult, charCount)
		result.CharActiveTime = make([]IntResult, charCount)
		result.DamageByChar = make([]map[string]FloatResult, charCount)
		result.DamageInstancesByChar = make([]map[string]IntResult, charCount)
		// Defined as a slice of int maps to make code similar to DamageByChar handling
		result.DamageByCharByTargets = make([]map[int]FloatResult, charCount)

		for i := range result.CharNames {
			result.CharNames[i] = chars[i]
			result.CharActiveTime[i].Min = math.MaxInt64
			result.CharActiveTime[i].Max = -1
			result.AbilUsageCountByChar[i] = make(map[string]IntResult)
			result.DamageByChar[i] = make(map[string]FloatResult)
			result.DamageInstancesByChar[i] = make(map[string]IntResult)
			result.DamageByCharByTargets[i] = make(map[int]FloatResult)
		}
	}
	// Used to aggregate individual damage instances across buckets first
	damageOverTimeByRun := make([]map[float64]float64, n)

	// var dd float64

	// Loop through each iteration to build overall statistics
	for iteration, v := range data {
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
			result.MinSeed = v.seed
		}
		if result.DPS.Max < v.DPS {
			result.DPS.Max = v.DPS
			result.MaxSeed = v.seed
		}
		result.DPS.Mean += v.DPS / float64(n)

		if !detailed {
			continue
		}

		damageOverTimeByRun[iteration] = make(map[float64]float64)
		// Damage Over Time - get data for all iterations first and summarized later
		for damageDetails, damage := range v.DamageDetailByTime {
			// Convert frame bucket value into seconds
			secBucket := float64(damageDetails.FrameBucket) / 60.0

			damageOverTimeByRun[iteration][secBucket] += damage
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

		//dmg instances by char
		for i, abil := range v.DamageInstancesByChar {
			for k, amt := range abil {
				x, ok := result.DamageInstancesByChar[i][k]
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

				result.DamageInstancesByChar[i][k] = x
			}
		}

		//dmg by char by target - saved in DPS terms already
		for i, target := range v.DamageByCharByTargets {
			for k, amt := range target {
				x, ok := result.DamageByCharByTargets[i][k]
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

				result.DamageByCharByTargets[i][k] = x
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

		//particles
		for c, amt := range v.ParticleCount {
			x, ok := result.ParticleCount[c]
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

			result.ParticleCount[c] = x
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

		//ele up time
		for t, m := range v.ElementUptime {
			if len(result.ElementUptime) == t {
				result.ElementUptime = append(result.ElementUptime, make(map[core.EleType]IntResult))
			}
			//go through m and add to our results
			for ele, amt := range m {
				x, ok := result.ElementUptime[t][ele]
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

				result.ElementUptime[t][ele] = x
			}
		}
	}

	// Get total DPS by Target stats here
	for i := 0; i < targetCount; i++ {
		for _, char := range result.DamageByCharByTargets {
			dpsResult := result.DPSByTarget[i]
			dpsResult.Mean += char[i].Mean
			result.DPSByTarget[i] = dpsResult
		}
	}

	// Get global mean for each time interval first
	for _, damageData := range damageOverTimeByRun {
		for bucket, dmgTotal := range damageData {
			bucketStr := fmt.Sprintf("%.2f", bucket)
			temp := result.DamageOverTime[bucketStr]
			temp.Mean += dmgTotal / float64(n)
			result.DamageOverTime[bucketStr] = temp
		}
	}

	// Build SD
	for _, damageData := range damageOverTimeByRun {
		for bucket, dmgTotal := range damageData {
			bucketStr := fmt.Sprintf("%.2f", bucket)
			temp := result.DamageOverTime[bucketStr]
			temp.SD += (dmgTotal - temp.Mean) * (dmgTotal - temp.Mean)
			result.DamageOverTime[bucketStr] = temp
		}
	}
	for bucket, resultData := range result.DamageOverTime {
		resultData.SD = math.Sqrt(resultData.SD / float64(n))
		result.DamageOverTime[bucket] = resultData
	}

	// Get standard deviations for statistics
	targetDamage := make(map[int]float64, targetCount)
	for _, v := range data {
		result.DPS.SD += (v.DPS - result.DPS.Mean) * (v.DPS - result.DPS.Mean)

		dd := float64(v.Duration) / 60 //sim reports in frames
		// Reset array
		for j := 0; j < targetCount; j++ {
			for _, charTargetDmg := range v.DamageByCharByTargets {
				targetDamage[j] += charTargetDmg[j] / float64(dd)
			}
			dpsTarget := result.DPSByTarget[j]
			dpsTarget.SD += (targetDamage[j] - dpsTarget.Mean) * (targetDamage[j] - dpsTarget.Mean)
			result.DPSByTarget[j] = dpsTarget

			// Reset
			targetDamage[j] = 0
		}
		if mode {
			result.Duration.SD += (float64(v.Duration) - result.Duration.Mean) * (float64(v.Duration) - result.Duration.Mean)
		}
	}

	result.DPS.SD = math.Sqrt(result.DPS.SD / float64(n))
	for j := 0; j < targetCount; j++ {
		dpsTargetRollup := result.DPSByTarget[j]
		dpsTargetRollup.SD = math.Sqrt(dpsTargetRollup.SD / float64(n))
		result.DPSByTarget[j] = dpsTargetRollup
	}

	// required ER

	if erCalcMode {

		/*
				initialize a two dimensional array, the first index representing the character
				every characters array is supposed to be a list of the minimum amount of "current energy during burst"
			    (read: the maximum amount of needed ER) of each iteration
				afterwards it is possible to use most statistical summary methods on those arrays, in this case we are
				using the mode to determine how much ER is needed in most cases
		*/

		accEnergy := make([][]float64, charCount)

		for i := 0; i < charCount; i++ {
			accEnergy[i] = make([]float64, 0, len(data))
		}

		for i := 0; i < len(data); i++ {

			for j := 0; j < charCount; j++ {
				current, _ := stats.Min(data[i].EnergyWhenBurst[j])

				// for simplcity we are already converting the current energies to the amount of ER needed in that case
				current = data[i].EnergyWhenBurst[j][0] / current

				accEnergy[j] = append(accEnergy[j], current)

			}

		}

		result.RequiredER = make([]float64, charCount)

		for i := 0; i < charCount; i++ {
			modes, _ := stats.Mode(accEnergy[i])
			result.RequiredER[i] = modes[0]
		}

	}

	return
}

func worker(src string, opt core.RunOpt, resp chan workerResp, req chan int64, done chan bool, cust ...func(*Simulation) error) {

	opt.Debug = false
	opt.DebugPaths = []string{}

	for {
		select {
		case seed := <-req:
			parser := parse.New("single", src)
			cfg, _, _ := parser.Parse()

			s, err := NewSim(cfg, seed, opt, cust...)
			if err != nil {
				resp <- workerResp{
					err: err,
				}
				return
			}

			stat, err := s.Run()
			stat.seed = seed

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

func (r *Result) PrettyPrint() string {

	var sb strings.Builder

	for i, t := range r.DamageByChar {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
		}
		sb.WriteString(fmt.Sprintf("%v contributed the following dps:\n", r.CharNames[i]))
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var total float64
		for _, k := range keys {
			v := t[k]
			damageInstances := r.DamageInstancesByChar[i][k]
			sb.WriteString(fmt.Sprintf("\t%v (%.2f%% of total, %.2f average damage procs): avg %.2f [min: %.2f | max: %.2f] \n", k, 100*v.Mean/r.DPS.Mean, damageInstances.Mean, v.Mean, v.Min, v.Max))
			total += v.Mean
		}

		sb.WriteString(fmt.Sprintf("%v total avg dps: %.2f; total percentage: %.0f%%\n", r.CharNames[i], total, 100*total/r.DPS.Mean))
	}
	for i, t := range r.AbilUsageCountByChar {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Character ability usage:\n")
		}
		sb.WriteString(fmt.Sprintf("%v used the following abilities:\n", r.CharNames[i]))
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
	for i, v := range r.CharActiveTime {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")

			sb.WriteString("Character field time:\n")
		}
		sb.WriteString(fmt.Sprintf("%v on average active for %.0f%% [min: %.0f%% | max: %.0f%%]\n", r.CharNames[i], 100*v.Mean/(r.Duration.Mean*60), float64(100*v.Min)/(r.Duration.Mean*60), float64(100*v.Max)/(r.Duration.Mean*60)))
	}
	pk := make([]string, 0, len(r.ParticleCount))
	for k := range r.ParticleCount {
		pk = append(pk, k)
	}
	sort.Strings(pk)
	for i, k := range pk {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Particle count:\n")
		}
		v := r.ParticleCount[k]
		sb.WriteString(fmt.Sprintf("\t%v: avg %.2f [min: %v max: %v]\n", k, v.Mean, v.Min, v.Max))
	}
	rk := make([]core.ReactionType, 0, len(r.ReactionsTriggered))
	for k := range r.ReactionsTriggered {
		rk = append(rk, k)
	}
	for i, k := range rk {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Reactions:\n")
		}
		v := r.ReactionsTriggered[k]
		sb.WriteString(fmt.Sprintf("\t%v: avg %.2f [min: %v max: %v]\n", k, v.Mean, v.Min, v.Max))
	}
	for i, m := range r.ElementUptime {
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Element up time:\n")
		}
		sb.WriteString(fmt.Sprintf("\tTarget #%v\n", i+1))
		for j, ele := range core.EleTypeString {
			v, ok := m[core.EleType(j)]
			if ok {
				if ele == "" {
					ele = "none"
				}
				sb.WriteString(fmt.Sprintf("\t\t%v: avg %.2f%% [min: %.2f%% max: %.2f%%]\n", ele, 100*v.Mean/(r.Duration.Mean*60), float64(100*v.Min)/(r.Duration.Mean*60), float64(100*v.Max)/(r.Duration.Mean*60)))
			}
		}
	}

	// Recommended ER, only in ER calc mode
	if r.RequiredER != nil {
		sb.WriteString("------------------------------------------\n")
		sb.WriteString("Recommended Total Energy Recharge:\n")

		for i, t := range r.RequiredER {
			sb.WriteString(fmt.Sprintf("\t%v: %.0f%% \n", r.CharNames[i], t*100))
		}

	}

	flagDamageByTargets := true
	for i, t := range r.DamageByCharByTargets {
		// Save some space if there is only one target - redundant information
		if len(t) == 1 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Damage by Target Omitted (Only 1 Target)\n")
			flagDamageByTargets = false
			break
		}
		if i == 0 {
			sb.WriteString("------------------------------------------\n")
			sb.WriteString("Damage by Target:\n")
		}
		sb.WriteString(fmt.Sprintf("%v contributed the following dps:\n", r.CharNames[i]))
		keys := make([]int, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		var total float64
		for _, k := range keys {
			v := t[k]
			sb.WriteString(fmt.Sprintf("\t%v (%.2f%% of total): avg %.2f [min: %.2f | max: %.2f] \n", k, 100*v.Mean/r.DPS.Mean, v.Mean, v.Min, v.Max))
			total += v.Mean
		}

		sb.WriteString(fmt.Sprintf("%v total avg dps: %.2f; total percentage: %.0f%%\n", r.CharNames[i], total, 100*total/r.DPS.Mean))
	}
	if flagDamageByTargets {
		for i := range r.DPSByTarget {
			sb.WriteString(fmt.Sprintf("%v (%.2f%% of total): Average %.2f DPS over %.2f seconds (std: %.2f)\n", i, 100*r.DPSByTarget[i].Mean/r.DPS.Mean, r.DPSByTarget[i].Mean, r.Duration.Mean, r.DPSByTarget[i].SD))
		}
	}

	sb.WriteString("------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Average %.2f damage over %.2f seconds, resulting in %.0f dps (min: %.2f max: %.2f std: %.2f) \n", r.Damage.Mean, r.Duration.Mean, r.DPS.Mean, r.DPS.Min, r.DPS.Max, r.DPS.SD))
	sb.WriteString(fmt.Sprintf("Simulation completed %v iterations in %v\n", r.Iterations, r.Runtime))

	return sb.String()
}
