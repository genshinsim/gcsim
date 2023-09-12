package keqing

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
var skillRecastFrames []int

const (
	skillHitmark       = 25
	skillRecastHitmark = 16
	stilettoKey        = "keqingstiletto"
	particleICDKey     = "keqing-particle-icd"
)

func init() {
	// skill -> x
	skillFrames = frames.InitAbilSlice(37)
	skillFrames[action.ActionAttack] = 36
	skillFrames[action.ActionSkill] = 35
	skillFrames[action.ActionDash] = 21
	skillFrames[action.ActionJump] = 21
	skillFrames[action.ActionSwap] = 28

	// skill (recast) -> x
	skillRecastFrames = frames.InitAbilSlice(43)
	skillRecastFrames[action.ActionAttack] = 42
	skillRecastFrames[action.ActionDash] = skillRecastHitmark
	skillRecastFrames[action.ActionJump] = skillRecastHitmark
	skillRecastFrames[action.ActionSwap] = 42
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// check if stiletto is on-field
	if c.Core.Status.Duration(stilettoKey) > 0 {
		return c.skillRecast(p)
	}
	return c.skillFirst(p)
}

func (c *char) skillFirst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Stellar Restoration",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 1.6),
		skillHitmark,
		skillHitmark,
	)

	if c.Base.Cons >= 6 {
		c.c6("skill")
	}

	// spawn after cd and stays for 5s
	c.Core.Status.Add(stilettoKey, 5*60+20)

	c.SetCDWithDelay(action.ActionSkill, 7*60+30, 20)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillRecast(p map[string]int) action.ActionInfo {
	// C1 DMG happens before Recast DMG
	if c.Base.Cons >= 1 {
		ai := combat.AttackInfo{
			Abil:       "Stellar Restoration (C1)",
			ActorIndex: c.Index,
			AttackTag:  attacks.AttackTagElementalArtHold,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       .5,
		}
		// 2 dmg instances at start and end
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2),
			3,
			3,
		)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 1.5}, 2),
			skillRecastHitmark,
			skillRecastHitmark,
		)
	}

	ai := combat.AttackInfo{
		Abil:             "Stellar Restoration (Slashing)",
		ActorIndex:       c.Index,
		AttackTag:        attacks.AttackTagElementalArt,
		ICDTag:           attacks.ICDTagElementalArt,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeSlash,
		Element:          attributes.Electro,
		Durability:       50,
		Mult:             skillPress[c.TalentLvlSkill()],
		HitlagHaltFrames: 0.09 * 60,
		HitlagFactor:     0.01,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 3),
		skillRecastHitmark,
		skillRecastHitmark,
		c.particleCB,
	)

	// add electro infusion
	c.a1()

	// despawn stiletto
	c.Core.Status.Delete(stilettoKey)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionDash], // earliest cancel
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
	c.AddStatus(particleICDKey, 0.6*60, true)

	count := 2.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 3
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Electro, c.ParticleDelay)
}
