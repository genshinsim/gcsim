package faruzan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
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

func (c *char) Aimed(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	skillActive := 0
	if c.StatusIsActive(skillKey) && c.hurricaneCount > 0 {
		// A1:
		// When Faruzan is in the Manifest Gale state created by Wind Realm of Nasamjnin,
		// the amount of time taken to charge a shot is decreased by 60%.
		if c.Base.Ascension >= 1 {
			skillActive = 1
		}
		c.hurricaneCount -= 1
		if c.hurricaneCount <= 0 {
			c.DeleteStatus(skillKey)
		}
	}

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Aim Charge Attack",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
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
			c.pressurizedCollapse(a.Target.Pos())
		}
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -0.5},
			0.1,
			1,
		),
		aimedHitmarks[skillActive],
		aimedHitmarks[skillActive]+travel,
		skillCb,
	)

	return action.Info{
		Frames:          func(next action.Action) int { return aimedFrames[skillActive][next] },
		AnimationLength: aimedFrames[skillActive][action.InvalidAction],
		CanQueueAfter:   aimedFrames[skillActive][action.ActionDash],
		State:           action.AimState,
	}, nil
}
