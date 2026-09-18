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

var burstTickHitmarks = []int{0, 0 + 33, 0 + 33 + 2 + 3, 0 + 33 + 2, 0 + 33 + 2 + 3 + 4}

const burstSpawnFrame = 36

func init() {
	burstFrames = make([][]int, 2)

	// Male, assuming same as female for now
	burstFrames[0] = frames.InitAbilSlice(75)
	burstFrames[0][action.ActionSkill] = 74 // Q -> E
	burstFrames[0][action.ActionJump] = 74  // Q -> J
	burstFrames[0][action.ActionSwap] = 73  // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(75)
	burstFrames[1][action.ActionSkill] = 74 // Q -> E
	burstFrames[1][action.ActionJump] = 74  // Q -> J
	burstFrames[1][action.ActionSwap] = 73  // Q -> Swap
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
	travel, ok := p["travel"]
	if !ok {
		travel = 6
	}

	attack := func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Frostbound Javelin",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       burst[c.TalentLvlBurst()] + flowGlowBonus[c.TalentLvlBurst()]*float64(c.frostglowStacks),
		}

		switch c.getRadiance() {
		case radianceStellarConduct:
			ai.Abil += stellarConductText
			ai.AttackTag = attacks.AttackTagDirectStellarConduct
			ai.Durability = 0
			ai.Mult = burstSSC[c.TalentLvlBurst()] + flowGlowBonusSSC[c.TalentLvlBurst()]*float64(c.frostglowStacks)
			ai.IgnoreDefPercent = 1
		case radianceStellarSwirl:
			ai.Abil += stellarSwirlText
			ai.AttackTag = attacks.AttackTagDirectStellarSwirl
			ai.Durability = 0
			ai.Mult = burstSSw[c.TalentLvlBurst()] + frostGlowBonusSSw[c.TalentLvlBurst()]*float64(c.frostglowStacks)
			ai.IgnoreDefPercent = 1
		default:
		}

		hits := 3
		if c.frostglowStacks == frostglowMax {
			hits += 2
		}

		for _, delay := range burstTickHitmarks[:hits] {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: -0.3}, 4.5),
				0,
				travel+delay,
			)
		}
		c.c6OnBurst(c.frostglowStacks)
		c.frostglowStacks = 0
	}

	c.QueueCharTask(attack, burstSpawnFrame)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(6)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
