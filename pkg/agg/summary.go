package agg

//go:generate go tool github.com/tinylib/msgp -io=false

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/stats"
)

// Summary is one iteration's stats.Result reduced to the values the aggregators read. The
// per-event loops are already done, so a Summary is a small fraction of the size of its Result.
// The web UI builds it in each sim worker and sends it to the aggregator worker instead of the
// full Result.
//
// Every value is computed with the same float operations, in the same order, as the aggregators
// used to compute it from the Result, so aggregated statistics stay bit-identical.
//
//msgp:tuple Summary CharacterSummary EnemySummary FailureTimes TargetDamage
type Summary struct {
	Seed          uint64
	Duration      int
	TotalDamage   float64
	DPS           float64
	DamageBuckets []float64
	ShieldResults stats.ShieldResult
	TargetOverlap bool
	EndStats      []stats.EndStats

	// Totals over all characters, summed in character order and then in event order.
	Reactions int     // reactions triggered
	Heal      float64 // healing done
	Energy    float64 // energy received, including energy wasted at full energy
	Failures  FailureTimes

	Characters []CharacterSummary
	Enemies    []EnemySummary
}

type CharacterSummary struct {
	ActiveTime int
	Failures   FailureTimes

	Actions   map[string]float64 // times each action was used
	Reactions map[string]float64 // reactions triggered, by reaction
	Energy    map[string]float64 // energy received (including wasted), by source

	Damage          float64            // total damage
	DamageByElement map[string]float64 // every element, 0 if it did no damage
	DamageByTarget  []TargetDamage
	// damage by source, with the reaction modifier appended to the key, e.g. "Skill (vaporize)"
	DamageBySource          map[string]float64
	DamageInstances         map[string]float64 // hits that did damage, same keys as DamageBySource
	DamageCumulativeContrib []float64
}

type EnemySummary struct {
	ReactionUptime   map[string]int
	CumulativeDamage []float64
}

// FailureTimes is how long actions failed to execute, in seconds, by reason.
type FailureTimes struct {
	Energy  float64
	Stamina float64
	Swap    float64
	Skill   float64
	Dash    float64
	BurstCD float64
}

type TargetDamage struct {
	Target int
	Damage float64
}

// Summarize reduces result to a Summary. The Summary shares slices and maps with result.
func Summarize(result *stats.Result) Summary {
	s := Summary{
		Seed:          result.Seed,
		Duration:      result.Duration,
		TotalDamage:   result.TotalDamage,
		DPS:           result.DPS,
		DamageBuckets: result.DamageBuckets,
		ShieldResults: result.ShieldResults,
		TargetOverlap: result.TargetOverlap,
		EndStats:      result.EndStats,
		Characters:    make([]CharacterSummary, len(result.Characters)),
		Enemies:       make([]EnemySummary, len(result.Enemies)),
	}

	for i := range result.Characters {
		char := &result.Characters[i]
		c := &s.Characters[i]

		c.ActiveTime = char.ActiveTime
		for _, fail := range char.FailedActions {
			c.Failures.add(fail)
			s.Failures.add(fail)
		}

		c.Actions = make(map[string]float64)
		for _, ev := range char.ActionEvents {
			c.Actions[ev.Action] += 1
		}

		s.Reactions += len(char.ReactionEvents)
		c.Reactions = make(map[string]float64)
		for _, ev := range char.ReactionEvents {
			c.Reactions[ev.Reaction] += 1
		}

		for _, h := range char.HealEvents {
			s.Heal += h.Heal
		}

		c.Energy = make(map[string]float64)
		for _, ev := range char.EnergyEvents {
			s.Energy += ev.Gained + ev.Wasted
			c.Energy[ev.Source] += ev.Gained + ev.Wasted
		}

		c.DamageByElement = make(map[string]float64)
		for _, ele := range attributes.ElementStrings() {
			c.DamageByElement[ele] = 0
		}
		c.DamageBySource = make(map[string]float64)
		c.DamageInstances = make(map[string]float64)
		for _, ev := range char.DamageEvents {
			c.DamageByTarget = addTargetDamage(c.DamageByTarget, ev.Target, ev.Damage)
			c.DamageByElement[ev.Element] += ev.Damage
			c.Damage += ev.Damage
			key := ev.Source
			if ev.ReactionModifier != "" {
				key += " (" + string(ev.ReactionModifier) + ")"
			}
			c.DamageBySource[key] += ev.Damage
			if ev.Damage > 0 {
				c.DamageInstances[key] += 1
			}
		}
		c.DamageCumulativeContrib = char.DamageCumulativeContrib
	}

	for i := range result.Enemies {
		s.Enemies[i] = EnemySummary{
			ReactionUptime:   result.Enemies[i].ReactionUptime,
			CumulativeDamage: result.Enemies[i].CumulativeDamage,
		}
	}
	return s
}

func (f *FailureTimes) add(fail stats.ActionFailInterval) {
	switch fail.Reason {
	case action.InsufficientEnergy.String():
		f.Energy += float64(fail.End-fail.Start) / 60
	case action.InsufficientStamina.String():
		f.Stamina += float64(fail.End-fail.Start) / 60
	case action.SwapCD.String():
		f.Swap += float64(fail.End-fail.Start) / 60
	case action.SkillCD.String():
		f.Skill += float64(fail.End-fail.Start) / 60
	case action.DashCD.String():
		f.Dash += float64(fail.End-fail.Start) / 60
	case action.BurstCD.String():
		f.BurstCD += float64(fail.End-fail.Start) / 60
	}
}

// addTargetDamage adds damage to target's entry like a map[int]float64 would, including starting
// from 0 (0 + -0 is 0, so a new entry doesn't simply store damage).
func addTargetDamage(targets []TargetDamage, target int, damage float64) []TargetDamage {
	i := 0
	for i < len(targets) && targets[i].Target != target {
		i++
	}
	if i == len(targets) {
		targets = append(targets, TargetDamage{Target: target})
	}
	targets[i].Damage += damage
	return targets
}
