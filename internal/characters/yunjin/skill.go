package yunjin

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var skillFrames [][]int
var skillHitmarks = []int{31, 81, 121}

func init() {
	skillFrames = make([][]int, 3)

	// skill (press) -> x
	skillFrames[0] = frames.InitAbilSlice(31)

	// skill (level=1) -> x
	skillFrames[1] = frames.InitAbilSlice(81)

	// skill (level=2) -> x
	skillFrames[2] = frames.InitAbilSlice(121)
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
		ActorIndex: c.Index,
		Abil:       "Opening Flourish Press (E)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       skillDmg[chargeLevel][c.TalentLvlSkill()],
		UseDef:     true,
	}

	// TODO: Fix hit frames when known
	// Particle should spawn after hit
	hitDelay := skillHitmarks[animIdx]
	switch chargeLevel {
	case 0:
		c.Core.QueueParticle("yunjin", 2, attributes.Geo, c.Core.Flags.ParticleDelay+hitDelay)
	case 1:
		// Currently believed to be 2-3 particles with the ratio 3:2
		if c.Core.Rand.Float64() < .6 {
			c.Core.QueueParticle("yunjin", 2, attributes.Geo, c.Core.Flags.ParticleDelay+hitDelay)
		} else {
			c.Core.QueueParticle("yunjin", 3, attributes.Geo, c.Core.Flags.ParticleDelay+hitDelay)
		}
		ai.Abil = "Opening Flourish Level 1 (E)"
	case 2:
		c.Core.QueueParticle("yunjin", 3, attributes.Geo, c.Core.Flags.ParticleDelay+hitDelay)
		ai.Durability = 100
		ai.Abil = "Opening Flourish Level 2 (E)"
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), hitDelay, hitDelay)

	// Add shield until skill unleashed (treated as frame when attack hits)
	c.Core.Player.Shields.Add(&shield.Tmpl{
		Src:        c.Core.F,
		ShieldType: shield.ShieldYunjinSkill,
		HP:         skillShieldPct[c.TalentLvlSkill()]*c.MaxHP() + skillShieldFlat[c.TalentLvlSkill()],
		Ele:        attributes.Geo,
		Expires:    c.Core.F + hitDelay,
	})

	if c.Base.Cons >= 1 {
		// 18% doesn't result in a whole number - 442.8 frames. We round up
		c.SetCD(action.ActionSkill, 443)
	} else {
		c.SetCD(action.ActionSkill, 9*60)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[animIdx]),
		AnimationLength: skillFrames[animIdx][action.InvalidAction],
		CanQueueAfter:   skillHitmarks[animIdx],
		State:           action.SkillState,
	}
}
