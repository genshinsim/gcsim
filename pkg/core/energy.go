package core

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
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
	var current [MaxTeamSize][weapon.EndWeaponClass]float64

	inc := []float64{
		0.05, //WeaponClassSword
		0.05, //WeaponClassClaymore
		0.04, //WeaponClassSpear
		0.01, //WeaponClassBow
		0.01, //WeaponClassCatalyst
	}

	//TODO: not sure if there's like a 0.2s icd on this. for now let's add it in to be safe
	icd := 0
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		//check icd
		if icd > c.F {
			return false
		}
		//check chance
		char := c.Player.ByIndex(atk.Info.ActorIndex)

		if c.Rand.Float64() > current[atk.Info.ActorIndex][char.Weapon.Class] {
			//increment chance
			current[atk.Info.ActorIndex][char.Weapon.Class] += inc[char.Weapon.Class]
			return false
		}

		//add energy
		char.AddEnergy("na-ca-on-hit", 1)
		// Add this log in sim if necessary to see as AddEnergy already generates a log
		c.Log.NewEvent("random energy on normal", glog.LogSimEvent, char.Index, "char", atk.Info.ActorIndex, "chance", current[atk.Info.ActorIndex][char.Weapon.Class])
		//set icd
		icd = c.F + 12
		current[atk.Info.ActorIndex][char.Weapon.Class] = 0
		return false
	}, "random-energy-restore-on-hit")

	//TODO: assuming we clear the probability on swap
	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		for i := range current {
			for j := range current[i] {
				current[i][j] = 0
			}
		}
		return false
	}, "random-energy-restore-on-hit-swap")

}
