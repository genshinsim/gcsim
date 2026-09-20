package zibai

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

var burstHitmarks = []int{95, 95}

func init() {
	burstFrames = frames.InitAbilSlice(103) // Q -> N1/E
	burstFrames[action.ActionAttack] = 97
	burstFrames[action.ActionSkill] = 98
	burstFrames[action.ActionDash] = 98
	burstFrames[action.ActionJump] = 99
	burstFrames[action.ActionSwap] = 96
}

// Zibai operates the Jadelight Canopy, dealing two instances of Geo DMG, with the second damage
// instance being considered Lunar-Crystallize Reaction DMG.
// When cast, if Zibai is in the Lunar Phase Shift, the duration of the current Lunar Phase Shift
// will extend by an additional 1.7s.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	// deal damage when created
	for i, mult := range burst {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Burst",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   75,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       mult[c.TalentLvlBurst()],
			UseDef:     true,
		}
		if i == 1 {
			ai.Abil += lunarCrystallizeAbil
			ai.AttackTag = attacks.AttackTagDirectLunarCrystallize
			ai.Durability = 0
			ai.IgnoreDefPercent = 1
			ai.PoiseDMG = 100
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.7),
			burstHitmarks[i],
			burstHitmarks[i],
		)
	}

	c.ExtendStatus(skillKey, 1.7*60)
	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
