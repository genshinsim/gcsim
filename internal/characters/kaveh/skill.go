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
	skillHitmark   = 35
	particleICDKey = "kaveh-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(52)
	skillFrames[action.ActionAttack] = 45
	skillFrames[action.ActionSkill] = 45
	skillFrames[action.ActionBurst] = 45
	skillFrames[action.ActionDash] = 34
	skillFrames[action.ActionJump] = 34
	skillFrames[action.ActionSwap] = 44
}

func (c *char) Skill(p map[string]int) action.Info {
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

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	c.Core.QueueAttack(ai, ap, 0, skillHitmark, c.particleCB)
	c.Core.Tasks.Add(func() { c.ruptureDendroCores(ap) }, skillHitmark+3)

	if c.Base.Cons >= 1 {
		c.Core.Tasks.Add(c.c1, 33)
	}

	c.SetCDWithDelay(action.ActionSkill, 360, 33)

	return action.Info{
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
	c.AddStatus(particleICDKey, 0.2*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Dendro, c.ParticleDelay)
}
