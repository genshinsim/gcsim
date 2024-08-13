package emilie

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	lumidouceArkheCD = "lumidouce-arkhe-cd"
	particleICDKey   = "skill-particle-icd"

	lumidouceSummonHitmark = 16
	lumidouceArkheHitmark  = 16
	particleICD            = 2.5 * 60
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(48)
	skillFrames[action.ActionAttack] = 31
	skillFrames[action.ActionBurst] = 31
	skillFrames[action.ActionDash] = 29
	skillFrames[action.ActionJump] = 30
	skillFrames[action.ActionSwap] = 29
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lumidouce Case (Summon)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillDMG[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 2.6}, 4.5),
		lumidouceSummonHitmark,
		lumidouceSummonHitmark,
	)

	c.spawnLumidouceCase(1)
	c.arkheAttack()
	c.SetCD(action.ActionSkill, int(skillCD[c.TalentLvlSkill()]*60))

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) arkheAttack() {
	if c.StatusIsActive(lumidouceArkheCD) {
		return
	}
	c.AddStatus(lumidouceArkheCD, int(skillArkeCD[c.TalentLvlBurst()]*60), true)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spiritbreath Thorn",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillArkeCD[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 2.6}, 4.5),
		lumidouceArkheHitmark,
		lumidouceArkheHitmark,
	)
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Dendro, c.ParticleDelay)
}
