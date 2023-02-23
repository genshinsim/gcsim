package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	skillHitmark   = 33
	particleICDKey = "ayaka-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(49)
	skillFrames[action.ActionBurst] = 48
	skillFrames[action.ActionDash] = 30
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 48
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Hyouka",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.a1()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4.5),
		0,
		skillHitmark,
		c.particleCB,
	)

	c.SetCD(action.ActionSkill, 600)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
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
	c.AddStatus(particleICDKey, 0.3*60, false)

	count := 4.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 5
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Cryo, c.ParticleDelay)
}
