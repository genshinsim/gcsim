package zibai

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	chargeFrames       []int
	chargeHitmarks     = []int{32, 32}
	chargeRadius       = []float64{2.3, 2.4}
	chargeOffsets      = []float64{1.4, 1.6}
	chargeHitHaltTime  = []float64{0.00, 0.06}
	chargeCanBeDefHalt = []bool{false, true}
)

func init() {
	chargeFrames = frames.InitAbilSlice(72)
	chargeFrames[action.ActionSkill] = 70
	chargeFrames[action.ActionBurst] = 70
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = 69
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillCharge()
	}

	for i, mult := range charge {
		ai := info.AttackInfo{
			Abil:               "Charge",
			ActorIndex:         c.Index(),
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeSlash,
			Element:            attributes.Physical,
			Durability:         25,
			Mult:               mult[c.TalentLvlAttack()],
			HitlagFactor:       chargeHitHaltTime[i],
			HitlagHaltFrames:   chargeHitHaltTime[i] * 60,
			CanBeDefenseHalted: chargeCanBeDefHalt[i],
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: chargeOffsets[i]}, chargeRadius[i])

		// this is okay because first multi-hit doesn't have hitlag
		c.Core.QueueAttack(
			ai,
			ap,
			chargeHitmarks[i],
			chargeHitmarks[i],
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) skillCharge() (action.Info, error) {
	for i, mult := range skillCharge {
		ai := info.AttackInfo{
			Abil:               "Charge",
			ActorIndex:         c.Index(),
			AttackTag:          attacks.AttackTagExtra,
			ICDTag:             attacks.ICDTagNormalAttack,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeBlunt,
			PoiseDMG:           60,
			Element:            attributes.Geo,
			IgnoreInfusion:     true,
			Durability:         25,
			Mult:               mult[c.TalentLvlSkill()],
			UseDef:             true,
			HitlagFactor:       chargeHitHaltTime[i],
			HitlagHaltFrames:   chargeHitHaltTime[i] * 60,
			CanBeDefenseHalted: chargeCanBeDefHalt[i],
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: chargeOffsets[i]}, chargeRadius[i])

		// this is okay because first multi-hit doesn't have hitlag
		c.Core.QueueAttack(
			ai,
			ap,
			chargeHitmarks[i],
			chargeHitmarks[i],
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}, nil
}
