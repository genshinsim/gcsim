package wriothesley

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillKey          = "wriothesley-e"
	skillCDDelay      = 1
	skillLoseHPICDKey = "wriothesley-lose-hp-icd"
	skillLoseHPICD    = 6
	particleICDKey    = "wriothesley-particle-icd"
	particleICD       = 2 * 60
)

func init() {
	skillFrames = frames.InitAbilSlice(29)
	skillFrames[action.ActionAttack] = 15
	skillFrames[action.ActionBurst] = 5
	skillFrames[action.ActionDash] = 21
	skillFrames[action.ActionJump] = 20
	skillFrames[action.ActionSwap] = 19
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// anything but NA/E -> E should reset savedNormalCounter
	// can't use CurrentState here since AnimationLength of Dash is the same as Dash -> Skill, so it switches to Idle instead of staying DashState
	switch c.Core.Player.LastAction.Type {
	case action.ActionAttack:
	case action.ActionSkill:
	default:
		c.savedNormalCounter = 0
	}

	c.resetA4()
	c.resetC1SkillExtension()

	c.AddStatus(skillKey, 10*60+skillCDDelay, true)
	c.applyA4(10*60 + skillCDDelay)
	c.SetCDWithDelay(action.ActionSkill, 16*60, skillCDDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Cryo, c.ParticleDelay)
}

func (c *char) chillingPenalty(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(skillLoseHPICDKey) {
		return
	}
	c.AddStatus(skillLoseHPICDKey, skillLoseHPICD, true)
	c.Core.Player.Drain(player.DrainInfo{
		ActorIndex: c.Index,
		Abil:       "Chilling Penalty",
		Amount:     0.045 * c.MaxHP(),
	})
}

func (c *char) onExit() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if !c.StatusIsActive(skillKey) {
			return false
		}
		prev := args[0].(int)
		if prev == c.Index {
			c.DeleteStatus(skillKey)
		}
		return false
	}, "wriothesley-exit")
}
