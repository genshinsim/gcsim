package core

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func (c *Core) QueueParticle(src string, num float64, ele attributes.Element, delay int) {
	p := character.Particle{
		Source: src,
		Num:    num,
		Ele:    ele,
	}
	if delay == 0 {
		c.Player.DistributeParticle(p)
		return
	}
	if delay < 0 {
		panic("queue particle called with delay < 0")
	}
	c.Tasks.Add(func() {
		c.Player.DistributeParticle(p)
	}, delay)
}

func (c *Core) SetupOnNormalHitEnergy() {
	var current [MaxTeamSize][info.EndWeaponClass]float64

	// https://genshin-impact.fandom.com/wiki/Energy#Energy_Generated_by_Normal_Attacks
	// Base Probability
	for i := range current {
		current[i][info.WeaponClassSword] = 0.10 // WeaponClassSword
	}
	// Probability Increase Per Fail
	inc := []float64{
		0.05, // WeaponClassSword
		0.10, // WeaponClassClaymore
		0.04, // WeaponClassSpear
		0.05, // WeaponClassBow
		0.10, // WeaponClassCatalyst
	}

	//TODO: not sure if there's like a 0.2s icd on this. for now let's add it in to be safe
	icd := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}
		// check icd
		if icd > c.F {
			return false
		}
		// check chance
		char := c.Player.ByIndex(atk.Info.ActorIndex)

		if c.Rand.Float64() > current[atk.Info.ActorIndex][char.Weapon.Class] {
			// increment chance
			current[atk.Info.ActorIndex][char.Weapon.Class] += inc[char.Weapon.Class]
			return false
		}

		// add energy
		char.AddEnergy("na-ca-on-hit", 1)
		// Add this log in sim if necessary to see as AddEnergy already generates a log
		c.Log.NewEvent("random energy on normal", glog.LogDebugEvent, char.Index).
			Write("char", atk.Info.ActorIndex).
			Write("chance", current[atk.Info.ActorIndex][char.Weapon.Class])
		// set icd
		icd = c.F + 12
		current[atk.Info.ActorIndex][char.Weapon.Class] = 0
		if char.Weapon.Class == info.WeaponClassSword {
			current[atk.Info.ActorIndex][char.Weapon.Class] = 0.10
		}
		return false
	}, "random-energy-restore-on-hit")

	//TODO: assuming we clear the probability on swap
	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		for i := range current {
			for j := range current[i] {
				current[i][j] = 0
			}
			current[i][info.WeaponClassSword] = 0.10
		}
		return false
	}, "random-energy-restore-on-hit-swap")

}
