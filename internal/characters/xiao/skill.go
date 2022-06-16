package xiao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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

// Skill attack damage queue generator
// Additionally implements A4
// Using Lemniscatic Wind Cycling increases the DMG of subsequent uses of Lemniscatic Wind Cycling by 15%. This effect lasts for 7s and has a maximum of 3 stacks. Gaining a new stack refreshes the duration of this effect.
func (c *char) Skill(p map[string]int) action.ActionInfo {

	// Add damage based on A4
	if c.a4Expiry <= c.Core.F {
		c.Tags["a4"] = 0
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lemniscatic Wind Cycling",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupXiaoDash,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(2, false, combat.TargettableEnemy), skillHitmark)

	// Text is not explicit, but assume that gaining a stack while at max still refreshes duration
	c.Tags["a4"]++
	c.a4Expiry = c.Core.F + 420
	if c.Tags["a4"] > 3 {
		c.Tags["a4"] = 3
	}

	// Cannot create energy during burst uptime
	if c.Core.Status.Duration("xiaoburst") > 0 {
	} else {
		c.Core.QueueParticle("xiao", 3, attributes.Anemo, skillHitmark+100)
	}

	// C6 handling - can use skill ignoring CD and without draining charges
	// Can simply return early
	if c.Base.Cons < 6 || c.Core.Status.Duration("xiaoc6") == 0 {
		c.SetCD(action.ActionSkill, 600)
	} else {
		c.Core.Log.NewEvent("xiao c6 active, Xiao E used, no charge used, no CD", glog.LogCharacterEvent, c.Index, "c6 remaining duration", c.Core.Status.Duration("xiaoc6"))
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSkill], // earliest cancel
		State:           action.SkillState,
	}
}
