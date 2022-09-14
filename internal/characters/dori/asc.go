package dori

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// A1:
// After a character connected to the Jinni triggers an Electro-Charged, Superconduct, Overloaded, Quicken, Aggravate, Hyperbloom,
//
//	or an Electro Swirl or Crystallize reaction, the CD of Spirit-Warding Lamp: Troubleshooter Cannon is decreased by 1s.
//
// This effect can be triggered once every 3s.
func (c *char) a1() {
	const icdKey = "dori-a1"
	icd := 180 // 3s * 60

	reduce := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		if c.Core.Player.Active() != atk.Info.ActorIndex { //only for on field character
			return false
		}
		if c.StatusIsActive(icdKey) {
			return false
		}
		c.AddStatus(icdKey, icd, true)
		c.ReduceActionCooldown(action.ActionSkill, 60)
		c.Core.Log.NewEvent("dori a1 proc", glog.LogCharacterEvent, c.Index).
			Write("reaction", atk.Info.Abil).
			Write("new cd", c.Cooldown(action.ActionSkill))
		return false
	}

	c.Core.Events.Subscribe(event.OnOverload, reduce, "dori-a1")
	c.Core.Events.Subscribe(event.OnElectroCharged, reduce, "dori-a1")
	c.Core.Events.Subscribe(event.OnSuperconduct, reduce, "dori-a1")
	//c.Core.Events.Subscribe(event.OnQuicken, reduce, "dori-a1") //TODO:save this for dendro branch folks
	//c.Core.Events.Subscribe(event.OnAggravate, reduce, "dori-a1")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, reduce, "dori-a1")
	c.Core.Events.Subscribe(event.OnSwirlElectro, reduce, "dori-a1")
}

// When the Troubleshooter Shots or After-Sales Service Rounds from Spirit-Warding Lamp: Troubleshooter Cannon hit opponents,
// Dori will restore 5 Elemental Energy for every 100% Energy Recharge possessed.
// Per Spirit-Warding Lamp: Troubleshooter Cannon, only one instance of Energy restoration can be triggered and a maximum of 15 Energy
//  can be restored this way.

func (c *char) a4() {
	//check a4CB on skill.go
}
