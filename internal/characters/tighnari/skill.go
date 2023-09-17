package tighnari

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillRelease           = 15
	particleICDKey         = "tighnari-particle-icd"
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

func (c *char) Skill(p map[string]int) action.Info {
	travel, ok := p["travel"]
	if !ok {
		travel = 5
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Vijnana-Phala Mine",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.skillArea = combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6)
	c.Core.QueueAttack(ai, c.skillArea, skillRelease, skillRelease+travel, c.particleCB)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	c.Core.Tasks.Add(func() {
		c.AddStatus(vijnanasuffusionStatus, 12*60, false)
		c.SetTag(wreatharrows, 3)
	}, 13)

	if c.Base.Cons >= 2 {
		c.QueueCharTask(c.c2, skillRelease+travel)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAim], // earliest cancel
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
	c.AddStatus(particleICDKey, 0.3*60, true)

	count := 3.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 4
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Dendro, c.ParticleDelay)
}
