package anemo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
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

func (c *Traveler) Burst(p map[string]int) (action.Info, error) {
	// first hit at 94, then 30 frames between hits. 9 anemo hits total
	// yes the game description scams you on the duration
	duration := burstHitmarks[c.gender] + 30*8

	c.Core.Status.Add("amcburst", duration)

	c.qAbsorb = attributes.NoElement
	c.qICDTag = attacks.ICDTagNone
	c.qAbsorbCheckLocation = combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1.5}, 2.5, 2.5)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gust Surge",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalArtAnemo,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	ap := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -1.5}, 3, 3)
	snap := c.Snapshot(&ai)

	aiAbs := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gust Surge (Absorbed)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.NoElement,
		Durability: 50,
		Mult:       burstAbsorbDot[c.TalentLvlBurst()],
	}
	apAbs := combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -1}, 2.5, 2.5)

	snapAbs := c.Snapshot(&aiAbs)

	var cb combat.AttackCBFunc
	if c.Base.Cons >= 6 {
		cb = c6cb(attributes.Anemo)
	}

	for i := 0; i < 9; i++ {
		c.Core.QueueAttackWithSnap(ai, snap, ap, 94+30*i, cb)

		c.Core.Tasks.Add(func() {
			if c.qAbsorb != attributes.NoElement {
				aiAbs.Element = c.qAbsorb
				var cbAbs combat.AttackCBFunc
				if c.Base.Cons >= 6 {
					cbAbs = c6cb(c.qAbsorb)
				}
				c.Core.QueueAttackWithSnap(aiAbs, snapAbs, apAbs, 0, cbAbs)
			}
			// check if infused
		}, 94+30*i)
	}

	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, duration/18), 39)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(3)

	// TODO: Fill these out later
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *Traveler) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qAbsorb = c.Core.Combat.AbsorbCheck(c.qAbsorbCheckLocation, attributes.Cryo, attributes.Pyro, attributes.Hydro, attributes.Electro)
		switch c.qAbsorb {
		case attributes.Cryo:
			c.qICDTag = attacks.ICDTagElementalBurstCryo
		case attributes.Pyro:
			c.qICDTag = attacks.ICDTagElementalBurstPyro
		case attributes.Electro:
			c.qICDTag = attacks.ICDTagElementalBurstElectro
		case attributes.Hydro:
			c.qICDTag = attacks.ICDTagElementalBurstHydro
		case attributes.NoElement:
			// otherwise queue up
			c.Core.Tasks.Add(c.absorbCheckQ(src, count+1, max), 18)
		}
	}
}
