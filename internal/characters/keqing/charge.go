package keqing

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int
var chargeHitmarks = []int{22, 24}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagExtra,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), chargeHitmarks[i], chargeHitmarks[i])
	}

	if c.Core.Status.Duration(stilettoKey) > 0 {
		// despawn stiletto
		c.Core.Status.Delete(stilettoKey)

		//2 hits
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Thunderclap Slash",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Electro,
			Durability: 50,
			Mult:       skillCA[c.TalentLvlSkill()],
		}
		for i := 0; i < 2; i++ {
			c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), chargeHitmarks[i], chargeHitmarks[i])
		}

		// TODO: Particle timing?
		if c.Core.Rand.Float64() < .5 {
			c.Core.QueueParticle("keqing", 2, attributes.Electro, 100)
		} else {
			c.Core.QueueParticle("keqing", 3, attributes.Electro, 100)
		}
	}

	if c.Base.Cons >= 6 {
		c.c6("charge")
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[len(chargeHitmarks)-1],
		Post:            chargeHitmarks[len(chargeHitmarks)-1],
		State:           action.ChargeAttackState,
	}
}
