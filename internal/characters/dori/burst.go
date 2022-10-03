package dori

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const (
	burstHitmark = 28
)

func init() {
	burstFrames = frames.InitAbilSlice(58) // Q
	burstFrames[action.ActionAttack] = 57  // Q -> N1
	burstFrames[action.ActionSkill] = 57   // Q -> E
	burstFrames[action.ActionJump] = 57    // Q -> J
	burstFrames[action.ActionSwap] = 56    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Alcazarzaray's Exactitude: Connector DMG",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDoriBurst,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	icdSrc := []int{math.MinInt32, math.MinInt32, math.MinInt32, math.MinInt32}
	// 32 damage ticks
	for i := 0; i < 32; i++ {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewDefBoxHit(1, -2, false, combat.TargettableEnemy, combat.TargettableGadget),
				0,
			) // TODO: accurate hitbox

			// dori self application
			c.Core.Events.Emit(event.OnCharacterHurt, 0)
			p, ok := c.Core.Combat.Player().(*avatar.Player)
			if !ok {
				panic("target 0 should be Player but is not!!")
			}
			idx := c.Core.Player.ActiveChar().Index
			if c.Core.F > icdSrc[idx]+combat.ICDGroupResetTimer[combat.ICDGroupDoriBurst] {
				dur := combat.Durability(25)
				if p.AuraCount() == 0 {
					dur = 20
				}
				p.ApplySelfInfusion(attributes.Electro, dur, 9.5*60) // TODO: find actual duration
				icdSrc[idx] = c.Core.F
			}

			if c.Base.Cons >= 4 {
				c.c4()
			}
		}, 24*i+burstHitmark)
	}

	c2Travel, ok := p["c2_travel"]
	if !ok {
		c2Travel = 10
	}

	for i := 0; i < 6; i++ {
		c.Core.Tasks.Add(func() {
			if c.Base.Cons >= 2 {
				c.c2(c2Travel)
			}
			// Heals
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "Alcazarzaray's Exactitude: Healing",
				Src:     bursthealpp[c.TalentLvlBurst()]*c.MaxHP() + bursthealflat[c.TalentLvlBurst()],
				Bonus:   snap.Stats[attributes.Heal],
			})
			// Energy regen to active char
			active := c.Core.Player.ActiveChar()
			active.AddEnergy("Alcazarzaray's Exactitude: Energy Regen", burstenergy[c.TalentLvlBurst()])
		}, 120*i+11)
	}

	c.ConsumeEnergy(4)
	c.SetCDWithDelay(action.ActionBurst, 1200, 1) // 20s * 60

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
