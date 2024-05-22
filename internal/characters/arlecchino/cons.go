package arlecchino

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2IcdKey = "arlecchino-c2-icd"
const c4IcdKey = "arlecchino-c4-icd"
const c6IcdKey = "arlecchino-c6-icd"
const c6Key = "arlecchino-c6"

func (c *char) c2() {
	c.initialDirectiveLevel = 1
	if c.Base.Cons < 2 || c.Base.Ascension < 1 {
		return
	}

	c.initialDirectiveLevel = 2
}

func (c *char) c2OnAbsorbDue() {
	// Check is redundant? Can't reach blood debt due without A1
	if c.Base.Cons < 2 || c.Base.Ascension < 1 {
		return
	}

	if c.StatusIsActive(c2IcdKey) {
		return
	}

	c.AddStatus(c2IcdKey, 10*60, true)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Balemoon Bloodfire (C2)",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       9.00,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: 3},
			6.5,
		),
		50,
		50,
	)
}

func (c *char) c4OnAbsorb() {
	if c.Base.Cons < 4 {
		return
	}

	if c.StatusIsActive(c4IcdKey) {
		return
	}

	c.AddStatus(c4IcdKey, 10*60, true)
	c.ReduceActionCooldown(action.ActionBurst, 2*60)
	c.AddEnergy("arlecchino-c4", 15)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}

		if ae.Info.AttackTag != attacks.AttackTagElementalBurst {
			return false
		}

		amt := c.getTotalAtk() * 7.0 * c.CurrentHPDebt() / c.MaxHP()
		c.Core.Log.NewEvent("Arlecchino C6 dmg add", glog.LogCharacterEvent, c.Index).
			Write("amt", amt)

		ae.Info.FlatDmg += amt

		return false
	}, "arlecchino-c6-burst")
}
func (c *char) c6skill() {
	if c.Base.Cons < 6 {
		return
	}

	if c.StatusIsActive(c6IcdKey) {
		return
	}
	c.AddStatus(c6IcdKey, 15*60, true)

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.1
	m[attributes.CD] = 0.7
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c6Key, 20*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalBurst, attacks.AttackTagNormal:
				return m, true
			}
			return nil, false
		},
	})
}
