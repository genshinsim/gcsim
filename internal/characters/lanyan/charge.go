package lanyan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	chargeFrames []int

	chargeHitmarks = []int{42, 49, 56}
)

func init() {
	chargeFrames = frames.InitAbilSlice(70) // CA -> Jump
	chargeFrames[action.ActionAttack] = 41
	chargeFrames[action.ActionSkill] = 44
	chargeFrames[action.ActionBurst] = 44
	chargeFrames[action.ActionDash] = 59
	chargeFrames[action.ActionWalk] = 59
	chargeFrames[action.ActionSwap] = 45
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Charged Attack",
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Anemo,
		Durability:   25,
		Mult:         charge[c.TalentLvlAttack()],
		IsDeployable: true,
	}

	for _, hitmark := range chargeHitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 2.5),
			hitmark,
			hitmark,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeFrames[action.ActionAttack],
		State:           action.ChargeAttackState,
	}, nil
}
