package yoimiya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var skillFrames []int

const (
	skillKey       = "yoimiyaskill"
	particleICDKey = "yoimiya-particle-icd"
	skillStart     = 11
)

func init() {
	skillFrames = frames.InitAbilSlice(34)
	skillFrames[action.ActionAttack] = 22
	skillFrames[action.ActionAim] = 22 // uses attack frames
	skillFrames[action.ActionBurst] = 23
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 31
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	c.AddStatus(skillKey, 600+skillStart, true) // activate for 10
	if !c.StatusIsActive(a1Key) {
		c.a1Stacks = 0
	}

	c.SetCDWithDelay(action.ActionSkill, 1080, 11)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Pyro, c.ParticleDelay)
}

func (c *char) onExit() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		next := args[1].(int)
		if prev == c.Index && next != c.Index {
			c.DeleteStatus(skillKey)
		}
		return false
	}, "yoimiya-exit")
}
