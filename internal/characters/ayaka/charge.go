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
	checkDelay := chargeHitmarks[0] - 1 // TODO: exact delay unknown
	singleCharge := func(pos geometry.Point, hitmark int) {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				pos,
				nil,
				1,
			),
			hitmark,
			hitmark,
			c.c1,
			c.c6,
		)
	}

	charge := func(target combat.Target) {
		for j := 0; j < 3; j++ {
			// queue up ca hits because target could move
			c.Core.Tasks.Add(func() {
				singleCharge(target.Pos(), 0)
			}, chargeHitmarks[j]-checkDelay)
		}
	}

	c.Core.Tasks.Add(func() {
		// look for enemies around the player
		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), nil)

		// don't do anything if there are no enemies in range
		if enemies == nil {
			return
		}

		// check for enemies around the enemy found
		anchorEnemy := enemies[0]
		chargeArea := combat.NewCircleHitOnTarget(anchorEnemy, nil, 4)
		enemies = c.Core.Combat.EnemiesWithinArea(chargeArea, func(t combat.Enemy) bool {
			return t.Key() != anchorEnemy.Key() // don't want to target the same enemy twice
		})
		enemyCount := len(enemies)

		// spawn attacks on enemies
		charge(anchorEnemy)
		chargeCount -= 1
		for i := 0; i < chargeCount; i++ {
			if i < enemyCount {
				charge(enemies[i])
			}
		}
		chargeCount -= enemyCount

		// queue up following ca hits

		// if less than 5 enemies were targeted, then check for gadgets
		if chargeCount > 0 {
			gadgets := c.Core.Combat.GadgetsWithinArea(chargeArea, nil)
			gadgetCount := len(gadgets)
			for i := 0; i < chargeCount; i++ {
				if i < gadgetCount {
					charge(gadgets[i])
				}
			}
		}
	}, checkDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
