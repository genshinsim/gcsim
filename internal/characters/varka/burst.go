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
	burstFrames   []int
	burstHitmark  = []int{112, 112 + 19}
	burstPoiseDmg = []float64{130, 70}
)

func init() {
	burstFrames = frames.InitAbilSlice(152)   // Q -> CA
	burstFrames[action.ActionAttack] = 137    // Q -> N1
	burstFrames[action.ActionSkill] = 138     // Q -> E
	burstFrames[action.ActionDash] = 157 - 19 // Q -> D
	burstFrames[action.ActionJump] = 167 - 30 // Q -> J
	burstFrames[action.ActionWalk] = 140      // Q -> Walk
	burstFrames[action.ActionSwap] = 136      // Q -> Swap
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) (action.Info, error) {
	ele := []attributes.Element{c.conversionElem, attributes.Anemo}
	for i, hitmark := range burstHitmark {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Northwind Avatar",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   burstPoiseDmg[i],
			Element:    ele[i],
			Durability: 25,
			Mult:       burst[i][c.TalentLvlBurst()],
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 2}, 7),
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
