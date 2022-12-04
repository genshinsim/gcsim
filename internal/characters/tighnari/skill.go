package tighnari

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const (
	skillRelease           = 15
	vijnanasuffusionStatus = "vijnanasuffusion"
	wreatharrows           = "wreatharrows"
)

func init() {
	skillFrames = frames.InitAbilSlice(30)
	skillFrames[action.ActionAttack] = 20
	skillFrames[action.ActionAim] = 20
	skillFrames[action.ActionBurst] = 22
	skillFrames[action.ActionDash] = 23
	skillFrames[action.ActionJump] = 23
	skillFrames[action.ActionSwap] = 21
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 5
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Vijnana-Phala Mine",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			6,
		),
		skillRelease,
		skillRelease+travel,
	)

	var count float64 = 3
	if c.Core.Rand.Float64() < 0.5 {
		count++
	}
	c.Core.QueueParticle("tighnari", count, attributes.Dendro, skillRelease+travel+c.ParticleDelay)
	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	c.Core.Tasks.Add(func() {
		c.AddStatus(vijnanasuffusionStatus, 12*60, false)
		c.SetTag(wreatharrows, 3)
	}, 13)

	if c.Base.Cons >= 2 {
		c.QueueCharTask(c.c2, skillRelease+travel)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAim], // earliest cancel
		State:           action.SkillState,
	}
}
