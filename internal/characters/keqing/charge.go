package keqing

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	chargeFrames   []int
	chargeHitmarks = []int{22, 24}
	chargeRadius   = []float64{2.2, 2.3}
	chargeOffsets  = []float64{1.5, 1.8}
)

func init() {
	chargeFrames = frames.InitAbilSlice(36)
	chargeFrames[action.ActionSkill] = 35
	chargeFrames[action.ActionBurst] = 35
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = chargeHitmarks[len(chargeHitmarks)-1]
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}
	c2CB := c.makeC2CB()
	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(
				c.Core.Combat.Player(),
				combat.Point{Y: chargeOffsets[i]},
				chargeRadius[i],
			),
			chargeHitmarks[i],
			chargeHitmarks[i],
			c2CB,
		)
	}

	if c.Core.Status.Duration(stilettoKey) > 0 {
		// despawn stiletto
		c.Core.Status.Delete(stilettoKey)

		//2 hits
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Thunderclap Slash",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeSlash,
			Element:    attributes.Electro,
			Durability: 50,
			Mult:       skillCA[c.TalentLvlSkill()],
		}
		for i := 0; i < 2; i++ {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5),
				chargeHitmarks[i],
				chargeHitmarks[i],
				c.particleCB,
			)
		}
	}

	if c.Base.Cons >= 6 {
		c.c6("charge")
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
