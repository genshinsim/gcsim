package skirk

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	chargeFrames   []int
	chargeHitmarks = []int{7, 11, 17} // CA-1 and CA-2 hit at the same time
	chargeOffsets  = []float64{1, 1.3, 1.3}
)

func init() {
	chargeFrames = frames.InitAbilSlice(46) // CA -> N1
	chargeFrames[action.ActionDash] = 25    // CA -> D
	chargeFrames[action.ActionJump] = 24    // CA -> J
	chargeFrames[action.ActionSwap] = 30    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.ChargeAttackSkill(p)
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: chargeOffsets[i]},
				2.2,
			),
			chargeHitmarks[i],
			chargeHitmarks[i],
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionJump], // earliest cancel
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) ChargeAttackSkill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		AttackTag:      attacks.AttackTagExtra,
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeSlash,
		Element:        attributes.Cryo,
		Durability:     25,
		IgnoreInfusion: true,
	}

	for i, mult := range skillCharge {
		ai.Mult = mult[c.TalentLvlSkill()]
		ai.Abil = fmt.Sprintf("Charge (Skill) %v", i)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				geometry.Point{Y: chargeOffsets[i]},
				2.2,
			),
			chargeHitmarks[i],
			chargeHitmarks[i],
			c.absorbVoidRiftCB,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionJump], // earliest cancel
		State:           action.ChargeAttackState,
	}, nil
}
