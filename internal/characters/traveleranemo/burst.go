package traveleranemo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(111)
	burstFrames[action.ActionAttack] = 106
	burstFrames[action.ActionSkill] = 104
	burstFrames[action.ActionDash] = 91
	burstFrames[action.ActionJump] = 91
	burstFrames[action.ActionSwap] = 96
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	//first hit at 94, then 30 frames between hits. 9 anemo hits total
	//yes the game description scams you on the duration
	duration := 94 + 30*8

	c.Core.Status.Add("amcburst", duration)

	c.qInfuse = attributes.NoElement
	c.qICDTag = combat.ICDTagNone
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gust Surge",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalArtAnemo,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	aiAbs := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gust Surge (Absorbed)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.NoElement,
		Durability: 50,
		Mult:       burstAbsorbDot[c.TalentLvlBurst()],
	}

	snapAbs := c.Snapshot(&aiAbs)

	var cb combat.AttackCBFunc
	if c.Base.Cons >= 6 {
		cb = c6cb(attributes.Anemo)
	}

	for i := 0; i < 9; i++ {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy), 94+30*i, cb)

		c.Core.Tasks.Add(func() {
			if c.qInfuse != attributes.NoElement {
				aiAbs.Element = c.qInfuse
				var cbAbs combat.AttackCBFunc
				if c.Base.Cons >= 6 {
					cbAbs = c6cb(c.qInfuse)
				}
				c.Core.QueueAttackWithSnap(aiAbs, snapAbs, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy), 0, cbAbs)
			}
			//check if infused
		}, 94+30*i)
	}

	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, int(duration/18)), 39)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 2)
	c.ConsumeEnergy(8)

	// TODO: Fill these out later
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qInfuse = c.Core.Combat.AbsorbCheck(c.infuseCheckLocation, attributes.Cryo, attributes.Pyro, attributes.Hydro, attributes.Electro)
		switch c.qInfuse {
		case attributes.Cryo:
			c.qICDTag = combat.ICDTagElementalBurstCryo
		case attributes.Pyro:
			c.qICDTag = combat.ICDTagElementalBurstPyro
		case attributes.Electro:
			c.qICDTag = combat.ICDTagElementalBurstElectro
		case attributes.Hydro:
			c.qICDTag = combat.ICDTagElementalBurstHydro
		case attributes.NoElement:
			//otherwise queue up
			c.Core.Tasks.Add(c.absorbCheckQ(src, count+1, max), 18)
		}
	}
}
