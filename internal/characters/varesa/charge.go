package varesa

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	chargeFrames      []int
	fieryChargeFrames []int
)

const (
	chargeHitmark      = 24
	fieryChargeHitmark = 19

	boostChargeAnim = 1.0 - 0.9
)

// TODO: update frames

func init() {
	// based on heizou frames
	chargeFrames = frames.InitAbilSlice(46)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionAttack] = 38
	chargeFrames[action.ActionSkill] = 38
	chargeFrames[action.ActionBurst] = 38

	// based on wriothesley frames
	fieryChargeFrames = frames.InitAbilSlice(52) // CA -> N1/E/Q
	fieryChargeFrames[action.ActionDash] = fieryChargeHitmark
	fieryChargeFrames[action.ActionJump] = fieryChargeHitmark
	fieryChargeFrames[action.ActionWalk] = 51
	fieryChargeFrames[action.ActionSwap] = 49
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
	framesCB := frames.NewAbilFunc
	if c.StatusIsActive(skillStatus) {
		ai.Abil += " (Follow-Up Strike)"
		hitmark = int(math.Round(chargeHitmark * boostChargeAnim))
		framesCB = quickAbilFunc
		c.DeleteStatus(skillStatus)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1.2}, 2.8, 3.6),
		hitmark,
		hitmark,
	)
	return action.Info{
		Frames:          framesCB(chargeFrames),
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
	framesCB := frames.NewAbilFunc
	if c.StatusIsActive(skillStatus) {
		ai.Abil += " (Follow-Up Strike)"
		hitmark = int(math.Round(fieryChargeHitmark * boostChargeAnim))
		framesCB = quickAbilFunc
		c.DeleteStatus(skillStatus)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1.2}, 2.8, 3.6),
		hitmark,
		hitmark,
	)
	return action.Info{
		Frames:          framesCB(fieryChargeFrames),
		AnimationLength: fieryChargeFrames[action.InvalidAction],
		CanQueueAfter:   hitmark,
		State:           action.ChargeAttackState,
	}
}

// TODO: custom frames or this function?
func quickAbilFunc(slice []int) func(action.Action) int {
	return func(next action.Action) int {
		return int(float64(slice[next]) * boostChargeAnim)
	}
}
