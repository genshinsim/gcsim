package thoma

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillHitmark   = 11
	particleICDKey = "thoma-particle-icd"
)

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
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.06 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	// snapshot unknown
	// snap := c.Snapshot(&ai)

	c.Core.Tasks.Add(func() {
		shieldamt := (shieldpp[c.TalentLvlSkill()]*c.MaxHP() + shieldflat[c.TalentLvlSkill()])
		c.genShield("Thoma Skill", shieldamt, false)
	}, 9)

	// damage component not final
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), combat.Point{Y: 1}, 3, 270),
		skillHitmark,
		skillHitmark,
		c.particleCB,
	)

	player, ok := c.Core.Combat.Player().(*avatar.Player)
	if !ok {
		panic("target 0 should be Player but is not!!")
	}

	c.Core.Tasks.Add(func() {
		player.ApplySelfInfusion(attributes.Pyro, 25, 30)
	}, 9)

	cd := 15
	// TODO: this should only active if a char protected by Thoma's shield is hit, should also proc on stuff like Dori Q self attack
	if c.Base.Cons >= 1 {
		cd = 12
	}
	c.SetCDWithDelay(action.ActionSkill, cd*60, 9)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash],
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
	c.AddStatus(particleICDKey, 0.3*60, true)

	count := 3.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 4
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Pyro, c.ParticleDelay)
}
