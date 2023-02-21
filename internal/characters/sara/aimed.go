package sara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames [][]int

var aimedHitmarks = []int{86, 50}

func init() {
	aimedFrames = make([][]int, 2)

	// outside of E status
	aimedFrames[0] = frames.InitAbilSlice(96)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// inside of E status
	aimedFrames[1] = frames.InitAbilSlice(60)
	aimedFrames[1][action.ActionBurst] = 62
	aimedFrames[1][action.ActionDash] = 52
	aimedFrames[1][action.ActionJump] = 52
}

// Aimed charge attack damage queue generator
// Additionally handles crowfeather state, E skill damage, and A4
// Has two parameters, "travel", used to set the number of frames that the arrow is in the air (default = 10)
// weak_point, used to determine if an arrow is hitting a weak point (default = 1 for true)
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	// A1:
	// While in the Crowfeather Cover state provided by Tengu Stormcall, Aimed Shot charge times are decreased by 60%.
	skillActive := 0
	if c.Base.Ascension >= 1 && c.Core.Status.Duration(coverKey) > 0 {
		skillActive = 1
	}

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Aim Charge Attack",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone,
		ICDGroup:             combat.ICDGroupDefault,
		StrikeType:           combat.StrikeTypePierce,
		Element:              attributes.Electro,
		Durability:           25,
		Mult:                 aimChargeFull[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     .12 * 60,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			combat.Point{Y: -0.5},
			0.1,
			1,
		),
		aimedHitmarks[skillActive],
		aimedHitmarks[skillActive]+travel,
	)

	// Cover state handling - drops crowfeather, which explodes after 1.5 seconds
	if c.Core.Status.Duration(coverKey) > 0 {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tengu Juurai: Ambush",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6)

		// TODO: snapshot?
		// Particles are emitted after the ambush thing hits
		c.Core.QueueAttack(ai, ap, aimedHitmarks[skillActive], aimedHitmarks[skillActive]+travel+90, c.makeA4CB(), c.particleCB)
		c.attackBuff(ap, aimedHitmarks[skillActive]+travel+90)

		c.Core.Status.Delete(coverKey)
	}

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return aimedFrames[skillActive][next] },
		AnimationLength: aimedFrames[skillActive][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[skillActive],
		State:           action.AimState,
	}
}
