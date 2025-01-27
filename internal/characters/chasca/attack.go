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
	attackHitmarks = [][]int{{12}, {8}, {14, 20}, {40}}

	attackSkillTapFrames []int
)

const normalHitNum = 4
const attackSkillTapHitmark = 18

func init() {
	attackFrames = make([][]int, normalHitNum)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 21) // N1 -> N2
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 19) // N2 -> N3
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 36) // N3 -> N4
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 73) // N4 -> N1

	attackSkillTapFrames = frames.InitAbilSlice(40)
	attackSkillTapFrames[action.ActionAim] = 38
	attackSkillTapFrames[action.ActionSkill] = 38
	attackSkillTapFrames[action.ActionBurst] = 20
	attackSkillTapFrames[action.ActionDash] = 20
	attackSkillTapFrames[action.ActionJump] = 33
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other bow characters
// Has "travel" parameter, used to set the number of frames that the arrow is in the air (default = 10)
func (c *char) Attack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.attackSkillTap(p), nil
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
			hitmark,
			hitmark+travel,
		)
	}

	defer c.AdvanceNormalIndex()

	return action.Info{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
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

	ap := combat.NewCircleHitFanAngle(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -3.0}, 8.0, 120)

	c.Core.QueueAttack(
		ai,
		ap,
		attackSkillTapHitmark,
		attackSkillTapHitmark,
	)

	defer c.AdvanceNormalIndex()
	return action.Info{
		Frames:          c.skillNextFrames(frames.NewAttackFunc(c.Character, attackFrames)),
		AnimationLength: attackSkillTapFrames[action.InvalidAction],
		CanQueueAfter:   attackSkillTapHitmark,
		State:           action.NormalAttackState,
	}
}
