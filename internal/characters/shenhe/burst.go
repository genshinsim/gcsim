package shenhe

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int
var burstTickOffset = []int{0, 2, 4, 0, 2, 4, 0, 2, 4}

const (
	burstStart   = 47
	burstHitmark = 78
	burstKey     = "shenheburst"
)

func init() {
	burstFrames = frames.InitAbilSlice(100) // Q -> E
	burstFrames[action.ActionAttack] = 99   // Q -> N1
	burstFrames[action.ActionDash] = 78     // Q -> D
	burstFrames[action.ActionJump] = 79     // Q -> J
	burstFrames[action.ActionWalk] = 98     // Q -> Walk
	burstFrames[action.ActionSwap] = 98     // Q -> Swap
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Divine Maiden's Deliverance (Initial)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2), burstHitmark, burstHitmark)

	// duration is 12 second (extended by c2 by 6s)
	count := 6
	burstDuration := 12 * 60
	if c.Base.Cons >= 2 {
		count += 3
		burstDuration += 6 * 60
	}
	c.AddStatus(burstKey, burstDuration, false)

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
	// DoT snapshot before A1
	c.Core.Tasks.Add(func() {
		snap := c.Snapshot(&ai)
		for i := 0; i < count; i++ {
			hitmark := 82 + i*117
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 5), hitmark)
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 5), hitmark+30+burstTickOffset[i])
		}
	}, burstStart)

	// assumes player/target is inside shenhe burst
	for i := burstStart; i < burstStart+burstDuration; i += 18 {
		c.Core.Tasks.Add(func() {
			buffDuration := 38
			active := c.Core.Player.ActiveChar()

			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("shenhe-a1", buffDuration),
				AffectedStat: attributes.CryoP,
				Amount: func() ([]float64, bool) {
					return c.burstBuff, true
				},
			})
			if c.Base.Cons >= 2 {
				c.c2(active, buffDuration)
			}

			for _, t := range c.Core.Combat.Enemies() {
				// skip non-enemy targets
				e, ok := t.(*enemy.Enemy)
				if !ok {
					continue
				}
				e.AddResistMod(enemy.ResistMod{
					Base:  modifier.NewBaseWithHitlag("shenhe-burst-shred-cryo", buffDuration),
					Ele:   attributes.Cryo,
					Value: -burstrespp[c.TalentLvlBurst()],
				})
				e.AddResistMod(enemy.ResistMod{
					Base:  modifier.NewBaseWithHitlag("shenhe-burst-shred-phys", buffDuration),
					Ele:   attributes.Physical,
					Value: -burstrespp[c.TalentLvlBurst()],
				})
			}

		}, i)
	}
	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(4)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
