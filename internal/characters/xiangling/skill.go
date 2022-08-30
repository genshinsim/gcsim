package xiangling

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       guobaTick[c.TalentLvlSkill()],
	}

	// guoba spawns at cd frame
	c.Core.Status.Add("xianglingguoba", 500+13)

	// lasts 7.3 seconds, shoots every 100 frames
	delay := 126 // first tick at 126
	snap := c.Snapshot(&ai)
	for i := 0; i < 4; i++ {
		c.Core.Tasks.Add(func() {
			done := false
			part := func(_ combat.AttackCB) {
				if done {
					return
				}
				done = true
				c.Core.QueueParticle("xiangling", 1, attributes.Pyro, c.ParticleDelay)
			}
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewCircleHit(c.Core.Combat.Player(), 0.5, false, combat.TargettableEnemy),
				10,
				c.c1,
				part,
			)
			c.Core.Log.NewEventBuildMsg(glog.LogElementEvent, c.Index, "guoba self infusion applied").
				SetEnded(c.Core.F+infuseWindow+1)
			c.guoba.Durability[attributes.Pyro] = infuseDurability
			c.Core.Tasks.Add(func() {
				c.guoba.Durability[attributes.Pyro] = 0
			}, infuseWindow+1) // +1 since infuse window is inclusive
		}, delay+i*100-10) // 10 frame window to swirl
	}

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
