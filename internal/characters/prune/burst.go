package prune

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const burstKey = "prune-burst"

func init() {
	burstFrames = frames.InitAbilSlice(73) // walk
	burstFrames[action.ActionAttack] = 73
	burstFrames[action.ActionCharge] = 82
	burstFrames[action.ActionSkill] = 72
	burstFrames[action.ActionDash] = 74
	burstFrames[action.ActionJump] = 73
	burstFrames[action.ActionSwap] = 71
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.a1ConvertEle = attributes.NoElement
	c.burstEnergyDrainDelay = 16

	duration := 813
	if c.Base.Cons >= 6 {
		duration += 4 * 60 // check frame
	}

	c.AddStatus(burstKey, duration, false)

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "The Bell Tolls! The Hunt Is On!",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5),
		0,
		34,
	)

	c.burstSrc = c.Core.F
	for i := 137; i < duration; i = i + 117 {
		c.Core.Tasks.Add(c.burstTick(c.burstSrc), i)
	}

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(c.burstEnergyDrainDelay)

	if c.Base.Cons >= 2 {
		c.c2Buff[attributes.ATKP] = 0.1
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTick(src int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		ai := info.AttackInfo{
			ActorIndex:   c.Index(),
			Abil:         "Witchlure Bell",
			AttackTag:    attacks.AttackTagElementalBurst,
			ICDTag:       attacks.ICDTagNone,
			ICDGroup:     attacks.ICDGroupDefault,
			Element:      attributes.Anemo,
			Durability:   25,
			Mult:         burstDot[c.TalentLvlBurst()],
			HitlagFactor: 0.01,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.6),
			0,
			0,
		)

		c.a1()

		if c.Base.Cons >= 2 {
			c.c2Buff[attributes.ATKP] = math.Min(c.c2Buff[attributes.ATKP]+0.05, 0.40)
		}
	}
}
