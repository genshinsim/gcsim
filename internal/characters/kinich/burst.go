package kinich

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var burstFrames []int
var ajawHitmarks = []int{253, 150}

const (
	cdStart            = 1
	consumeEnergyDelay = 5

	burstHitMark = 161
	ajawDuration = 15*60 + burstHitMark

	burstKey = "ajaw"
)

func init() {
	burstFrames = frames.InitAbilSlice(126) // Q -> E
	burstFrames[action.ActionAttack] = 125
	burstFrames[action.ActionDash] = 124
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.ajawSrc = c.Core.F
	c.AddStatus(burstKey, ajawDuration, false)
	if c.nightsoulState.HasBlessing() {
		// extend Nightsoul's Blessing time limit countdown
		duration := (c.exitStateF - c.Core.F) + 1.7*60
		c.setNightsoulExitTimer(duration)
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Hail to the Almighty Dragonlord (Skill DMG)",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagElementalBurst,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Dendro,
		Durability:     25,
		Mult:           burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4), burstHitMark, burstHitMark)
	c.Core.Tasks.Add(c.QueueLaser(1, c.ajawSrc), ajawHitmarks[0])
	c.ConsumeEnergy(consumeEnergyDelay)
	c.SetCDWithDelay(action.ActionBurst, 18*60, cdStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) QueueLaser(step, src int) func() {
	return func() {
		if c.ajawSrc != src {
			return
		}
		// duration expired
		if !c.StatusIsActive(burstKey) {
			return
		}
		// condition to track number of hits just in case
		if step == 7 {
			c.DeleteStatus(burstKey)
			return
		}
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Hail to the Almighty Dragonlord (Dragon Breath DMG)",
			AttackTag:      attacks.AttackTagElementalBurst,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagElementalBurst,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Dendro,
			Durability:     25,
			Mult:           burst[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, 0)
		c.Core.Tasks.Add(c.QueueLaser(step+1, src), ajawHitmarks[1])
	}
}
