package ningguang

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	chargeFrames       [][]int
	chargeHitmarks     []int
	chargeJadeHitmarks []int
	chargeC6Hitmarks   []int
	chargeOptions      = map[attackType][]attackType{
		attackTypeLeft:  {attackTypeRight, attackTypeTwirl},
		attackTypeRight: {attackTypeLeft, attackTypeTwirl},
		attackTypeTwirl: {attackTypeLeft, attackTypeRight},
	}
)

func init() {
	chargeHitmarks = make([]int, endAttackType)
	chargeHitmarks[attackTypeLeft] = 35
	chargeHitmarks[attackTypeRight] = 58
	chargeHitmarks[attackTypeTwirl] = 66

	chargeJadeHitmarks = make([]int, endAttackType)
	chargeJadeHitmarks[attackTypeLeft] = 45
	chargeJadeHitmarks[attackTypeRight] = 66
	chargeJadeHitmarks[attackTypeTwirl] = 74

	chargeC6Hitmarks = make([]int, endAttackType)
	chargeC6Hitmarks[attackTypeLeft] = 40
	chargeC6Hitmarks[attackTypeRight] = 58
	chargeC6Hitmarks[attackTypeTwirl] = 67

	chargeFrames = make([][]int, endAttackType)
	// CA Left > x
	chargeFrames[attackTypeLeft] = frames.InitAbilSlice(52)
	chargeFrames[attackTypeLeft][action.ActionAttack] = 46
	chargeFrames[attackTypeLeft][action.ActionCharge] = 47
	chargeFrames[attackTypeLeft][action.ActionSkill] = 48
	chargeFrames[attackTypeLeft][action.ActionBurst] = 47
	chargeFrames[attackTypeLeft][action.ActionDash] = 32
	chargeFrames[attackTypeLeft][action.ActionJump] = 32
	chargeFrames[attackTypeLeft][action.ActionSwap] = 46
	// CA Right > x
	chargeFrames[attackTypeRight] = frames.InitAbilSlice(74)
	chargeFrames[attackTypeRight][action.ActionAttack] = 68
	chargeFrames[attackTypeRight][action.ActionCharge] = 69
	chargeFrames[attackTypeRight][action.ActionSkill] = 72
	chargeFrames[attackTypeRight][action.ActionBurst] = 69
	chargeFrames[attackTypeRight][action.ActionDash] = 54
	chargeFrames[attackTypeRight][action.ActionJump] = 54
	chargeFrames[attackTypeRight][action.ActionSwap] = 67
	// CA Twirl > x
	chargeFrames[attackTypeTwirl] = frames.InitAbilSlice(82)
	chargeFrames[attackTypeTwirl][action.ActionAttack] = 76
	chargeFrames[attackTypeTwirl][action.ActionCharge] = 76
	chargeFrames[attackTypeTwirl][action.ActionSkill] = 77
	chargeFrames[attackTypeTwirl][action.ActionBurst] = 77
	chargeFrames[attackTypeTwirl][action.ActionDash] = 61
	chargeFrames[attackTypeTwirl][action.ActionJump] = 62
	chargeFrames[attackTypeTwirl][action.ActionSwap] = 76
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	chargeType := chargeOptions[c.prevAttack][c.Core.Rand.Intn(2)]
	if c.Core.Player.CurrentState() == action.NormalAttackState {
		switch c.prevAttack {
		case attackTypeLeft,
			attackTypeTwirl:
			chargeType = attackTypeLeft
		case attackTypeRight:
			chargeType = attackTypeRight
		}
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Charge (%s)", chargeType),
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagExtraAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	// skip CA windup if we're in NA/CA animation
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.ChargeAttackState:
		windup = 15
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			1.5,
		),
		chargeHitmarks[chargeType]-windup,
		chargeHitmarks[chargeType]-windup+travel,
	)

	ai = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Charge Gem (%s)", chargeType),
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagExtraAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Geo,
		Durability:         50,
		Mult:               jade[c.TalentLvlAttack()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 0.5)

	jadeHitmarks := chargeJadeHitmarks
	if c.jadeCount == 7 {
		jadeHitmarks = chargeC6Hitmarks
	}
	for i := 0; i < c.jadeCount; i++ {
		c.Core.QueueAttack(
			ai,
			ap,
			jadeHitmarks[chargeType]-windup,
			jadeHitmarks[chargeType]-windup+travel,
		)
	}
	c.jadeCount = 0

	canQueueAfter := math.MaxInt32
	for _, f := range chargeFrames[chargeType] {
		if f < canQueueAfter {
			canQueueAfter = f
		}
	}

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return chargeFrames[chargeType][next] - windup
		},
		AnimationLength: chargeFrames[chargeType][action.InvalidAction] - windup,
		CanQueueAfter:   canQueueAfter - windup,
		State:           action.ChargeAttackState,
	}
}
