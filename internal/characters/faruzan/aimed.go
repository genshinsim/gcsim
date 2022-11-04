package faruzan

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
// A4 effect is: When Tengu Juurai: Ambush hits opponents, Kujou faruzan will restore 1.2 Energy to all party members for every 100% Energy Recharge she has. This effect can be triggered once every 3s.
// Has two parameters, "travel", used to set the number of frames that the arrow is in the air (default = 10)
// weak_point, used to determine if an arrow is hitting a weak point (default = 1 for true)
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	skillActive := 0
	if c.StatusIsActive(skillKey) && c.hurricaneCount > 0 {
		skillActive = 1
	}

	if skillActive == 0 {
		ai := combat.AttackInfo{
			ActorIndex:           c.Index,
			Abil:                 "Aim Charge Attack",
			AttackTag:            combat.AttackTagExtra,
			ICDTag:               combat.ICDTagNone,
			ICDGroup:             combat.ICDGroupDefault,
			StrikeType:           combat.StrikeTypePierce,
			Element:              attributes.Anemo,
			Durability:           25,
			Mult:                 aimChargeFull[c.TalentLvlAttack()],
			HitWeakPoint:         weakspot == 1,
			HitlagHaltFrames:     .12 * 60,
			HitlagOnHeadshotOnly: true,
			IsDeployable:         true,
		}
		c.Core.QueueAttack(
			ai,
			combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget),
			aimedHitmarks[skillActive],
			aimedHitmarks[skillActive]+travel,
		)
	} else {
		c.Core.Tasks.Add(func() {
			c.hurricaneArrow(travel, weakspot == 1)
		}, aimedHitmarks[skillActive])
		c.hurricaneCount--
	}

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return aimedFrames[skillActive][next] },
		AnimationLength: aimedFrames[skillActive][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[skillActive],
		State:           action.AimState,
	}
}
