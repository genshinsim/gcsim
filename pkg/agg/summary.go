package agg

//go:generate go tool github.com/tinylib/msgp -io=false

import (
	"iter"

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
//msgp:tuple Summary CharacterSummary EnemySummary FailureTimes TargetDamage Sums
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

	Actions   Sums // times each action was used
	Reactions Sums // reactions triggered, by reaction
	Energy    Sums // energy received (including wasted), by source

	Damage          float64 // total damage
	DamageByElement Sums    // every element, 0 if it did no damage
	DamageByTarget  []TargetDamage
	// damage by source, with the reaction modifier appended to the key, e.g. "Skill (vaporize)"
	DamageBySource          Sums
	DamageInstances         Sums // hits that did damage, same keys as DamageBySource
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

// Sums holds a float64 per string key, like a map[string]float64, as keys and values in the
// order the keys were first seen. An iteration has few keys per sum, so a linear search is
// cheaper than hashing (wasm has no fast string hash), and Sums encode and decode without
// building maps.
type Sums struct {
	Keys   []string
	Values []float64
}

// All iterates over the keys and their values.
func (s *Sums) All() iter.Seq2[string, float64] {
	return func(yield func(string, float64) bool) {
		for i, k := range s.Keys {
			if !yield(k, s.Values[i]) {
				return
			}
		}
	}
}

// add adds v to key's value, which starts from 0 like a map's zero value.
func (s *Sums) add(key string, v float64) {
	i := 0
	for i < len(s.Keys) && s.Keys[i] != key {
		i++
	}
	if i == len(s.Keys) {
		s.Keys = append(s.Keys, key)
		s.Values = append(s.Values, 0)
	}
	s.Values[i] += v
}

// sourceIndex returns the index of the damage aggregator's key for a source and reaction
// modifier: the source, followed by the modifier in parentheses if there is one. The key is
// built only when it's new, with value 0.
func (s *Sums) sourceIndex(source string, modifier stats.ReactionModifier) int {
	m := string(modifier)
	for i, k := range s.Keys {
		if isSourceKey(k, source, m) {
			return i
		}
	}
	key := source
	if m != "" {
		key += " (" + m + ")"
	}
	s.Keys = append(s.Keys, key)
	s.Values = append(s.Values, 0)
	return len(s.Keys) - 1
}

// isSourceKey reports whether key is source + " (" + modifier + ")", or source if there is no
// modifier, without building that string.
func isSourceKey(key, source, modifier string) bool {
	if modifier == "" {
		return key == source
	}
	n := len(source)
	return len(key) == n+len(modifier)+3 && key[:n] == source && key[n:n+2] == " (" &&
		key[n+2:len(key)-1] == modifier && key[len(key)-1] == ')'
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

		for _, ev := range char.ActionEvents {
			c.Actions.add(ev.Action, 1)
		}

		s.Reactions += len(char.ReactionEvents)
		for _, ev := range char.ReactionEvents {
			c.Reactions.add(ev.Reaction, 1)
		}

		for _, h := range char.HealEvents {
			s.Heal += h.Heal
		}

		for _, ev := range char.EnergyEvents {
			s.Energy += ev.Gained + ev.Wasted
			c.Energy.add(ev.Source, ev.Gained+ev.Wasted)
		}

		elements := attributes.ElementStrings()
		c.DamageByElement = Sums{Keys: elements, Values: make([]float64, len(elements))}
		var instances []float64 // by DamageBySource index
		for _, ev := range char.DamageEvents {
			c.DamageByTarget = addTargetDamage(c.DamageByTarget, ev.Target, ev.Damage)
			c.DamageByElement.add(ev.Element, ev.Damage)
			c.Damage += ev.Damage
			j := c.DamageBySource.sourceIndex(ev.Source, ev.ReactionModifier)
			c.DamageBySource.Values[j] += ev.Damage
			if j == len(instances) {
				instances = append(instances, 0)
			}
			if ev.Damage > 0 {
				instances[j] += 1
			}
		}
		for j, n := range instances {
			if n > 0 {
				c.DamageInstances.Keys = append(c.DamageInstances.Keys, c.DamageBySource.Keys[j])
				c.DamageInstances.Values = append(c.DamageInstances.Values, n)
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
