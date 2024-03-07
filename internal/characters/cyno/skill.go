package cyno

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
	skillBName     = "Mortuary Rite"
	particleICDKey = "cyno-particle-icd"
)

var (
	skillCD       = 450
	skillBCD      = 180
	skillCDDelay  = 17
	skillBCDDelay = 26
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

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(burstKey) {
		return c.skillB()
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Secret Rite: Chasmic Soulfarer",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
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

	c.Core.Tasks.Add(c.triggerSkillCD, skillCDDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillB() (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             skillBName,
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeBlunt,
		PoiseDMG:         75,
		Element:          attributes.Electro,
		Durability:       25,
		Mult:             skillB[c.TalentLvlSkill()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.03 * 60,
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.5},
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
		c.c6Init()

		c.Core.QueueAttack(ai, ap, skillBHitmark, skillBHitmark, particleCB)
		// Apply the extra hit
		ai.Abil = "Duststalker Bolt"
		ai.Mult = 1.0
		ai.FlatDmg = c.a4Bolt()
		ai.AttackTag = attacks.AttackTagElementalArtHold
		ai.ICDTag = attacks.ICDTagElementalArt
		ai.ICDGroup = attacks.ICDGroupCynoBolt
		ai.StrikeType = attacks.StrikeTypeSlash
		ai.PoiseDMG = 25
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
		c.ExtendStatus(burstKey, 240) // 4s*60
		c.burstExtension++
	}

	c.tryBurstPPSlide(skillBHitmark)

	c.Core.Tasks.Add(c.triggerSkillCD, skillBCDDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillBFrames),
		AnimationLength: skillBFrames[action.InvalidAction],
		CanQueueAfter:   skillBFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) triggerSkillCD() {
	c.ResetActionCooldown(action.ActionSkill)
	c.SetCD(action.ActionSkill, skillCD)
	c.SetCD(action.ActionLowPlunge, skillBCD)
}

func (c *char) makeParticleCB(burst bool) combat.AttackCBFunc {
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(particleICDKey) {
			return
		}
		c.AddStatus(particleICDKey, 0.5*60, true)

		var count float64
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
