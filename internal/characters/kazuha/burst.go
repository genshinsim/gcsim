package kazuha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const (
	burstHitmark   = 82
	burstFirstTick = 140
)

const burstStatus = "kazuha-q"

func init() {
	burstFrames = frames.InitAbilSlice(93) // Q -> J
	burstFrames[action.ActionAttack] = 92  // Q -> N1
	burstFrames[action.ActionSkill] = 92   // Q -> E
	burstFrames[action.ActionDash] = 92    // Q -> D
	burstFrames[action.ActionSwap] = 90    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	player := c.Core.Combat.Player()
	c.qAbsorb = attributes.NoElement
	c.qAbsorbCheckLocation = combat.NewCircleHitOnTarget(player, combat.Point{Y: 1}, 8)

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Kazuha Slash",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Anemo,
		Durability:         50,
		Mult:               burstSlash[c.TalentLvlBurst()],
		HitlagHaltFrames:   0.05 * 60,
		HitlagFactor:       0.05,
		CanBeDefenseHalted: false,
	}
	ap := combat.NewCircleHitOnTarget(player, combat.Point{Y: 1}, 9)

	c.Core.QueueAttack(ai, ap, 0, burstHitmark)

	//apply dot and check for absorb
	ai.Abil = "Kazuha Slash (Dot)"
	ai.StrikeType = combat.StrikeTypeDefault
	ai.Mult = burstDot[c.TalentLvlBurst()]
	ai.Durability = 25
	// no more hitlag after initial slash
	ai.HitlagHaltFrames = 0
	snap := c.Snapshot(&ai)

	aiAbsorb := ai
	aiAbsorb.Abil = "Kazuha Slash (Absorb Dot)"
	aiAbsorb.Mult = burstEleDot[c.TalentLvlBurst()]
	aiAbsorb.Element = attributes.NoElement
	snapAbsorb := c.Snapshot(&aiAbsorb)

	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, int(310/18)), 10)

	// make sure that this task gets executed:
	// - inside Q hitlag
	// - before kazuha can get affected by any more hitlag
	c.QueueCharTask(func() {
		// queue up ticks
		// from kisa's count: ticks starts at 147, + 117 gap each roughly; 5 ticks total
		// updated to 140 based on koli's count: https://docs.google.com/spreadsheets/d/1uEbP13O548-w_nGxFPGsf5jqj1qGD3pqFZ_AiV4w3ww/edit#gid=775340159
		for i := 0; i < 5; i++ {
			c.Core.Tasks.Add(func() {
				if c.qAbsorb != attributes.NoElement {
					aiAbsorb.Element = c.qAbsorb
					c.Core.QueueAttackWithSnap(aiAbsorb, snapAbsorb, ap, 0)
				}
				c.Core.QueueAttackWithSnap(ai, snap, ap, 0)
			}, (burstFirstTick-(burstHitmark+1))+117*i)
		}
		// C2:
		// TODO: Not sure when it lasts from and until
		// -> For now, assume that it lasts from Initial Hit hitlag end to the last Q tick.
		// TODO: Does it apply to Kazuha's initial hit?
		// -> For now, assume that it doesn't.
		c.Core.Status.Add(burstStatus, (burstFirstTick-(burstHitmark+1))+117*4)
		if c.Base.Cons >= 2 {
			c.qFieldSrc = c.Core.F
			c.Core.Tasks.Add(c.c2(c.Core.F), 30) // start checking in 0.5s
		}
		// C6:
		// TODO: when does the infusion kick in?
		// -> For now, assume that it starts on Initial Hit hitlag end.
		if c.Base.Cons >= 6 {
			c.c6()
		}
	}, burstHitmark+1)

	//reset skill cd
	if c.Base.Cons >= 1 {
		c.ResetActionCooldown(action.ActionSkill)
	}

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(4)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) absorbCheckQ(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.qAbsorb = c.Core.Combat.AbsorbCheck(c.qAbsorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

		if c.qAbsorb != attributes.NoElement {
			return
		}
		//otherwise queue up
		c.Core.Tasks.Add(c.absorbCheckQ(src, count+1, max), 18)
	}
}
