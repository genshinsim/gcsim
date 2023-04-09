package result

import (
	"fmt"
	"sort"
	"strings"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

// TODO: only here for support with the existing UI
// TODO: this entire results class will be refactored with the stats redesign
// TODO: initial refactor: https://github.com/unleashurgeek/gcsim/blob/18b9356f51294270affb846382bc1f292ffb1c23/pkg/result/results.go

type Summary struct {
	//version stuff
	V2                    bool                       `json:"v2"`
	Version               string                     `json:"version"`
	BuildDate             string                     `json:"build_date"`
	IsDamageMode          bool                       `json:"is_damage_mode"`
	ActiveChar            string                     `json:"active_char"`
	CharNames             []string                   `json:"char_names"`
	DamageByChar          []map[string]agg.FloatStat `json:"damage_by_char"`
	DamageInstancesByChar []map[string]agg.IntStat   `json:"damage_instances_by_char"`
	DamageByCharByTargets []map[int]agg.FloatStat    `json:"damage_by_char_by_targets"`
	CharActiveTime        []agg.IntStat              `json:"char_active_time"`
	AbilUsageCountByChar  []map[string]agg.IntStat   `json:"abil_usage_count_by_char"`
	ParticleCount         map[string]agg.FloatStat   `json:"particle_count"`
	ReactionsTriggered    map[string]agg.IntStat     `json:"reactions_triggered"`
	ElementUptime         []map[string]agg.IntStat   `json:"ele_uptime"`
	RequiredER            []float64                  `json:"required_er"`
	IncompleteChars       []string                   `json:"incomplete_chars"`

	Duration agg.FloatStat `json:"sim_duration"`
	//final result
	Damage         agg.FloatStat            `json:"damage"`
	DPS            agg.FloatStat            `json:"dps"`
	DPSByTarget    map[int]agg.FloatStat    `json:"dps_by_target"`
	DamageOverTime map[string]agg.FloatStat `json:"damage_over_time"`

	Iterations int     `json:"iter"`
	Runtime    float64 `json:"runtime"`
	//other info
	NumTargets    int                          `json:"num_targets"` //TODO: to deprecate this
	CharDetails   []simulation.CharacterDetail `json:"char_details"`
	TargetDetails []enemy.EnemyProfile         `json:"target_details"`
	//for tracking min/max run
	MinSeed int64 `json:"-"`
	MaxSeed int64 `json:"-"`
	//put these last so result is kinda readable by human
	Config         string                   `json:"config_file"`
	Text           string                   `json:"text"`
	Debug          []map[string]interface{} `json:"debug"`
	DebugMinDPSRun []map[string]interface{} `json:"debug_min_dps_run,omitempty"`
	DebugMaxDPSRun []map[string]interface{} `json:"debug_max_dps_run,omitempty"`
}

// TODO: very temporary mess to have the new aggregate system connect to the legacy summary
func (r *Summary) Map(cfg *ast.ActionList, result *agg.Result) {
	for _, v := range cfg.Characters {
		r.CharNames = append(r.CharNames, v.Base.Key.String())

		if !IsCharacterComplete(v.Base.Key) {
			r.IncompleteChars = append(r.IncompleteChars, v.Base.Key.String())
		}
	}

	// metadata agg
	r.MinSeed = int64(result.MinSeed)
	r.MaxSeed = int64(result.MaxSeed)
	r.Duration = result.Duration

	// overview agg
	r.DPS = result.DPS
	r.Damage = result.TotalDamage

	// legacy agg
	r.DamageByChar = result.Legacy.DamageByChar
	r.DamageInstancesByChar = result.Legacy.DamageInstancesByChar
	r.DamageByCharByTargets = result.Legacy.DamageByCharByTargets
	r.CharActiveTime = result.Legacy.CharActiveTime
	r.AbilUsageCountByChar = result.Legacy.AbilUsageCountByChar
	r.ParticleCount = result.Legacy.ParticleCount
	r.ReactionsTriggered = result.Legacy.ReactionsTriggered
	r.ElementUptime = result.Legacy.ElementUptime
	r.DPSByTarget = result.Legacy.DPSByTarget
	r.DamageOverTime = result.Legacy.DamageOverTime
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
	rk := make([]string, 0, len(r.ReactionsTriggered))
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
		for k, v := range m {
			if k == "" {
				k = "none"
			}
			sb.WriteString(fmt.Sprintf("\t\t%v: avg %.2f%% [min: %.2f%% max: %.2f%%]\n", k, 100*v.Mean/(r.Duration.Mean*60), float64(100*v.Min)/(r.Duration.Mean*60), float64(100*v.Max)/(r.Duration.Mean*60)))
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
