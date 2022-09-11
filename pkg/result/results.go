package result

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

type Summary struct {
	//version stuff
	V2                    bool                               `json:"v2"`
	Version               string                             `json:"version"`
	BuildDate             string                             `json:"build_date"`
	IsDamageMode          bool                               `json:"is_damage_mode"`
	ActiveChar            string                             `json:"active_char"`
	CharNames             []string                           `json:"char_names"`
	DamageByChar          []map[string]FloatResult           `json:"damage_by_char"`
	DamageInstancesByChar []map[string]IntResult             `json:"damage_instances_by_char"`
	DamageByCharByTargets []map[int]FloatResult              `json:"damage_by_char_by_targets"`
	CharActiveTime        []IntResult                        `json:"char_active_time"`
	AbilUsageCountByChar  []map[string]IntResult             `json:"abil_usage_count_by_char"`
	ParticleCount         map[string]FloatResult             `json:"particle_count"`
	ReactionsTriggered    map[combat.ReactionType]IntResult  `json:"reactions_triggered"`
	Duration              FloatResult                        `json:"sim_duration"`
	ElementUptime         []map[attributes.Element]IntResult `json:"ele_uptime"`
	RequiredER            []float64                          `json:"required_er"`
	//final result
	Damage         FloatResult            `json:"damage"`
	DPS            FloatResult            `json:"dps"`
	DPSByTarget    map[int]FloatResult    `json:"dps_by_target"`
	DamageOverTime map[string]FloatResult `json:"damage_over_time"`
	Iterations     int                    `json:"iter"`
	Runtime        float64                `json:"runtime"`
	//other info
	NumTargets    int                     `json:"num_targets"` //TODO: to deprecate this
	CharDetails   []simulation.CharDetail `json:"char_details"`
	TargetDetails []enemy.EnemyProfile    `json:"target_details"`
	//for tracking min/max run
	MinSeed int64 `json:"-"`
	MaxSeed int64 `json:"-"`
	//put these last so result is kinda readable by human
	Config         string `json:"config_file"`
	Text           string `json:"text"`
	Debug          string `json:"debug"`
	DebugMinDPSRun string `json:"debug_min_dps_run"`
	DebugMaxDPSRun string `json:"debug_max_dps_run"`
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

func CollectResult(data []simulation.Result, mode bool, chars []string, detailed bool) (result Summary) {

	n := len(data)

	// TODO: Kind of brittle - maybe track something separate for this?
	charCount := len(chars)
	result.DPS.Min = math.MaxFloat64
	result.DPS.Max = -1
	result.DPSByTarget = make(map[int]FloatResult)
	result.DamageOverTime = make(map[string]FloatResult)
	if detailed {
		result.ReactionsTriggered = make(map[combat.ReactionType]IntResult)
		result.CharNames = make([]string, charCount)
		result.ParticleCount = make(map[string]FloatResult)
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
			result.Damage.Min = v.Damage

		}
		if v.Damage > result.Damage.Max {
			result.Damage.Max = v.Damage
		}
		result.Damage.Mean += v.Damage / float64(n)

		//dps
		if result.DPS.Min > v.DPS {
			result.DPS.Min = v.DPS
			result.MinSeed = v.Seed
		}
		if result.DPS.Max < v.DPS {
			result.DPS.Max = v.DPS
			result.MaxSeed = v.Seed
		}
		result.DPS.Mean += v.DPS / float64(n)

		if !detailed {
			continue
		}

		damageOverTimeByRun[iteration] = make(map[float64]float64)
		// Damage Over Time - get data for all iterations first and summarized later
		for frameBucket, damage := range v.DamageDetailByTime {
			// Convert frame bucket value into seconds
			secBucket := float64(frameBucket) / 60.0

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
				result.ElementUptime = append(result.ElementUptime, make(map[attributes.Element]IntResult))
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
	for _, char := range result.DamageByCharByTargets {
		for idxTarget, dpsFigure := range char {
			dpsResult := result.DPSByTarget[idxTarget]
			dpsResult.Mean += dpsFigure.Mean
			result.DPSByTarget[idxTarget] = dpsResult
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
	targetDamage := make(map[int]float64)
	for _, v := range data {
		result.DPS.SD += (v.DPS - result.DPS.Mean) * (v.DPS - result.DPS.Mean)

		dd := float64(v.Duration) / 60 //sim reports in frames
		for _, charDmg := range v.DamageByCharByTargets {
			for idxTarget, dmg := range charDmg {
				targetDamage[idxTarget] += dmg / float64(dd)
			}
		}
		for idxTarget, dmg := range targetDamage {
			dpsTarget := result.DPSByTarget[idxTarget]
			dpsTarget.SD += (dmg - dpsTarget.Mean) * (dmg - dpsTarget.Mean)
			result.DPSByTarget[idxTarget] = dpsTarget

			// Reset
			targetDamage[idxTarget] = 0
		}
		if mode {
			result.Duration.SD += (float64(v.Duration) - result.Duration.Mean) * (float64(v.Duration) - result.Duration.Mean)
		}
	}

	result.DPS.SD = math.Sqrt(result.DPS.SD / float64(n))

	for idxTarget := range result.DPSByTarget {
		dpsTargetRollup := result.DPSByTarget[idxTarget]
		dpsTargetRollup.SD = math.Sqrt(dpsTargetRollup.SD / float64(n))
		result.DPSByTarget[idxTarget] = dpsTargetRollup
	}

	return
}

func (r *Summary) PrettyPrint() string {

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
	rk := make([]combat.ReactionType, 0, len(r.ReactionsTriggered))
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
		for j, ele := range attributes.ElementString {
			v, ok := m[attributes.Element(j)]
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
		keys := make([]int, 0, len(r.DPSByTarget))
		for k := range r.DPSByTarget {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		for _, i := range keys {
			sb.WriteString(fmt.Sprintf("%v (%.2f%% of total): Average %.2f DPS over %.2f seconds (std: %.2f)\n", i, 100*r.DPSByTarget[i].Mean/r.DPS.Mean, r.DPSByTarget[i].Mean, r.Duration.Mean, r.DPSByTarget[i].SD))
		}
	}

	sb.WriteString("------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Average duration of %.2f seconds (min: %.2f max: %.2f std: %.2f)\n", r.Duration.Mean, r.Duration.Min, r.Duration.Max, r.Duration.SD))
	sb.WriteString(fmt.Sprintf("Average %.2f damage over %.2f seconds, resulting in %.0f dps (min: %.2f max: %.2f std: %.2f) \n", r.Damage.Mean, r.Duration.Mean, r.DPS.Mean, r.DPS.Min, r.DPS.Max, r.DPS.SD))
	sb.WriteString(fmt.Sprintf("Simulation completed %v iterations in %.3f seconds\n", r.Iterations, r.Runtime/1000000000))

	return sb.String()
}
