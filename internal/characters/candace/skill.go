package candace

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var (
	skillFrames   [][]int
	skillHitmarks = []int{26, 88} // TODO: find correct hitmarks
	skillCDStarts = []int{13, 40} // TODO: find correct CD starts
	skillCD       = []int{360, 540}
)

func init() {
	skillFrames = make([][]int, 2) // TODO: find correct frames
	// Tap E
	skillFrames[0] = frames.InitAbilSlice(26)
	// Hold E
	skillFrames[1] = frames.InitAbilSlice(88)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// Hold parameter gets used in action frames to get earliest possible release frame
	chargeLevel := p["hold"]
	if chargeLevel > 1 {
		chargeLevel = 1
	}
	animIdx := chargeLevel
	if p["perfect"] == 1 {
		animIdx = 0
		chargeLevel = 1
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Sacred Rite: Heron's Sanctum (E)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Hydro,
		Durability:         25,
		FlatDmg:            skillDmg[chargeLevel][c.TalentLvlSkill()] * c.MaxHP(),
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	// Particle should spawn after hit
	hitDelay := skillHitmarks[animIdx]
	switch chargeLevel {
	case 0:
		ai.HitlagHaltFrames = 0.06 * 60 // TODO: check hitlag frames
		c.Core.QueueParticle("candace", 2, attributes.Hydro, c.ParticleDelay+hitDelay)
	case 1:
		c.Core.QueueParticle("candace", 3, attributes.Hydro, c.ParticleDelay+hitDelay)
		ai.Abil = "Sacred Rite: Heron's Sanctum Charged Up (E)"
		ai.HitlagHaltFrames = 0.09 * 60 // TODO: check hitlag frames
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy), hitDelay, hitDelay)

	// Add shield until skill unleashed (treated as frame when attack hits)
	c.Core.Player.Shields.Add(&shield.Tmpl{
		Src:        c.Core.F,
		ShieldType: shield.ShieldCandaceSkill,
		HP:         skillShieldPct[c.TalentLvlSkill()]*c.MaxHP() + skillShieldFlat[c.TalentLvlSkill()],
		Ele:        attributes.Hydro,
		Expires:    c.Core.F + hitDelay,
	})

	c.SetCDWithDelay(action.ActionSkill, skillCD[animIdx], skillCDStarts[animIdx])

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[animIdx]),
		AnimationLength: skillFrames[animIdx][action.InvalidAction],
		CanQueueAfter:   skillFrames[animIdx][action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}
