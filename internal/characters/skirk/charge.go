package skirk

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	chargeFrames   []int
	chargeHitmarks = []int{27, 27 + 7}

	chargeSkillFrames   []int
	chargeSkillHitmarks = []int{28, 28 + 9, 28 + 9 + 9}
)

func init() {
	chargeFrames = frames.InitAbilSlice(53) // CA -> W
	chargeFrames[action.ActionAttack] = 43  // CA -> N1
	chargeFrames[action.ActionCharge] = 43  // CA -> CA
	chargeFrames[action.ActionSkill] = 43   // CA -> E
	chargeFrames[action.ActionBurst] = 43   // CA -> Q
	chargeFrames[action.ActionDash] = 27    // CA -> D
	chargeFrames[action.ActionJump] = 28    // CA -> J
	chargeFrames[action.ActionSwap] = 42    // CA -> Swap

	chargeSkillFrames = frames.InitAbilSlice(54) // CA -> N1
	chargeSkillFrames[action.ActionCharge] = 53  // CA -> CA
	chargeSkillFrames[action.ActionBurst] = 43   // CA -> Q
	chargeSkillFrames[action.ActionDash] = 27    // CA -> D
	chargeSkillFrames[action.ActionJump] = 28    // CA -> J
	chargeSkillFrames[action.ActionWalk] = 52    // CA -> W
	chargeSkillFrames[action.ActionSwap] = 42    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.ChargeAttackSkill(p)
	}

	// If the previous attack is N1 to N4 of skill state, this CA is also in skill state
	if c.Core.Player.CurrentState() == action.NormalAttackState && c.Core.Player.ActiveChar().NormalCounter > 0 && c.prevNASkillState {
		return c.ChargeAttackSkill(p)
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
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
		// TODO: y=1 offset when no target in range. What is the range?
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.PrimaryTarget(),
				nil,
				2.8,
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
	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
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
		// TODO: y=1 offset when no target in range. What is the range?
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				nil,
				3.3,
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
