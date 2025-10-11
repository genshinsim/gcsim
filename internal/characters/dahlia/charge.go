package dahlia

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	chargeFrames   []int
	chargeHitmarks = []int{10, 10} // CA-1 and CA-2 hit at the same time
)

func init() {
	chargeFrames = frames.InitAbilSlice(67)                                 // CA -> N1
	chargeFrames[action.ActionSkill] = 42                                   // CA -> tE / hE (TO-DO: check if this includes both)
	chargeFrames[action.ActionBurst] = 42                                   // CA -> Q
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1] // CA -> D
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1] // CA -> J
	chargeFrames[action.ActionSwap] = 42                                    // CA -> Swap
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	for i, mult := range charge {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.2), // TO-DO: copied from Bennett
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
