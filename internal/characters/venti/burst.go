package venti

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstStart = 94

func init() {
	burstFrames = frames.InitAbilSlice(95) // Q -> N1/CA/E/D
	burstFrames[action.ActionJump] = 94    // Q -> J
	burstFrames[action.ActionSwap] = 93    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// reset location
	c.qAbsorb = attributes.NoElement
	player := c.Core.Combat.Player()
	c.qPos = combat.CalcOffsetPoint(player.Pos(), combat.Point{Y: 5}, player.Direction())
	c.absorbCheckLocation = combat.NewBoxHitOnTarget(c.qPos, combat.Point{Y: -1}, 2.5, 2.5)

	//8 second duration, tick every .4 second
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Wind's Grand Ode",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurstAnemo,
		ICDGroup:   combat.ICDGroupVenti,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.qPos, nil, 4)

	c.aiAbsorb = ai
	c.aiAbsorb.Abil = "Wind's Grand Ode (Absorbed)"
	c.aiAbsorb.Mult = burstAbsorbDot[c.TalentLvlBurst()]
	c.aiAbsorb.Element = attributes.NoElement

	// snapshot is around cd frame and 1st tick?
	var snap combat.Snapshot
	c.Core.Tasks.Add(func() {
		snap = c.Snapshot(&ai)
		c.snapAbsorb = c.Snapshot(&c.aiAbsorb)
	}, 104)

	var cb combat.AttackCBFunc
	if c.Base.Cons >= 6 {
		cb = c.c6(attributes.Anemo)
	}

	// starts at 106 with 24f interval between ticks. 20 total
	for i := 0; i < 20; i++ {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackWithSnap(ai, snap, ap, 0, cb)
		}, 106+24*i)
	}
	// Infusion usually occurs after 4 ticks of anemo according to KQM library
	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, int((480-24*4)/18)), 106+24*3)

	// a4: restore 15 energy on burst end
	c.Core.Tasks.Add(func() {
		c.a4()
	}, 480+burstStart)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 81)
	c.ConsumeEnergy(84)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstAbsorbedTicks() {
	var cb combat.AttackCBFunc
	if c.Base.Cons >= 6 {
		cb = c.c6(c.qAbsorb)
	}

	ap := combat.NewCircleHitOnTarget(c.qPos, nil, 6)
	// ticks at 24f. 15 total
	for i := 0; i < 15; i++ {
		c.Core.QueueAttackWithSnap(c.aiAbsorb, c.snapAbsorb, ap, i*24, cb)
	}
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qAbsorb = c.Core.Combat.AbsorbCheck(c.absorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)
		if c.qAbsorb != attributes.NoElement {
			c.aiAbsorb.Element = c.qAbsorb
			switch c.qAbsorb {
			case attributes.Pyro:
				c.aiAbsorb.ICDTag = combat.ICDTagElementalBurstPyro
			case attributes.Hydro:
				c.aiAbsorb.ICDTag = combat.ICDTagElementalBurstHydro
			case attributes.Electro:
				c.aiAbsorb.ICDTag = combat.ICDTagElementalBurstElectro
			case attributes.Cryo:
				c.aiAbsorb.ICDTag = combat.ICDTagElementalBurstCryo
			}
			//trigger dmg ticks here
			c.burstAbsorbedTicks()
			return
		}
		//otherwise queue up
		c.Core.Tasks.Add(c.absorbCheckQ(src, count+1, max), 18)
	}
}
