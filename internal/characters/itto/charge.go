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

func init() {
	chargeFrames = make([][]int, EndSlashType)

	// ActionCharge frames are different per each slash
	// the frames for CA1/CA2 -> CAF are handled in ActionInfo.Frames

	// CA0 -> x
	chargeFrames[SaichiSlash] = frames.InitAbilSlice(131)
	attackFrames[SaichiSlash][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it
	chargeFrames[SaichiSlash][action.ActionDash] = chargeHitmarks[SaichiSlash]
	chargeFrames[SaichiSlash][action.ActionJump] = chargeHitmarks[SaichiSlash]
	chargeFrames[SaichiSlash][action.ActionSwap] = 130

	// CA1 -> x
	chargeFrames[RightSlash] = frames.InitAbilSlice(104) // NA frames
	chargeFrames[RightSlash][action.ActionCharge] = 57   // CA2 frames
	chargeFrames[RightSlash][action.ActionSkill] = chargeHitmarks[RightSlash]
	chargeFrames[RightSlash][action.ActionBurst] = chargeHitmarks[RightSlash]
	chargeFrames[RightSlash][action.ActionDash] = chargeHitmarks[RightSlash]
	chargeFrames[RightSlash][action.ActionJump] = chargeHitmarks[RightSlash]
	chargeFrames[RightSlash][action.ActionSwap] = chargeHitmarks[RightSlash]

	// CA2 -> x
	chargeFrames[LeftSlash] = frames.InitAbilSlice(77) // NA frames
	chargeFrames[LeftSlash][action.ActionCharge] = 29  // CA1 frames
	chargeFrames[LeftSlash][action.ActionSkill] = chargeHitmarks[LeftSlash]
	chargeFrames[LeftSlash][action.ActionBurst] = chargeHitmarks[LeftSlash]
	chargeFrames[LeftSlash][action.ActionDash] = chargeHitmarks[LeftSlash]
	chargeFrames[LeftSlash][action.ActionJump] = chargeHitmarks[LeftSlash]
	chargeFrames[LeftSlash][action.ActionSwap] = chargeHitmarks[LeftSlash]

	// CAF -> x
	chargeFrames[FinalSlash] = frames.InitAbilSlice(110) // CA0 frames
	chargeFrames[FinalSlash][action.ActionAttack] = 109
	chargeFrames[FinalSlash][action.ActionSkill] = 76
	chargeFrames[FinalSlash][action.ActionBurst] = 76
	chargeFrames[FinalSlash][action.ActionDash] = chargeHitmarks[FinalSlash]
	chargeFrames[FinalSlash][action.ActionJump] = chargeHitmarks[FinalSlash]
	chargeFrames[FinalSlash][action.ActionSwap] = 76
}

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

func (c *char) windupFrames(prevSlash, curSlash SlashType) int {
	switch animState := c.Core.Player.CurrentState(); animState {
	// attack -> x
	case action.NormalAttackState:
		switch curSlash {
		// NA -> CA0
		case SaichiSlash:
			switch c.NormalCounter - 1 {
			case 0:
				return 14
			case 1, 2:
				return 21
			}
		// NA -> CA1/CAF
		case RightSlash, FinalSlash:
			return 10
		}

	// charge -> x
	case action.ChargeAttackState:
		switch curSlash {
		// CAF->CA0
		case SaichiSlash:
			if prevSlash == FinalSlash {
				return 14
			}
		// CA2 -> CA1
		case RightSlash:
			if prevSlash == LeftSlash {
				return 28
			}
		// CA1/CA2 -> CAF
		case FinalSlash:
			if prevSlash == RightSlash || prevSlash == LeftSlash {
				return 25
			}
		}

	// skill -> x
	case action.SkillState:
		switch curSlash {
		// E->CA0
		case SaichiSlash:
			return 14
		// E->CA1/CAF
		case RightSlash, FinalSlash:
			return 17
		}
	}
	return 0
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
		HitlagHaltFrames:   0.10 * 60, // defaults to SaichiSlash/FinalSlash hitlag
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	prevSlash := c.slashState
	stacks := c.Tags[strStackKey]
	c.slashState = prevSlash.Next(stacks)

	// figure out how many frames we need to skip
	windup := c.windupFrames(prevSlash, c.slashState)

	// handle hitlag and talent%
	ai.Abil = fmt.Sprintf("%v (Stacks %v)", c.slashState, stacks)
	switch c.slashState {
	case SaichiSlash:
		ai.Mult = saichiSlash[c.TalentLvlAttack()]
	case RightSlash, LeftSlash:
		ai.Mult = akCombo[c.TalentLvlAttack()]
		haltFrames := 0.03 // stacks >= 4
		switch stacks {
		case 2:
			haltFrames = 0.07
		case 3:
			haltFrames = 0.03
		}
		ai.HitlagHaltFrames = haltFrames * 60
	case FinalSlash:
		ai.Mult = akFinal[c.TalentLvlAttack()]
	}

	// A4: Arataki Kesagiri DMG is increased by 35% of Arataki Itto's DEF.
	if c.slashState != SaichiSlash {
		ai.FlatDmg = 0.35*c.Base.Def*(1+c.Stat(attributes.DEFP)) + c.Stat(attributes.DEF)
	}

	// Unsure of range, it's huge though
	r := 1.0
	if c.StatModIsActive(burstBuffKey) {
		r = 3
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), r, false, combat.TargettableEnemy),
		chargeHitmarks[c.slashState]-windup,
		chargeHitmarks[c.slashState]-windup,
	)

	// C6: has a 50% chance to not consume stacks of Superlative Superstrength.
	if c.Base.Cons < 6 || c.Core.Rand.Float64() < 0.5 {
		c.addStrStack(-1)
	}

	return action.ActionInfo{
		Frames: func(next action.Action) int {
			// handle CA1/CA2 -> CAF frames
			if next == action.ActionCharge && c.slashState.Next(c.Tags[strStackKey]) == FinalSlash {
				switch c.slashState {
				case RightSlash: // CA1 -> CAF
					return 59
				case LeftSlash: // CA2 -> CAF
					return 32
				}
			}
			return chargeFrames[c.slashState][next] - windup
		},
		AnimationLength: chargeFrames[c.slashState][action.InvalidAction] - windup,
		CanQueueAfter:   chargeHitmarks[c.slashState] - windup,
		State:           action.ChargeAttackState,
	}
}
