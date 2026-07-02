package varka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames  []int
	burstHitmark = []int{95, 100}
)

func init() {
	burstFrames = frames.InitAbilSlice(110) // Q -> E
	burstFrames[action.ActionAttack] = 105  // Q -> N1
	burstFrames[action.ActionDash] = 105    // Q -> D
	burstFrames[action.ActionJump] = 105    // Q -> J
	burstFrames[action.ActionWalk] = 105    // Q -> Walk
	burstFrames[action.ActionSwap] = 105    // Q -> Swap
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) (action.Info, error) {
	ele := []attributes.Element{c.conversionElem, attributes.Anemo}
	for i, hitmark := range burstHitmark {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Divine Maiden's Deliverance (Initial)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   250,
			Element:    ele[i],
			Durability: 25,
			Mult:       burst[i][c.TalentLvlBurst()],
		}
		burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 2}, 8)
		burstPos := burstArea.Shape.Pos()
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTargetFanAngle(burstPos, nil, 8, 120),
			hitmark,
			hitmark,
		)
	}

	// apparently extends E by 2.3s even though it's not in the description
	c.ExtendStatus(skillKey, 2.3*60)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(4)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
