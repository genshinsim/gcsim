package arlecchino

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var chargeFrames []int

const chargeHitmark = 37

func init() {
	chargeFrames = frames.InitAbilSlice(60)
	chargeFrames[action.ActionAttack] = 42
	chargeFrames[action.ActionCharge] = 53
	chargeFrames[action.ActionSkill] = 42
	chargeFrames[action.ActionBurst] = 42
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 58
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.swapError {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Charged Attack with Swap", c.CharWrapper.Base.Key)
	}

	if c.chargeEarlyCancelled {
		return action.Info{}, fmt.Errorf("%v: Cannot early cancel Charged Attack with Charged Attack", c.CharWrapper.Base.Key)
	}

	windup := 0
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		windup = 12
	}

	early, ok := p["early_cancel"]
	if !ok {
		early = 0
	}

	c.QueueCharTask(func() {
		c.absorbDirectives()
	}, 12-windup)

	if early > 0 {
		c.chargeEarlyCancelled = true
		// TODO: error if the user waits until after hitmark to do the dash/jump
		return action.Info{
			Frames:          func(next action.Action) int { return 13 - windup },
			AnimationLength: chargeHitmark - 1,
			CanQueueAfter:   13 - windup,
			State:           action.ChargeAttackState,
			OnRemoved: func(next action.AnimationState) {
				if next == action.SwapState {
					c.swapError = true
				}
			},
		}, nil
	}

	c.QueueCharTask(func() {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Charge",
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagExtraAttack,
			ICDGroup:           attacks.ICDGroupPoleExtraAttack,
			StrikeType:         attacks.StrikeTypeSpear,
			Element:            attributes.Physical,
			Durability:         25,
			HitlagHaltFrames:   0.02,
			HitlagFactor:       0.01,
			CanBeDefenseHalted: true,
			Mult:               charge[c.TalentLvlAttack()],
		}

		if c.StatusIsActive(naBuffKey) {
			ai.Element = attributes.Pyro
			ai.IgnoreInfusion = true
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: 0.9},
				4,
			),
			0,
			0,
		)
	}, chargeHitmark-windup)
	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeHitmark - windup,
		State:           action.ChargeAttackState,
	}, nil
}
