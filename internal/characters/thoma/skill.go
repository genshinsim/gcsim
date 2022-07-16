package thoma

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 44

func init() {
	skillFrames = frames.InitAbilSlice(44)
}

// Skill attack damage queue generator
func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Blazing Blessing",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	// snapshot unknown
	// snap := c.Snapshot(&ai)

	//3 or 4, 3:2 ratio
	var count float64 = 3
	if c.Core.Rand.Float64() < 0.4 {
		count = 4
	}
	c.Core.QueueParticle("thoma", count, attributes.Pyro, skillHitmark+c.Core.Flags.ParticleDelay)

	shieldamt := (shieldpp[c.TalentLvlSkill()]*c.MaxHP() + shieldflat[c.TalentLvlSkill()])
	c.genShield("Thoma Skill", shieldamt)

	// damage component not final
	x, y := c.Core.Combat.Target(0).Pos()
	c.Core.QueueAttack(ai, combat.NewCircleHit(x, y, 2, false, combat.TargettableEnemy), skillHitmark, skillHitmark)

	cd := 15
	if c.Base.Cons >= 1 {
		cd = 12 //the CD reduction activates when a character protected by Thoma's shield is hit. Since it is almost impossible for this not to activate, we set the duration to 12 for sim purposes.
	}
	c.SetCD(action.ActionSkill, cd*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.InvalidAction],
		State:           action.SkillState,
	}
}
