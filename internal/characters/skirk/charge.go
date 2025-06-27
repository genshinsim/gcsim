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
	chargeHitmarks = []int{27, 27 + 7}
	chargeOffsets  = []float64{1, 1.3}

	chargeSkillFrames   []int
	chargeSkillHitmarks = []int{27, 27 + 7, 27 + 7 + 7}
	chargeSkillOffsets  = []float64{1, 1.3, 1.3}
)

func init() {
	chargeFrames = frames.InitAbilSlice(53) // CA -> W
	chargeFrames[action.ActionAttack] = 43  // CA -> N1
	chargeFrames[action.ActionCharge] = 43  // CA -> CA
	chargeFrames[action.ActionSkill] = 31   // CA -> E
	chargeFrames[action.ActionBurst] = 31   // CA -> Q
	chargeFrames[action.ActionDash] = 27    // CA -> D
	chargeFrames[action.ActionJump] = 28    // CA -> J
	chargeFrames[action.ActionSwap] = 42    // CA -> Swap

	chargeSkillFrames = frames.InitAbilSlice(54) // CA -> N1
	chargeSkillFrames[action.ActionCharge] = 53  // CA -> CA
	chargeSkillFrames[action.ActionSkill] = 31   // CA -> E
	chargeSkillFrames[action.ActionBurst] = 31   // CA -> Q
	chargeSkillFrames[action.ActionDash] = 27    // CA -> D
	chargeSkillFrames[action.ActionJump] = 28    // CA -> J
	chargeSkillFrames[action.ActionWalk] = 52    // CA -> W
	chargeSkillFrames[action.ActionSwap] = 42    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.ChargeAttackSkill(p)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagExtraAttack,
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
		Frames:          func(next action.Action) int { return chargeFrames[next] },
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionDash],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) ChargeAttackSkill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		AttackTag:      attacks.AttackTagExtra,
		ICDTag:         attacks.ICDTagExtraAttack,
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
				geometry.Point{Y: chargeSkillOffsets[i]},
				2.2,
			),
			chargeSkillHitmarks[i],
			chargeSkillHitmarks[i],
			c.absorbVoidRiftCB,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeSkillFrames),
		AnimationLength: chargeSkillFrames[action.InvalidAction],
		CanQueueAfter:   chargeSkillFrames[action.ActionDash], // earliest cancel
		State:           action.ChargeAttackState,
	}, nil
}
