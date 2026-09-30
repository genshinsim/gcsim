package agg

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func TestSummarize(t *testing.T) {
	result := stats.Result{
		Seed:     7,
		Duration: 600,
		Characters: []stats.CharacterResult{
			{
				ActiveTime: 300,
				DamageEvents: []stats.DamageEvent{
					{Source: "Skill", Target: 1, Element: "pyro", ReactionModifier: stats.Vaporize, Damage: 100},
					{Source: "Skill", Target: 2, Element: "pyro", Damage: 50},
					{Source: "Skill", Target: 1, Element: "pyro", ReactionModifier: stats.Vaporize, Damage: 0},
					{Source: "Normal 0", Target: 1, Element: "physical", Damage: 25},
					// a source whose name looks like source plus modifier shares its key, as with a map
					{Source: "Burst (melt)", Target: 1, Element: "pyro", Damage: 10},
					{Source: "Burst", Target: 1, Element: "pyro", ReactionModifier: stats.Melt, Damage: 5},
					{Source: "Burst", Target: 1, Element: "pyro", Damage: 1},
				},
				ActionEvents:   []stats.ActionEvent{{Action: "skill"}, {Action: "attack"}, {Action: "skill"}},
				ReactionEvents: []stats.ReactionEvent{{Reaction: "vaporize"}},
				EnergyEvents: []stats.EnergyEvent{
					{Source: "a", Gained: 2, Wasted: 1},
					{Source: "a", Gained: 3},
					{Source: "b", Gained: 1},
				},
				HealEvents: []stats.HealEvent{{Heal: 10}},
				FailedActions: []stats.ActionFailInterval{
					{Start: 0, End: 30, Reason: action.SkillCD.String()},
					{Start: 60, End: 120, Reason: action.InsufficientEnergy.String()},
				},
			},
			{
				ReactionEvents: []stats.ReactionEvent{{Reaction: "melt"}, {Reaction: "melt"}},
				HealEvents:     []stats.HealEvent{{Heal: 5}},
				FailedActions:  []stats.ActionFailInterval{{Start: 0, End: 90, Reason: action.SkillCD.String()}},
			},
		},
	}

	s := Summarize(&result)

	if s.Seed != 7 || s.Duration != 600 {
		t.Errorf("seed/duration = %v/%v, want 7/600", s.Seed, s.Duration)
	}
	if s.Reactions != 3 || s.Heal != 15 || s.Energy != 7 {
		t.Errorf("totals: reactions %v heal %v energy %v, want 3 15 7", s.Reactions, s.Heal, s.Energy)
	}
	if s.Failures.Skill != 2 || s.Failures.Energy != 1 {
		t.Errorf("total failures %+v, want skill 2 and energy 1", s.Failures)
	}

	c := s.Characters[0]
	if c.ActiveTime != 300 || c.Damage != 191 {
		t.Errorf("active time/damage = %v/%v, want 300/191", c.ActiveTime, c.Damage)
	}
	if c.Failures.Skill != 0.5 || c.Failures.Energy != 1 {
		t.Errorf("failures %+v, want skill 0.5 and energy 1", c.Failures)
	}
	wantSums(t, "actions", c.Actions, map[string]float64{"skill": 2, "attack": 1})
	wantSums(t, "reactions", c.Reactions, map[string]float64{"vaporize": 1})
	wantSums(t, "energy", c.Energy, map[string]float64{"a": 6, "b": 1})
	wantSums(t, "damage by source", c.DamageBySource, map[string]float64{
		"Skill (vaporize)": 100, "Skill": 50, "Normal 0": 25, "Burst (melt)": 15, "Burst": 1,
	})
	wantSums(t, "damage instances", c.DamageInstances, map[string]float64{
		"Skill (vaporize)": 1, "Skill": 1, "Normal 0": 1, "Burst (melt)": 2, "Burst": 1,
	})

	byElement := map[string]float64{}
	for _, ele := range attributes.ElementStrings() {
		byElement[ele] = 0
	}
	byElement["pyro"] = 166
	byElement["physical"] = 25
	wantSums(t, "damage by element", c.DamageByElement, byElement)

	targets := map[int]float64{}
	for _, td := range c.DamageByTarget {
		if _, ok := targets[td.Target]; ok {
			t.Errorf("target %v listed twice", td.Target)
		}
		targets[td.Target] = td.Damage
	}
	if len(targets) != 2 || targets[1] != 141 || targets[2] != 50 {
		t.Errorf("damage by target %v, want 1: 141, 2: 50", c.DamageByTarget)
	}

	wantSums(t, "second character reactions", s.Characters[1].Reactions, map[string]float64{"melt": 2})
}

func wantSums(t *testing.T, name string, got Sums, want map[string]float64) {
	t.Helper()
	m := map[string]float64{}
	for k, v := range got.All() {
		if _, ok := m[k]; ok {
			t.Errorf("%v: key %q listed twice", name, k)
		}
		m[k] = v
	}
	if len(m) != len(want) || len(got.Values) != len(got.Keys) {
		t.Errorf("%v: got %v, want %v", name, got, want)
		return
	}
	for k, v := range want {
		if g, ok := m[k]; !ok || g != v {
			t.Errorf("%v: got %v, want %v", name, m, want)
			return
		}
	}
}
