package yunjin

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var (
	skillFrames   [][]int
	skillHitmarks = []int{13, 50, 93}
	skillCDStarts = []int{11, 48, 90}
)

func init() {
	skillFrames = make([][]int, 3)

	// Tap E
	skillFrames[0] = frames.InitAbilSlice(62) // Tap E -> N1/Q
	skillFrames[0][action.ActionDash] = 49    // Tap E -> D
	skillFrames[0][action.ActionJump] = 48    // Tap E -> J
	skillFrames[0][action.ActionSwap] = 59    // Tap E -> Swap

	// Hold E Lv. 1
	skillFrames[1] = frames.InitAbilSlice(97) // Hold E Lv. 1 -> Q
	skillFrames[1][action.ActionAttack] = 96  // Hold E Lv. 1 -> N1
	skillFrames[1][action.ActionDash] = 85    // Hold E Lv. 1 -> D
	skillFrames[1][action.ActionJump] = 85    // Hold E Lv. 1 -> J
	skillFrames[1][action.ActionSwap] = 95    // Hold E Lv. 1 -> Swap

	// Hold E Lv. 2
	skillFrames[2] = frames.InitAbilSlice(141) // Hold E Lv. 2 -> Q
	skillFrames[2][action.ActionAttack] = 140  // Hold E Lv. 2 -> N1
	skillFrames[2][action.ActionDash] = 129    // Hold E Lv. 2 -> D
	skillFrames[2][action.ActionJump] = 129    // Hold E Lv. 2 -> J
	skillFrames[2][action.ActionSwap] = 138    // Hold E Lv. 2 -> Swap
}

// Skill - modelled after Beidou E
// Has two parameters:
// perfect = 1 if you are doing a perfect counter
// hold = 1 or 2 for regular charging up to level 1 or 2
func (c *char) Skill(p map[string]int) action.ActionInfo {
	// Hold parameter gets used in action frames to get earliest possible release frame
	chargeLevel := p["hold"]
	if chargeLevel > 2 {
		chargeLevel = 2
	}
	animIdx := chargeLevel
	if p["perfect"] == 1 {
		animIdx = 0
		chargeLevel = 2
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Opening Flourish Press (E)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSpear,
		Element:            attributes.Geo,
		Durability:         50,
		Mult:               skillDmg[chargeLevel][c.TalentLvlSkill()],
		UseDef:             true,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	// Particle should spawn after hit
	hitDelay := skillHitmarks[animIdx]
	radius := 4.0
	switch chargeLevel {
	case 0:
		ai.HitlagHaltFrames = 0.06 * 60
		c.Core.QueueParticle("yunjin", 2, attributes.Geo, c.ParticleDelay+hitDelay)
	case 1:
		// 2 or 3, 1:1 ratio
		if c.Core.Rand.Float64() < 0.5 {
			c.Core.QueueParticle("yunjin", 2, attributes.Geo, c.ParticleDelay+hitDelay)
		} else {
			c.Core.QueueParticle("yunjin", 3, attributes.Geo, c.ParticleDelay+hitDelay)
		}
		ai.Abil = "Opening Flourish Level 1 (E)"
		ai.HitlagHaltFrames = 0.09 * 60
		radius = 6
	case 2:
		c.Core.QueueParticle("yunjin", 3, attributes.Geo, c.ParticleDelay+hitDelay)
		ai.Durability = 100
		ai.Abil = "Opening Flourish Level 2 (E)"
		ai.HitlagHaltFrames = 0.12 * 60
		radius = 8
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, radius), hitDelay, hitDelay)

	// Add shield until skill unleashed (treated as frame when attack hits)
	c.Core.Player.Shields.Add(&shield.Tmpl{
		Src:        c.Core.F,
		Name:       "Yun Jin Skill",
		ShieldType: shield.ShieldYunjinSkill,
		HP:         skillShieldPct[c.TalentLvlSkill()]*c.MaxHP() + skillShieldFlat[c.TalentLvlSkill()],
		Ele:        attributes.Geo,
		Expires:    c.Core.F + hitDelay,
	})

	if c.Base.Cons >= 1 {
		// 18% doesn't result in a whole number - 442.8 frames. We round up
		c.SetCDWithDelay(action.ActionSkill, 443, skillCDStarts[animIdx])
	} else {
		c.SetCDWithDelay(action.ActionSkill, 9*60, skillCDStarts[animIdx])
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[animIdx]),
		AnimationLength: skillFrames[animIdx][action.InvalidAction],
		CanQueueAfter:   skillFrames[animIdx][action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}
