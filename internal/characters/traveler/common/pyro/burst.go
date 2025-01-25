package pyro

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames [][]int

const (
	burstHitmark = 38
)

func init() {
	burstFrames = make([][]int, 2)

	// Male
	burstFrames[0] = frames.InitAbilSlice(58)
	burstFrames[0][action.ActionSwap] = 57 // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(58)
	burstFrames[1][action.ActionSwap] = 57 // Q -> Swap
}

func (c *Traveler) Burst(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		c.SetCD(action.ActionBurst, 18*60)
	}, 2)
	c.ConsumeEnergy(5)

	c.c4AddMod()

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Plains Scorcher",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           burst[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5),
		burstHitmark,
		burstHitmark,
	)

	c.nightsoulGainFunc(0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *Traveler) nightsoulGainFunc(count int) func() {
	return func() {
		if count > 3 {
			return
		}
		c.nightsoulState.GeneratePoints(7)
		c.QueueCharTask(c.nightsoulGainFunc(count+1), 60)
	}
}
