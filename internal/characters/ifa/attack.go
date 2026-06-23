package ifa

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
	attackFrames          [][]int
	skillAttackTapFrames  []int
	skillAttackHoldFrames []int
	attackHitmarks        = []int{10, 13, 42}
)

const (
	attackRadius        = 0.7
	normalHitNum        = 3
	attackSkillInterval = 54
)

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 27) // N1 -> Walk
	attackFrames[0][action.ActionAttack] = 15
	attackFrames[0][action.ActionCharge] = 18

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 31) // N2 -> Walk
	attackFrames[1][action.ActionAttack] = 20
	attackFrames[1][action.ActionCharge] = 19

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 90) // N3 -> N1
	attackFrames[2][action.ActionWalk] = 88
	attackFrames[2][action.ActionCharge] = 81

	skillAttackTapFrames = frames.InitNormalCancelSlice(0, attackSkillInterval)
	skillAttackTapFrames[action.ActionCharge] = chargeSkillInterval
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.attackHoldSkillState(p), nil
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.PrimaryTarget(),
		nil,
		attackRadius,
	)

	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, ap, 0, 0)
	}, attackHitmarks[c.NormalCounter])
	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) attackTapSkillState(p map[string]int) action.Info {
	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Tonic Shot",
		AttackTag:      attacks.AttackTagNormal,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagIfaSkill,
		ICDGroup:       attacks.ICDGroupIfaSkillHit,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           skill_dmg[c.TalentLvlSkill()],
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.PrimaryTarget(),
		nil,
		3,
	)

	c.QueueCharTask(func() {
		if !c.nightsoulState.HasBlessing() {
			return
		}
		c.Core.QueueAttack(
			ai,
			ap,
			0,
			0,
			c.particleCB,
			c.healCB,
			c.c1CB,
		)
	}, 3)

	atkspd := c.Stat(attributes.AtkSpd)

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.AtkSpdAdjust(skillAttackTapFrames[next], atkspd)
		},
		AnimationLength: attackSkillInterval,
		CanQueueAfter:   0, // can run out of nightsoul and start falling earlier
		State:           action.NormalAttackState,
	}
}
