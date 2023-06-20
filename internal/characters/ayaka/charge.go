package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int
var chargeHitmarks = []int{27, 33, 39}

func init() {
	chargeFrames = frames.InitAbilSlice(71)
	chargeFrames[action.ActionSkill] = 62
	chargeFrames[action.ActionBurst] = 63
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = chargeHitmarks[len(chargeHitmarks)-1]
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Charge",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
		ICDGroup:   attacks.ICDGroupAyakaExtraAttack,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       ca[c.TalentLvlAttack()],
	}

	// spawn up to 5 attacks
	// priority: enemy > gadget
	chargeCount := 5
	chargeArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4)
	charge := func(pos geometry.Point, i int) {
		for j := 0; j < 3; j++ {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(
					pos,
					nil,
					1,
				),
				chargeHitmarks[j],
				chargeHitmarks[j],
				c.c1,
				c.c6,
			)
		}
	}

	// target up to 5 enemies
	enemies := c.Core.Combat.EnemiesWithinArea(chargeArea, nil)
	enemyCount := len(enemies)
	for i := 0; i < chargeCount; i++ {
		if i < enemyCount {
			charge(enemies[i].Pos(), i)
		}
	}
	chargeCount -= enemyCount

	// if less than 5 enemies were targeted, then check for gadgets
	if chargeCount > 0 {
		gadgets := c.Core.Combat.GadgetsWithinArea(chargeArea, nil)
		gadgetCount := len(gadgets)
		for i := 0; i < chargeCount; i++ {
			if i < gadgetCount {
				charge(gadgets[i].Pos(), i)
			}
		}
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
