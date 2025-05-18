package klee

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	a1IcdKey   = "a1-icd"
	a1SparkKey = "a1-spark"
)

// When Jumpy Dumpty and Normal Attacks deal DMG, Klee has a 50% chance to obtain an Explosive Spark.
// This Explosive Spark is consumed by the next Charged Attack, which costs no Stamina and deals 50% increased DMG.
func (c *char) makeA1CB() combat.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	return func(a combat.AttackCB) {
		if c.StatusIsActive(a1IcdKey) {
			return
		}
		if c.Core.Rand.Float64() < 0.5 {
			return
		}
		c.AddStatus(a1IcdKey, 60*5, true)

		if !c.StatusIsActive(a1SparkKey) {
			c.AddStatus(a1SparkKey, 60*30, true)
		}
	}
}

const a4ICDKey = "klee-a4-icd"

// When Klee's Charged Attack results in a CRIT Hit, all party members gain 2 Elemental Energy.
func (c *char) makeA4CB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if !a.IsCrit {
			return
		}
		if c.StatusIsActive(a4ICDKey) {
			return
		}
		c.AddStatus(a4ICDKey, 0.6*60, true)
		for _, x := range c.Core.Player.Chars() {
			x.AddEnergy("klee-a4", 2)
		}
	}
}
