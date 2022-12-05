package xiangling

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	infuseWindow     = 30
	infuseDurability = 20
)

func init() {
	skillFrames = frames.InitAbilSlice(39)
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionJump] = 14
	skillFrames[action.ActionSwap] = 38
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Guoba",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       guobaTick[c.TalentLvlSkill()],
	}

	// guoba spawns at cd frame
	c.Core.Status.Add("xianglingguoba", 500+13)

	// lasts 7.3 seconds, shoots every 100 frames
	// delay := 126 // first tick at 126
	// snap := c.Snapshot(&ai)
	guoba := c.newGuoba(ai)
	c.Core.Tasks.Add(func() {
		c.Core.Combat.AddGadget(guoba)
	}, 13)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
