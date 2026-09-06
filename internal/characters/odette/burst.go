package odette

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	burstFrames   []int
	burstHitmarks = []int{113, 113 + 5, 113 + 5 + 5}
	finalHitmark  = 113 + 5 + 5 + 5
)

const (
	burstSummonFrame    = 99
	burstFirstTickDelay = 134
	swansDreamKey       = "odette-snow-swans-dream"
)

func init() {
	burstFrames = frames.InitAbilSlice(126) // Q -> W
	burstFrames[action.ActionAttack] = 108  // Q -> N1
	burstFrames[action.ActionSkill] = 107   // Q -> E
	burstFrames[action.ActionDash] = 108    // Q -> D
	burstFrames[action.ActionJump] = 109    // Q -> J
	burstFrames[action.ActionSwap] = 106    // Q -> Swap
}

// With quick, lively dance steps, Odette deals multiple instances of AoE Cryo DMG, and summons her
// Solo Dance Double. She also gains Snow Swan's Dream, which increases the Stellar Glimmer reaction
// DMG Odette deals.
//
// Additionally, for 6s after unleashing the Elemental Burst Presto: Bluebird Finale, Odette's
// Elemental Skill Adagio: Phantom Night Dancers will be replaced with the special Elemental Skill
// Adagio: Coda at Dawn's Tolling instead.
//
// If there is a Solo Dance Double summoned by Odette on the field, it will be summoned to her side
// with its duration refreshed.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Presto: Bluebird Finale (Slash)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: -5},
		14,
		12,
	)

	for _, delay := range burstHitmarks {
		c.QueueCharTask(func() { c.Core.QueueAttack(ai, ap, 0, 0) }, delay)
	}

	c.QueueCharTask(func() {
		aiFinal := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Presto: Bluebird Finale (Final)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       burstFinal[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(aiFinal, ap, 0, 0)

		c.addSwansDreamBuff()
	}, finalHitmark)

	c.QueueCharTask(func() { c.summonDanceDouble(burstFirstTickDelay) }, burstSummonFrame)

	c.AddStatus(skillRecastKey, 6*60+burstSummonFrame, false)
	c.ConsumeEnergy(7)
	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) addSwansDreamBuff() {
	buff := burstBuff[c.TalentLvlBurst()]
	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBaseWithHitlag(swansDreamKey, 20*60),
		Amount: func(ai info.AttackInfo) float64 {
			switch ai.AttackTag {
			case attacks.AttackTagDirectStellarConduct,
				attacks.AttackTagDirectStellarSwirl,
				attacks.AttackTagReactionStellarSwirl:
				return buff
			default:
				return 0
			}
		},
	})
	c.c4OnBurst(buff)
}
