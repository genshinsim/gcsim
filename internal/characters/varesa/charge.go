package varesa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

// TODO: update hitlags/hitboxes
var (
	chargeFrames      []int
	fieryChargeFrames []int
)

const (
	chargeHitmark      = 69
	fieryChargeHitmark = 69

	fastChargeHitmark      = 11
	fastFieryChargeHitmark = 11

	fastChargedFrame      = 57
	fastFieryChargedFrame = 58
)

func init() {
	chargeFrames = frames.InitAbilSlice(143) // CA -> Q/Dash/Jump
	chargeFrames[action.ActionAttack] = 142
	chargeFrames[action.ActionCharge] = 141
	chargeFrames[action.ActionSkill] = 142
	chargeFrames[action.ActionWalk] = 142
	chargeFrames[action.ActionSwap] = 139
	chargeFrames[action.ActionHighPlunge] = 77

	fieryChargeFrames = frames.InitAbilSlice(143) // CA -> CA/Jump
	fieryChargeFrames[action.ActionAttack] = 142
	fieryChargeFrames[action.ActionSkill] = 141
	fieryChargeFrames[action.ActionBurst] = 141
	fieryChargeFrames[action.ActionDash] = 142
	fieryChargeFrames[action.ActionWalk] = 141
	fieryChargeFrames[action.ActionSwap] = 138
	fieryChargeFrames[action.ActionHighPlunge] = 81
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		return c.fieryChargeAttack(), nil
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charged Attack",
		AdditionalTags:     []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagVaresaCombatCycle,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               charged[c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.09 * 60,
		CanBeDefenseHalted: false,
	}

	hitmark := chargeHitmark
	framesCB := frames.NewAbilFunc(chargeFrames)
	if c.StatusIsActive(skillStatus) {
		ai.Abil += " (Follow-Up Strike)"
		hitmark = fastChargeHitmark
		framesCB = quickAbilFunc(chargeFrames, fastChargedFrame)
		c.DeleteStatus(skillStatus)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1.2}, 2.8, 3.6),
		hitmark,
		hitmark,
	)
	return action.Info{
		Frames:          framesCB,
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   hitmark,
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) fieryChargeAttack() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Fiery Passion Charged Attack",
		AdditionalTags:   []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:        attacks.AttackTagExtra,
		ICDTag:           attacks.ICDTagVaresaCombatCycle,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Electro,
		Durability:       25,
		Mult:             fieryCharged[c.TalentLvlAttack()],
		HitlagFactor:     0.01,
		HitlagHaltFrames: 0.09 * 60,
	}

	hitmark := fieryChargeHitmark
	framesCB := frames.NewAbilFunc(fieryChargeFrames)
	if c.StatusIsActive(skillStatus) {
		ai.Abil += " (Follow-Up Strike)"
		hitmark = fastFieryChargeHitmark
		framesCB = quickAbilFunc(fieryChargeFrames, fastFieryChargedFrame)
		c.DeleteStatus(skillStatus)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1.2}, 2.8, 3.6),
		hitmark,
		hitmark,
	)
	return action.Info{
		Frames:          framesCB,
		AnimationLength: fieryChargeFrames[action.InvalidAction],
		CanQueueAfter:   hitmark,
		State:           action.ChargeAttackState,
	}
}

func quickAbilFunc(slice []int, skip int) func(action.Action) int {
	return func(next action.Action) int {
		if next == action.ActionHighPlunge {
			return slice[next] - skip
		}
		return slice[next]
	}
}
