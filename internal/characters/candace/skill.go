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
	skillHitmarks = []int{26, 91}
	skillCDStarts = []int{14, 89}
	skillCD       = []int{360, 540}
)

func init() {
	skillFrames = make([][]int, 2)
	// Tap E
	skillFrames[0] = frames.InitAbilSlice(26)
	skillFrames[0][action.ActionBurst] = 25
	skillFrames[0][action.ActionSwap] = 25
	// Hold E
	skillFrames[1] = frames.InitAbilSlice(113)
	skillFrames[0][action.ActionAttack] = 112
	skillFrames[0][action.ActionBurst] = 112
	skillFrames[0][action.ActionJump] = 111
	skillFrames[0][action.ActionSwap] = 110
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	chargeLevel := p["hold"]
	if chargeLevel > 1 {
		chargeLevel = 1
	}
	perfect := p["perfect"] != 0
	if perfect {
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
		HitlagHaltFrames:   0.05 * 60,
		CanBeDefenseHalted: true,
	}

	// Particle should spawn after hit
	hitDelay := skillHitmarks[chargeLevel]
	if perfect {
		hitDelay -= 55
	}

	switch chargeLevel {
	case 0:
		c.Core.QueueParticle("candace", 2, attributes.Hydro, c.ParticleDelay+hitDelay)
	case 1:
		c.Core.QueueParticle("candace", 3, attributes.Hydro, c.ParticleDelay+hitDelay)
		ai.Abil = "Sacred Rite: Heron's Sanctum Charged Up (E)"
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
		hitDelay,
		hitDelay,
		func(_ combat.AttackCB) {
			if c.Base.Cons >= 2 {
				c.c2()
			}
		},
	)

	// Add shield until skill unleashed (treated as frame when attack hits)
	c.Core.Player.Shields.Add(&shield.Tmpl{
		Src:        c.Core.F,
		ShieldType: shield.ShieldCandaceSkill,
		HP:         skillShieldPct[c.TalentLvlSkill()]*c.MaxHP() + skillShieldFlat[c.TalentLvlSkill()],
		Ele:        attributes.Hydro,
		Expires:    c.Core.F + hitDelay,
	})

	cd := skillCD[chargeLevel]
	if c.Base.Cons >= 4 {
		cd = skillCD[0]
	}
	c.SetCDWithDelay(action.ActionSkill, cd, skillCDStarts[chargeLevel])

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			f := skillFrames[chargeLevel][next]
			if perfect {
				f -= 55
			}
			return f
		},
		AnimationLength: skillFrames[chargeLevel][action.InvalidAction],
		CanQueueAfter:   skillFrames[chargeLevel][action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}
