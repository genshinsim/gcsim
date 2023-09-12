package sucrose

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const particleICDKey = "sucrose-particle-icd"

func init() {
	skillFrames = frames.InitAbilSlice(57)
	skillFrames[action.ActionCharge] = 56
	skillFrames[action.ActionSkill] = 56
	skillFrames[action.ActionDash] = 11
	skillFrames[action.ActionJump] = 11
	skillFrames[action.ActionSwap] = 56
}

func (c *char) Skill(p map[string]int) action.Info {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Astable Anemohypostasis Creation-6308",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	done := false
	a4CB := func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.a4()
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 5}, 6),
		0,
		42,
		a4CB,
		c.particleCB,
	)

	// reduce charge by 1
	c.SetCDWithDelay(action.ActionSkill, 900, 9)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.4*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Anemo, c.ParticleDelay)
}
