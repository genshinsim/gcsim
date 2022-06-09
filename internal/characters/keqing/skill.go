package keqing

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 25

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
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), skillHitmark, skillHitmark)

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
		Post:            skillFrames[action.ActionDash], // earliest cancel
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

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), skillRecastHitmark, skillRecastHitmark)

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
			c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), skillRecastHitmark, skillRecastHitmark)
		}
	}

	// TODO: Particle timing?
	if c.Core.Rand.Float64() < .5 {
		c.Core.QueueParticle("keqing", 2, attributes.Electro, 100)
	} else {
		c.Core.QueueParticle("keqing", 3, attributes.Electro, 100)
	}

	// despawn stiletto
	c.Core.Status.Delete(stilettoKey)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionDash], // earliest cancel
		Post:            skillRecastFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
