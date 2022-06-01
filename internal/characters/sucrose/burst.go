package sucrose

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//tag a4
	//first hit at 137, then 113 frames between hits
	duration := 360
	if c.Base.Cons >= 2 {
		duration = 480
	}

	c.qInfused = attributes.NoElement

	c.Core.Status.Add("sucroseburst", duration)
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Forbidden Creation-Isomer 75/Type II",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	//TODO: does sucrose burst snapshot?
	snap := c.Snapshot(&ai)
	//TODO: does burst absorb snapshot
	aiAbs := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Forbidden Creation-Isomer 75/Type II (Absorb)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.NoElement,
		Durability: 25,
		Mult:       burstAbsorb[c.TalentLvlBurst()],
	}
	snapAbs := c.Snapshot(&aiAbs)

	lockout := 0
	cb := func(a combat.AttackCB) {
		//lockout for 1 frame to prevent triggering multiple times on one attack
		if lockout > c.Core.F {
			return
		}
		lockout = c.Core.F + 1
		c.a4()
	}

	for i := 137; i <= duration+5; i += 113 {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(5, false, combat.TargettableEnemy), i, cb)

		c.Core.Tasks.Add(func() {
			if c.qInfused != attributes.NoElement {
				aiAbs.Element = c.qInfused
				c.Core.QueueAttackWithSnap(aiAbs, snapAbs, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0)
			}
			//check if infused
		}, i)
	}

	c.Core.Tasks.Add(c.absorbCheck(c.Core.F, 0, int(duration/18)), 136)

	c.SetCDWithDelay(action.ActionBurst, 1200, 18)
	c.ConsumeEnergy(21)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		Post:            burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) absorbCheck(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfused = c.Core.AbsorbCheck(c.infuseCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

		if c.qInfused != attributes.NoElement {
			if c.Base.Cons >= 6 {
				c.c6()
			}
			return
		}
		//otherwise queue up
		c.Core.Tasks.Add(c.absorbCheck(src, count+1, max), 18)
	}
}
