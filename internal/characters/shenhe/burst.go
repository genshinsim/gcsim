package shenhe

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int

const burstHitmark = 99

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// TODO: Not 100% sure if this shares ICD with the DoT, currently coded that it does
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Divine Maiden's Deliverance (Initial)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	// assumes player is target 0
	x, y := c.Core.Combat.Target(0).Pos()

	//duration is 12 second (extended by c2 by 6s)
	dur := 12 * 60
	count := 6
	if c.Base.Cons >= 2 {
		dur += 6 * 60
		count += 3

		// Active characters within the skill's field deals 15% increased Cryo CRIT DMG.
		// TODO: Exact mechanics of how this works is unknown. Not sure if it works like Gorou E/Bennett Q
		// For now, assume that it operates like Kazuha C2, and extends for 2s after burst ends like the res shred
		m := make([]float64, attributes.EndStatType)
		m[attributes.CD] = 0.15
		for _, char := range c.Core.Player.Chars() {
			this := char
			char.AddAttackMod("shenhe-c2", dur+2*60, func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if ae.Info.Element != attributes.Cryo {
					return nil, false
				}

				switch this.Index {
				case c.Core.Player.Active(), c.Index:
					return m, true
				}
				return nil, false
			})
		}
	}
	// Res shred persists for 2 seconds after burst ends
	cb := func(a combat.AttackCB) {
		e, ok := a.Target.(core.Enemy)
		if !ok {
			return
		}
		e.AddResistMod("shenhe-burst-shred-cryo", 2*60, attributes.Cryo, -burstrespp[c.TalentLvlBurst()])
		e.AddResistMod("shenhe-burst-shred-phys", 2*60, attributes.Physical, -burstrespp[c.TalentLvlBurst()])
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(x, y, 2, false, combat.TargettableEnemy), 0, 15, cb)

	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Divine Maiden's Deliverance (DoT)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burstdot[c.TalentLvlBurst()],
	}

	mA1 := make([]float64, attributes.EndStatType)
	mA1[attributes.CryoP] = 0.15
	c.Core.Tasks.Add(func() {
		snap := c.Snapshot(&ai)
		c.Core.Status.Add("shenheburst", dur)
		// inspired from barbara c2
		// TODO: this isn't right.. it should only apply for active char
		// TODO: technically always assumes you are inside shenhe burst
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod("shenhe-a1", dur, attributes.CryoP, func() ([]float64, bool) {
				return mA1, true
			})
		}
		//TODO: check this accuracy? Siri's sheet has 137 per
		// dot every 2 second, double tick shortly after another
		for i := 0; i < count; i++ {
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(0, 0, 5, false, combat.TargettableEnemy), i*120+50)
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(0, 0, 5, false, combat.TargettableEnemy), i*120+80)
		}
	}, burstHitmark+2)

	c.SetCDWithDelay(action.ActionBurst, 20*60, 11)
	c.ConsumeEnergy(11)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		Post:            burstHitmark,
		State:           action.BurstState,
	}
}
