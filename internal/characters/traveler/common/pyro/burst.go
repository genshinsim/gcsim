package pyro

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames [][]int

const (
	burstHitmark = 40
)

var nightsoulGainDelays = []int{42, 52, 59, 59}

func init() {
	burstFrames = make([][]int, 2)

	// Male
	burstFrames[0] = frames.InitAbilSlice(49) // Q -> N1
	burstFrames[0][action.ActionSkill] = 48   // Q -> E, Eh
	burstFrames[0][action.ActionSwap] = 47

	// Female
	burstFrames[1] = frames.InitAbilSlice(49) // Q -> N1
	burstFrames[1][action.ActionSkill] = 48   // Q -> E, Eh
	burstFrames[1][action.ActionSwap] = 47
}

func (c *Traveler) Burst(p map[string]int) (action.Info, error) {
	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(19)

	c.c4AddMod()

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Plains Scorcher",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4.5),
		burstHitmark,
		burstHitmark,
	)

	c.QueueCharTask(c.nightsoulGainFunc(0), nightsoulGainDelays[0])

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *Traveler) nightsoulGainFunc(count int) func() {
	return func() {
		c.nightsoulState.GeneratePoints(7)
		if count < 3 {
			c.QueueCharTask(c.nightsoulGainFunc(count+1), nightsoulGainDelays[count+1])
		}
	}
}
