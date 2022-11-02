package keqing

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 25
const stilettoKey = "keqingstiletto"

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
	skillRecastFrames[action.ActionDash] = 15
	skillRecastFrames[action.ActionJump] = 16
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
		Abil:               "Stellar Restoration",
		ActorIndex:         c.Index,
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.09 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1), skillHitmark, skillHitmark)

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

var skillRecastFrames []int

const skillRecastHitmark = 27

func (c *char) skillRecast(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Stellar Restoration (Slashing)",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 50,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1), skillRecastHitmark, skillRecastHitmark)

	//add electro infusion
	c.a1()

	if c.Base.Cons >= 1 {
		//2 tick dmg at start to end
		hits, ok := p["c1"]
		if !ok {
			hits = 1 //default 1 hit
		}
		ai := combat.AttackInfo{
			Abil:       "Stellar Restoration (C1)",
			ActorIndex: c.Index,
			AttackTag:  combat.AttackTagElementalArtHold,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       .5,
		}
		// TODO: this should be 1st hit on cast and 2nd at end
		for i := 0; i < hits; i++ {
			c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2), skillRecastHitmark, skillRecastHitmark)
		}
	}

	// TODO: Particle timing?
	count := 2.0
	if c.Core.Rand.Float64() < .5 {
		count = 3
	}
	c.Core.QueueParticle("keqing", count, attributes.Electro, skillRecastHitmark+c.ParticleDelay)

	// despawn stiletto
	c.Core.Status.Delete(stilettoKey)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
