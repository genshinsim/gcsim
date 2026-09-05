package cryo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames [][]int

var burstTickHitmarks = []int{60, 60 + 15, 60 + 15 + 6, 60 + 15 + 6 + 6, 60 + 15 + 6 + 6 + 6}

const burstSpawnFrame = 55

func init() {
	burstFrames = make([][]int, 2)

	// TODO: Placeholder using DMC frames

	// Male
	burstFrames[0] = frames.InitAbilSlice(58)
	burstFrames[0][action.ActionSwap] = 57 // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(58)
	burstFrames[1][action.ActionSwap] = 57 // Q -> Swap
}

// Generates an ice javelin with the power of Cryo, then directs it at enemies to deal multiple
// instances of Cryo DMG.
//
// When cast, the Traveler consumes all existing stacks of Frostglow, which in turn increases the
// DMG dealt by this Elemental Burst. When 8 stacks of Frostglow are consumed, the number of DMG
// instances caused by the javelins is also increased.
//
// Radiance: Stellar Glimmer: DMG from the current Elemental Burst is changed to Cryo DMG of the
// corresponding Stellar Glimmer reaction type.
func (c *Traveler) Burst(p map[string]int) (action.Info, error) {
	attack := func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Frostbound Javelin",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       burst[c.TalentLvlBurst()] + flowGlowBonus[c.TalentLvlBurst()]*float64(c.flostglowStacks),
		}

		switch c.getRadiance() {
		case radianceStellarConduct:
			ai.Abil += stellarConductText
			ai.AttackTag = attacks.AttackTagDirectStellarConduct
			ai.Durability = 0
			ai.Mult = burstSSC[c.TalentLvlBurst()] + flowGlowBonusSSC[c.TalentLvlBurst()]*float64(c.flostglowStacks)
			ai.IgnoreDefPercent = 1
		case radianceStellarSwirl:
			ai.Abil += stellarSwirlText
			ai.AttackTag = attacks.AttackTagDirectStellarSwirl
			ai.Durability = 0
			ai.Mult = burstSSw[c.TalentLvlBurst()] + flowGlowBonusSSw[c.TalentLvlBurst()]*float64(c.flostglowStacks)
			ai.IgnoreDefPercent = 1
		default:
		}

		hits := 3
		if c.flostglowStacks == skillStacksMax {
			hits += 2
		}

		for _, delay := range burstTickHitmarks[:hits] {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4.5),
				delay-burstSpawnFrame,
				delay-burstSpawnFrame,
			)
		}
		c.c6OnBurst(c.flostglowStacks)
		c.flostglowStacks = 0
	}

	c.QueueCharTask(attack, burstSpawnFrame)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
