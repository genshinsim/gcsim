package dori

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

// After a character connected to the Jinni triggers an Electro-Charged, Superconduct, Overloaded, Quicken, Aggravate, Hyperbloom,
// or an Electro Swirl or Crystallize reaction, the CD of Spirit-Warding Lamp: Troubleshooter Cannon is decreased by 1s.
// This effect can be triggered once every 3s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	const icdKey = "dori-a1"
	icd := 180 // 3s * 60
	//nolint:unparam // ignoring for now, event refactor should get rid of bool return of event sub
	reduce := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		if c.Core.Player.Active() != atk.Info.ActorIndex { // only for on field character
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

	reduceNoGadget := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		return reduce(args...)
	}

	c.Core.Events.Subscribe(event.OnOverload, reduceNoGadget, "dori-a1")
	c.Core.Events.Subscribe(event.OnElectroCharged, reduceNoGadget, "dori-a1")
	c.Core.Events.Subscribe(event.OnSuperconduct, reduceNoGadget, "dori-a1")
	c.Core.Events.Subscribe(event.OnQuicken, reduceNoGadget, "dori-a1")
	c.Core.Events.Subscribe(event.OnAggravate, reduceNoGadget, "dori-a1")
	c.Core.Events.Subscribe(event.OnHyperbloom, reduce, "dori-a1")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, reduceNoGadget, "dori-a1")
	c.Core.Events.Subscribe(event.OnSwirlElectro, reduceNoGadget, "dori-a1")
}

// When the Troubleshooter Shots or After-Sales Service Rounds from Spirit-Warding Lamp: Troubleshooter Cannon hit opponents,
// Dori will restore 5 Elemental Energy for every 100% Energy Recharge possessed.
// Per Spirit-Warding Lamp: Troubleshooter Cannon, only one instance of Energy restoration can be triggered
// and a maximum of 15 Energy can be restored this way.
func (c *char) makeA4CB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}

	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		a4Energy := a.AttackEvent.Snapshot.Stats[attributes.ER] * 5
		if a4Energy > 15 {
			a4Energy = 15
		}
		c.AddEnergy("dori-a4", a4Energy)
	}
}
