package thoma

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 11

func init() {
	skillFrames = frames.InitAbilSlice(46)
	skillFrames[action.ActionDash] = 32
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 44
}

// Skill attack damage queue generator
func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Blazing Blessing",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.06 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	// snapshot unknown
	// snap := c.Snapshot(&ai)

	// 3 or 4, 1:1 ratio
	var count float64 = 3
	if c.Core.Rand.Float64() < 0.5 {
		count = 4
	}
	c.Core.QueueParticle("thoma", count, attributes.Pyro, skillHitmark+c.ParticleDelay)

	c.Core.Tasks.Add(func() {
		shieldamt := (shieldpp[c.TalentLvlSkill()]*c.MaxHP() + shieldflat[c.TalentLvlSkill()])
		c.genShield("Thoma Skill", shieldamt, false)
	}, 9)

	// damage component not final
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 2),
		skillHitmark,
		skillHitmark,
	)

	player, ok := c.Core.Combat.Player().(*avatar.Player)
	if !ok {
		panic("target 0 should be Player but is not!!")
	}

	c.Core.Tasks.Add(func() {
		player.ApplySelfInfusion(attributes.Pyro, 25, 30)
	}, 9)

	cd := 15
	if c.Base.Cons >= 1 {
		cd = 12 // the CD reduction activates when a character protected by Thoma's shield is hit. Since it is almost impossible for this not to activate, we set the duration to 12 for sim purposes.
	}
	c.SetCDWithDelay(action.ActionSkill, cd*60, 9)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash],
		State:           action.SkillState,
	}
}
