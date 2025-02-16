package chasca

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
	attackFrames   [][]int
	attackHitmarks = [][]int{{10}, {8}, {14, 20}, {39}}

	attackSkillTapFrames []int
)

const normalHitNum = 4
const attackSkillTapHitmark = 11

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 17) // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 19) // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 36) // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 73) // N4 -> N1

	attackFrames[0][action.ActionWalk] = 27
	attackFrames[1][action.ActionWalk] = 29
	attackFrames[2][action.ActionWalk] = 53
	attackFrames[3][action.ActionWalk] = 62

	attackSkillTapFrames = frames.InitAbilSlice(39)
	attackSkillTapFrames[action.ActionAttack] = 34
	attackSkillTapFrames[action.ActionSkill] = 31
	attackSkillTapFrames[action.ActionBurst] = attackSkillTapHitmark
	attackSkillTapFrames[action.ActionDash] = attackSkillTapHitmark
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other bow characters
// Has "travel" parameter, used to set the number of frames that the arrow is in the air (default = 10)
func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.attackSkillTap(p), nil
	}

	windup := 5
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.AimState:
		windup = 0
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	for hitmark := range attackHitmarks[c.NormalCounter] {
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			windup+hitmark,
			windup+hitmark+travel,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          func(next action.Action) int { return frames.NewAttackFunc(c.Character, attackFrames)(next) + windup },
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}, nil
}

func (c *char) attackSkillTap(_ map[string]int) action.Info {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:      attacks.AttackTagNormal,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagChascaTap,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	windup := 7
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		windup = 0
	}

	ap := combat.NewCircleHitFanAngle(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -3.0}, 8.0, 120)
	c.QueueCharTask(func() {
		if !c.nightsoulState.HasBlessing() {
			return
		}
		c.Core.QueueAttack(
			ai,
			ap,
			0,
			0,
		)
	}, windup+attackSkillTapHitmark)

	defer c.AdvanceNormalIndex()
	return action.Info{
		Frames:          c.skillNextFrames(frames.NewAttackFunc(c.Character, attackFrames), 0),
		AnimationLength: attackSkillTapFrames[action.InvalidAction],
		CanQueueAfter:   1, // can run out of nightsoul and start falling earlier
		State:           action.NormalAttackState,
	}
}
