package diona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const burstStart = 58 // Initial Hit

func init() {
	burstFrames = frames.InitAbilSlice(64) // Q -> N1/E
	burstFrames[action.ActionDash] = 43    // Q -> D
	burstFrames[action.ActionJump] = 44    // Q -> J
	burstFrames[action.ActionSwap] = 41    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	// Initial Hit
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Signature Mix (Initial)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 3), 0, burstStart)

	// Ticks
	ai.Abil = "Signature Mix (Tick)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6.5)

	snap := c.Snapshot(&ai)
	hpplus := snap.Stats[attributes.Heal]
	maxhp := c.MaxHP()
	heal := burstHealPer[c.TalentLvlBurst()]*maxhp + burstHealFlat[c.TalentLvlBurst()]

	c.burstBuffArea = combat.NewCircleHitOnTarget(ap.Shape.Pos(), nil, 7)
	// apparently lasts for 12.5
	// TODO: assumes that field starts when it lands (which is dynamic ingame)
	c.Core.Tasks.Add(func() {
		// add burst status for C4 check
		c.Core.Status.Add("diona-q", 750)
		//ticks every 2s, first tick at t=2s (relative to field start), then t=4,6,8,10,12; lasts for 12.5s from field start
		for i := 0; i < 6; i++ {
			c.Core.Tasks.Add(func() {
				// attack
				c.Core.QueueAttackWithSnap(ai, snap, ap, 0)
				// heal
				if !c.Core.Combat.Player().IsWithinArea(c.burstBuffArea) {
					return
				}
				c.Core.Player.Heal(player.HealInfo{
					Caller:  c.Index,
					Target:  c.Core.Player.Active(),
					Message: "Drunken Mist",
					Src:     heal,
					Bonus:   hpplus,
				})
			}, 120+i*120)
		}
		// C6
		if c.Base.Cons >= 6 {
			c.c6()
		}
	}, burstStart)

	// C1
	if c.Base.Cons >= 1 {
		//15 energy after ends, flat not affected by ER
		c.Core.Tasks.Add(func() {
			c.AddEnergy("diona-c1", 15)
		}, burstStart+750)
	}

	c.SetCDWithDelay(action.ActionBurst, 1200, 41)
	c.ConsumeEnergy(43)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstStart,
		State:           action.BurstState,
	}
}
