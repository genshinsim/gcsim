package heizou

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	//tag a4
	//first hit at 137, then 113 frames between hits

	duration := 360
	if c.Base.Cons >= 2 {
		duration = 480
	}

	c.burstTaggedCount = 0

	c.Core.Status.AddStatus("heizouburst", duration)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Fudou Style Vacuum Slugger",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	//TODO: does heizou burst snapshot?
	snap := c.Snapshot(&ai)

	burstCB := func(a core.AttackCB) {
		//check if enemy
		if a.Target.Type() != core.TargettableEnemy {
			return
		}
		//max 4 tagged
		if c.burstTaggedCount == 4 {
			return
		}
		//check for element and queue attack
		c.burstTaggedCount++
		if c.Base.Cons >= 4 {
			c.c4(c.burstTaggedCount)
		}
		c.irisDmg(a.Target)
	}

	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), f, burstCB)

	//TODO: Check CD with or without delay, check energy consume frame
	c.SetCD(core.ActionBurst, 720)
	c.ConsumeEnergy(21)
	return f, a
}

//When Vacuum Slugger hits opponents affected by Hydro/Pyro/Cryo/Electro,
//these opponents will be afflicted with Windmuster Iris.
//This Windmuster Iris will explode after a moment and dissipate,
//dealing AoE DMG of the corresponding aforementioned elemental type.
func (c *char) irisDmg(t core.Target) {

	//TODO: does burst iris snapshot
	aiAbs := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Windmuster Iris",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.NoElement,
		Durability: 25,
		Mult:       burstIris[c.TalentLvlBurst()],
	}
	//TODO: Iris timing; looks to be 0.6s after hitmark
	x, y := t.Shape().Pos()

	switch ele := t.AuraType(); ele {
	case core.Pyro, core.Hydro, core.Electro, core.Cryo:
		aiAbs.Element = ele
	case core.EC:
		aiAbs.Element = core.Hydro
	case core.Frozen:
		aiAbs.Element = core.Cryo
	default:
		c.Core.Log.NewEvent(
			"No valid aura detected, omiting iris",
			core.LogCharacterEvent,
			c.Index,
			"aura type", t.AuraType(),
		)
		return
	}

	c.Core.Combat.QueueAttack(aiAbs, core.NewCircleHit(x, y, 2.5, false, core.TargettableEnemy), 1, 1)

}
