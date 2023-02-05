package cyno

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const (
	skillBName     = "Mortuary Rite"
	particleICDKey = "cyno-particle-icd"
)

var (
	skillCD       = 450
	skillHitmark  = 21
	skillBHitmark = 28
	skillFrames   []int
	skillBFrames  []int
)

func init() {
	skillFrames = frames.InitAbilSlice(43)
	skillFrames[action.ActionDash] = 31
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 42

	// burst frames
	skillBFrames = frames.InitAbilSlice(34)
	skillBFrames[action.ActionDash] = 30
	skillBFrames[action.ActionJump] = 31
	skillBFrames[action.ActionSwap] = 33
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if c.StatusIsActive(BurstKey) {
		return c.skillB()
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Secret Rite: Chasmic Soulfarer",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			1,
		),
		skillHitmark,
		skillHitmark,
		c.makeParticleCB(false),
	)

	c.lastSkillCast = c.Core.F + 17
	c.SetCDWithDelay(action.ActionSkill, skillCD, 17)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillB() action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             skillBName,
		AttackTag:        combat.AttackTagElementalArt,
		ICDTag:           combat.ICDTagNone,
		ICDGroup:         combat.ICDGroupDefault,
		StrikeType:       combat.StrikeTypeBlunt,
		Element:          attributes.Electro,
		Durability:       25,
		Mult:             skillB[c.TalentLvlSkill()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.03 * 60,
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		combat.Point{Y: 1.5},
		6,
	)
	particleCB := c.makeParticleCB(true)
	if !c.StatusIsActive(a1Key) { // check for endseer buff
		c.Core.QueueAttack(ai, ap, skillBHitmark, skillBHitmark, particleCB)
	} else {
		// apply the extra damage on skill
		c.a1Buff()
		if c.Base.Cons >= 1 && c.StatusIsActive(c1Key) {
			c.c1()
		}
		if c.Base.Cons >= 6 { // constellation 6 giving 4 stacks on judication
			c.c6Stacks += 4
			c.AddStatus(c6Key, 480, true) // 8s*60
			if c.c6Stacks > 8 {
				c.c6Stacks = 8
			}
		}

		c.Core.QueueAttack(ai, ap, skillBHitmark, skillBHitmark, particleCB)
		// Apply the extra hit
		ai.Abil = "Duststalker Bolt"
		ai.Mult = 1.0
		ai.FlatDmg = c.a4Bolt()
		ai.ICDTag = combat.ICDTagCynoBolt
		ai.ICDGroup = combat.ICDGroupCynoBolt
		ai.StrikeType = combat.StrikeTypeSlash
		ai.HitlagFactor = 0
		ai.HitlagHaltFrames = 0

		// 3 instances
		for i := 0; i < 3; i++ {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					0.3,
				),
				skillBHitmark,
				skillBHitmark,
				particleCB,
			)
		}

	}
	if c.burstExtension < 2 { // burst can only be extended 2 times per burst cycle (up to 18s, 10s base and +4 each time)
		c.ExtendStatus(BurstKey, 240) // 4s*60
		c.burstExtension++
	}

	c.tryBurstPPSlide(skillBHitmark)

	c.lastSkillCast = c.Core.F + 26
	c.SetCDWithDelay(action.ActionSkill, 180, 26)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillBFrames),
		AnimationLength: skillBFrames[action.InvalidAction],
		CanQueueAfter:   skillBFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) makeParticleCB(burst bool) combat.AttackCBFunc {
	var count float64
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.StatusIsActive(particleICDKey) {
			return
		}
		c.AddStatus(particleICDKey, 0.5*60, true)

		if burst {
			count = 1
			if c.Core.Rand.Float64() < 0.33 {
				count = 2
			}
		} else {
			count = 3
		}
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Electro, c.ParticleDelay)
	}
}
