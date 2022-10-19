package traveleranemo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstHitmarks = []int{96, 94}
var burstFrames [][]int

func init() {
	burstFrames = make([][]int, 2)

	// Male
	burstFrames[0] = frames.InitAbilSlice(110) // Q -> N1
	burstFrames[0][action.ActionSkill] = 109   // Q -> E
	burstFrames[0][action.ActionDash] = 96     // Q -> D
	burstFrames[0][action.ActionJump] = 96     // Q -> J
	burstFrames[0][action.ActionSwap] = 100    // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(105) // Q -> N1
	burstFrames[1][action.ActionSkill] = 104   // Q -> E
	burstFrames[1][action.ActionDash] = 90     // Q -> D
	burstFrames[1][action.ActionJump] = 90     // Q -> J
	burstFrames[1][action.ActionSwap] = 95     // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	//first hit at 94, then 30 frames between hits. 9 anemo hits total
	//yes the game description scams you on the duration
	duration := burstHitmarks[c.gender] + 30*8

	c.Core.Status.Add("amcburst", duration)

	c.qAbsorb = attributes.NoElement
	c.qICDTag = combat.ICDTagNone
	self_absorb, ok := p["self_absorb"]
	if !ok {
		self_absorb = 0
	}
	if self_absorb == 0 {
		c.absorbCheckLocation = combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 0.1, false, combat.TargettableEnemy, combat.TargettableGadget)
	} else {
		c.absorbCheckLocation = combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 0.1, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableGadget)
	}

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
			if c.qAbsorb != attributes.NoElement {
				aiAbs.Element = c.qAbsorb
				var cbAbs combat.AttackCBFunc
				if c.Base.Cons >= 6 {
					cbAbs = c6cb(c.qAbsorb)
				}
				c.Core.QueueAttackWithSnap(aiAbs, snapAbs, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy, combat.TargettableGadget), 0, cbAbs)
			}
			//check if infused
		}, 94+30*i)
	}

	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, int(duration/18)), 39)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(3)

	// TODO: Fill these out later
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qAbsorb = c.Core.Combat.AbsorbCheck(c.absorbCheckLocation, attributes.Cryo, attributes.Pyro, attributes.Hydro, attributes.Electro)
		switch c.qAbsorb {
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
