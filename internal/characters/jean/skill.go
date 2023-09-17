package jean

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
	skillHitmark        = 21
	baseParticleICDKey  = "jean-base-particle-icd"
	extraParticleICDKey = "jean-extra-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(46)
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionJump] = 28
	skillFrames[action.ActionSwap] = 45
}

func (c *char) Skill(p map[string]int) action.Info {
	hold := p["hold"]
	// hold for p up to 5 seconds
	if hold > 300 {
		hold = 300
	}
	hitmark := skillHitmark + hold

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gale Blade",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	if c.Base.Cons >= 1 && p["hold"] >= 60 {
		c.c1(&snap)
	}

	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 4, 4.1),
		hitmark,
		c.baseParticleCB,
		c.extraParticleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 360, hitmark-2)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] + hold },
		AnimationLength: skillFrames[action.InvalidAction] + hold,
		CanQueueAfter:   skillFrames[action.ActionDash] + hold, // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) baseParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(baseParticleICDKey) {
		return
	}
	c.AddStatus(baseParticleICDKey, 0.3*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Anemo, c.ParticleDelay)
}

func (c *char) extraParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(extraParticleICDKey) {
		return
	}
	c.AddStatus(extraParticleICDKey, 1*60, true)
	if c.Core.Rand.Float64() < 0.67 {
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Anemo, c.ParticleDelay)
	}
}
