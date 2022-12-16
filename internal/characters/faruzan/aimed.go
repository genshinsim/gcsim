package faruzan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	aimedFrames   [][]int
	aimedHitmarks = []int{86, 49}
)

func init() {
	aimedFrames = make([][]int, 2)

	// outside of E status
	aimedFrames[0] = frames.InitAbilSlice(96)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// inside of E status
	aimedFrames[1] = frames.InitAbilSlice(60)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	skillActive := 0
	if c.StatusIsActive(skillKey) && c.hurricaneCount > 0 {
		skillActive = 1
		c.hurricaneCount -= 1
		if c.hurricaneCount <= 0 {
			c.DeleteStatus(skillKey)
		}
	}

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

	var skillCb func(a combat.AttackCB)

	if skillActive == 1 {
		ai.Abil = "Hurricane Arrow"
		done := false
		skillCb = func(a combat.AttackCB) {
			if done {
				return
			}
			c.pressurizedCollapse(a.Target)
		}
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 0.5),
		aimedHitmarks[skillActive],
		aimedHitmarks[skillActive]+travel,
		skillCb,
	)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return aimedFrames[skillActive][next] },
		AnimationLength: aimedFrames[skillActive][action.InvalidAction],
		CanQueueAfter:   aimedFrames[skillActive][action.ActionDash],
		State:           action.AimState,
	}
}
