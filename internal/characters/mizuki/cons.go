package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key               = "mizuki-c1"
	c1Interval          = 3.5 * 60
	c1Duration          = 3 * 60
	c1Multiplier        = 11.0
	c1Range             = 30
	c2Key               = "mizuki-c2"
	c2EMMultiplier      = 0.0004
	c4EnergyGenerations = 4
	c4Key               = "mizuki-c4"
	c4Energy            = 5
	c6Key               = "mizuki-c6"
	c6CR                = 0.3
	c6CD                = 1.0
)

// When Yumemizuki Mizuki is in the Dreamdrifter state, she will continuously apply the "Twenty-Three Nights' Awaiting"
// effect to nearby opponents for 3s every 3.5s. When an opponent is affected by Anemo DMG-triggered Swirl reactions
// while the aforementioned effect is active, the effect will be canceled and this Swirl instance has its DMG against
// this opponent increased by 1100% of Mizuki's Elemental Mastery.
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		// Check if enemy has the debuff
		e := args[0].(*enemy.Enemy)
		if !e.StatusIsActive(c1Key) {
			return false
		}

		// Only on swirls. The swirl source does not matter, it can be either mizuki or another anemo char.
		switch atk.Info.AttackTag {
		case attacks.AttackTagSwirlCryo:
		case attacks.AttackTagSwirlElectro:
		case attacks.AttackTagSwirlHydro:
		case attacks.AttackTagSwirlPyro:
		default:
			return false
		}

		// do not proc on 0 DMG swirls (e.g. hydro AOE swirls or swirl ICD)
		if atk.Info.FlatDmg == 0 {
			return false
		}

		additionalDmg := c1Multiplier * c.Stat(attributes.EM)

		c.Core.Log.NewEvent("mizuki c1 proc", glog.LogPreDamageMod, atk.Info.ActorIndex).
			Write("before", atk.Info.FlatDmg).
			Write("addition", additionalDmg).
			Write("final", atk.Info.FlatDmg+additionalDmg)

		atk.Info.FlatDmg += additionalDmg
		atk.Info.Abil += " (Mizuki C1)"

		// Cancel the effect
		e.DeleteStatus(c1Key)

		return false
	}, c1Key)
}

func (c *char) applyC1Effect() {
	var c1Func func()
	c1Func = func() {
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return
		}

		// The range cannot be determined exactly, but from testing it seems that it is at least 25 meters or more.
		// For simulation purposes we can assume a bit more than that.
		area := combat.AttackPattern{
			Shape: geometry.NewCircle(c.Core.Combat.Player().Pos(), c1Range, geometry.DefaultDirection(), 360),
		}
		for _, target := range c.Core.Combat.EnemiesWithinArea(area, nil) {
			if e, ok := target.(*enemy.Enemy); ok {
				// is it even possible to verify if it is affected by hitlag?
				e.AddStatus(c1Key, c1Duration, true)
			}
		}
		c.QueueCharTask(c1Func, c1Interval)
	}

	// Debuff does not take 3.5s to apply but does not trigger on initial skill activation swirl according to testing.
	// First cloud (0.45s after skill activation) can trigger it so queue it a few frames later
	c.QueueCharTask(c1Func, skillHitmark+2)
}

// When Yumemizuki Mizuki enters the Dreamdrifter state, every Elemental Mastery point she has will increase all nearby
// party members' Pyro, Hydro, Cryo, and Electro DMG Bonuses by 0.04% until the Dreamdrifter state ends.
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2Buff = make([]float64, attributes.EndStatType)

	for _, char := range c.Core.Player.Chars() {
		if char.Index == c.Index {
			continue
		}
		// TODO: Test whether this is indeed a static buff once we have C2
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBase(c2Key, -1),
			Amount: func() ([]float64, bool) {
				if !c.StatusIsActive(dreamDrifterStateKey) {
					return nil, false
				}
				dmgBonus := c.Stat(attributes.EM) * c2EMMultiplier
				c.c2Buff[attributes.PyroP] = dmgBonus
				c.c2Buff[attributes.HydroP] = dmgBonus
				c.c2Buff[attributes.ElectroP] = dmgBonus
				c.c2Buff[attributes.CryoP] = dmgBonus
				return c.c2Buff, true
			},
		})
	}
}

// While Yumemizuki Mizuki is in the Dreamdrifter state, Swirl DMG dealt by nearby party members can Crit,
// with CRIT Rate fixed at 30%, and CRIT DMG fixed at 100%.
func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		ae := args[1].(*combat.AttackEvent)

		// Only on swirls. The swirl source does not matter, it can be either mizuki or other anemo char.
		switch ae.Info.AttackTag {
		case attacks.AttackTagSwirlPyro:
		case attacks.AttackTagSwirlCryo:
		case attacks.AttackTagSwirlHydro:
		case attacks.AttackTagSwirlElectro:
		default:
			return false
		}

		// The effect is only when mizuki is in dreamDrifter state
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return false
		}

		// Crit rate/DMG is fixed to 30% CR and 100% CD
		ae.Snapshot.Stats[attributes.CR] = c6CR
		ae.Snapshot.Stats[attributes.CD] = c6CD

		return false
	}, c6Key)
}
