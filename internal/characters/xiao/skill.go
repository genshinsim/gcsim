package xiao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillHitmark = 4

func init() {
	skillFrames = frames.InitAbilSlice(37)
	skillFrames[action.ActionAttack] = 24
	skillFrames[action.ActionSkill] = 24
	skillFrames[action.ActionBurst] = 24
	skillFrames[action.ActionDash] = 35
	skillFrames[action.ActionSwap] = 35
}

const a4BuffKey = "xiao-a4"

// Skill attack damage queue generator
// Additionally implements A4
// Using Lemniscatic Wind Cycling increases the DMG of subsequent uses of Lemniscatic Wind Cycling by 15%. This effect lasts for 7s and has a maximum of 3 stacks. Gaining a new stack refreshes the duration of this effect.
func (c *char) Skill(p map[string]int) action.ActionInfo {

	// Add damage based on A4
	if !c.StatModIsActive(a4BuffKey) {
		c.a4stacks = 0
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lemniscatic Wind Cycling",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupXiaoDash,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.8,
		),
		skillHitmark,
	)

	// apply A4 0.25s after cast
	c.Core.Tasks.Add(func() {
		// Text is not explicit, but assume that gaining a stack while at max still refreshes duration
		c.a4stacks++
		if c.a4stacks > 3 {
			c.a4stacks = 3
		}
		c.a4buff[attributes.DmgP] = float64(c.a4stacks) * 0.15
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(a4BuffKey, 420),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				return c.a4buff, atk.Info.AttackTag == combat.AttackTagElementalArt
			},
		})
	}, 15)

	// Cannot create energy during burst uptime
	if !c.StatusIsActive(burstBuffKey) {
		c.Core.QueueParticle("xiao", 3, attributes.Anemo, skillHitmark+c.ParticleDelay)
	}

	// C6 handling - can use skill ignoring CD and without draining charges
	// Can simply return early
	if c.Base.Cons >= 6 && c.StatusIsActive(c6BuffKey) {
		c.Core.Log.NewEvent("xiao c6 active, Xiao E used, no charge used, no CD", glog.LogCharacterEvent, c.Index).
			Write("c6 remaining duration", c.Core.Status.Duration("xiaoc6"))
	} else {
		c.SetCD(action.ActionSkill, 600)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSkill], // earliest cancel
		State:           action.SkillState,
	}
}
