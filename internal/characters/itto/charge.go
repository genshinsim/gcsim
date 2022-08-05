package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames [][]int
var chargeHitmarks = []int{89, 51, 24, 71}

type SlashType int

const (
	InvalidSlash SlashType = iota - 1
	SaichiSlash            // CA0
	RightSlash             // CA1
	LeftSlash              // CA2
	FinalSlash             // CAF
	EndSlashType
)

var slashName = []string{
	"Saichimonji Slash",
	"Arataki Kesagiri Right Slash",
	"Arataki Kesagiri Left Slash",
	"Arataki Kesagiri Final Slash",
}

func (s SlashType) String() string {
	return slashName[s]
}

func (s SlashType) Next(stacks int) SlashType {
	switch s {

	// idle -> charge (based on stacks)
	case InvalidSlash:
		if stacks == 1 {
			return FinalSlash
		} else if stacks > 1 {
			return RightSlash
		}
		return SaichiSlash

	// loops CA1/CA2 until stacks=1
	case RightSlash:
		if stacks == 1 {
			return FinalSlash
		}
		return LeftSlash
	case LeftSlash:
		if stacks == 1 {
			return FinalSlash
		}
		return RightSlash

	// CA0/CAF -> x
	case SaichiSlash, FinalSlash:
		fallthrough
	default:
		return SaichiSlash
	}
}

func init() {
	chargeFrames = make([][]int, EndSlashType)

	// CA0 -> x
	chargeFrames[SaichiSlash] = frames.InitAbilSlice(131)
	attackFrames[SaichiSlash][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
	chargeFrames[SaichiSlash][action.ActionDash] = chargeHitmarks[SaichiSlash]
	chargeFrames[SaichiSlash][action.ActionJump] = chargeHitmarks[SaichiSlash]
	chargeFrames[SaichiSlash][action.ActionSwap] = 130

	// CA1 -> x
	chargeFrames[RightSlash] = frames.InitAbilSlice(104)
	chargeFrames[RightSlash][action.ActionCharge] = 57
	chargeFrames[RightSlash][action.ActionSkill] = chargeHitmarks[RightSlash]
	chargeFrames[RightSlash][action.ActionBurst] = chargeHitmarks[RightSlash]
	chargeFrames[RightSlash][action.ActionDash] = chargeHitmarks[RightSlash]
	chargeFrames[RightSlash][action.ActionJump] = chargeHitmarks[RightSlash]
	chargeFrames[RightSlash][action.ActionSwap] = chargeHitmarks[RightSlash]

	// CA2 -> x
	chargeFrames[LeftSlash] = frames.InitAbilSlice(77)
	chargeFrames[LeftSlash][action.ActionCharge] = 29
	chargeFrames[LeftSlash][action.ActionSkill] = chargeHitmarks[LeftSlash]
	chargeFrames[LeftSlash][action.ActionBurst] = chargeHitmarks[LeftSlash]
	chargeFrames[LeftSlash][action.ActionDash] = chargeHitmarks[LeftSlash]
	chargeFrames[LeftSlash][action.ActionJump] = chargeHitmarks[LeftSlash]
	chargeFrames[LeftSlash][action.ActionSwap] = chargeHitmarks[LeftSlash]

	// CAF -> x
	chargeFrames[FinalSlash] = frames.InitAbilSlice(110)
	chargeFrames[FinalSlash][action.ActionAttack] = 109
	chargeFrames[FinalSlash][action.ActionSkill] = 76
	chargeFrames[FinalSlash][action.ActionBurst] = 76
	chargeFrames[FinalSlash][action.ActionDash] = chargeHitmarks[FinalSlash]
	chargeFrames[FinalSlash][action.ActionJump] = chargeHitmarks[FinalSlash]
	chargeFrames[FinalSlash][action.ActionSwap] = 76
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagHaltFrames:   0.10 * 60, // FIXME: adjust haltframes based on stacks/combo
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	prevSlash := c.slashState
	stacks := c.Tags[strStackKey]
	c.slashState = prevSlash.Next(stacks)

	switch c.slashState {
	case SaichiSlash:
		ai.Mult = saichiSlash[c.TalentLvlAttack()]
	case RightSlash, LeftSlash:
		ai.Mult = akCombo[c.TalentLvlAttack()]
		ai.FlatDmg = 0.35*c.Base.Def*(1+c.Stat(attributes.DEFP)) + c.Stat(attributes.DEF)
	case FinalSlash:
		ai.Mult = akFinal[c.TalentLvlAttack()]
		ai.FlatDmg = 0.35*c.Base.Def*(1+c.Stat(attributes.DEFP)) + c.Stat(attributes.DEF)
	}

	// FIXME: hitlag checks goes here. can prob skip CA frames here?
	ai.Abil = fmt.Sprintf("%v (stacks %v)", c.slashState, stacks)

	if c.Base.Cons < 6 || c.Core.Rand.Float64() < 0.5 {
		c.addStrStack(-1)
	}

	r := 1.0
	if c.StatModIsActive(burstBuffKey) {
		// Unsure of range, it's huge though
		r = 3
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), r, false, combat.TargettableEnemy), chargeHitmarks[c.slashState], chargeHitmarks[c.slashState])

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			return chargeFrames[c.slashState][next]
		},
		AnimationLength: chargeFrames[c.slashState][action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[c.slashState],
		State:           action.ChargeAttackState,
	}
}
