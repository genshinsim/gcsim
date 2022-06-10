package xiangling

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

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
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       guobaTick[c.TalentLvlSkill()],
	}

	// guoba spawns at cd frame
	c.Core.Status.Add("xianglingguoba", 500+13)

	//lasts 7.3 seconds, shoots every 100 frames
	delay := 126 //first tick at 126
	snap := c.Snapshot(&ai)
	for i := 0; i < 4; i++ {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), 10, c.c1)
			c.guoba.pyroWindowStart = c.Core.F
			c.guoba.pyroWindowEnd = c.Core.F + 20
		}, delay+i*100-10) //10 frame window to swirl
		//TODO: check guoba particle generation
		c.Core.QueueParticle("xiangling", 1, attributes.Pyro, delay+i*100+150)
	}

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		Post:            skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
