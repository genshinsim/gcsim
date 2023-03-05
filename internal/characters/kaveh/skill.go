package kaveh

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

var skillFrames []int

const (
	skillHitmark   = 32
	particleICDKey = "kaveh-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(46)
	skillFrames[action.ActionDash] = 32
	skillFrames[action.ActionJump] = 32
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Artistic Ingenuity (E)",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4.5) // TODO: correct skill hitbox
	c.Core.QueueAttack(ai, ap, 0, skillHitmark, c.particleCB)
	c.Core.Tasks.Add(func() { c.ruptureDendroCores(ap) }, skillHitmark)

	if c.Base.Cons >= 1 {
		c.c1()
	}

	c.SetCD(action.ActionSkill, 360)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) ruptureDendroCores(ap combat.AttackPattern) {
	for _, g := range c.Core.Combat.Gadgets() {
		seed, ok := g.(*reactable.DendroCore)
		if !ok {
			continue
		}
		if willLand, _ := seed.AttackWillLand(ap); !willLand {
			continue
		}
		seed.Duration = 1
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Dendro, c.ParticleDelay)
}
