package mizuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstActivationDmgName  = "Anraku Secret Spring Therapy"
	burstKey                = "mizuki-burst"
	burstHitmark            = 93
	burstDurability         = 50
	burstPoise              = 100
	burstDuration           = 12 * 60
	burstCdDelay            = 1
	burstEnergyDrainDelay   = 4
	burstCd                 = 15 * 60
	burstRadius             = 8
	snackDmgName            = "Munen Shockwave"
	snackHealName           = "Snack Pick-Up"
	snackInterval           = 1.5 * 60
	snackHitmark            = 22
	snackDurability         = 25
	snackPoise              = 30
	snackRadius             = 4
	snackHealTriggerHpRatio = 0.7
)

func init() {
	burstFrames = frames.InitAbilSlice(93)
	burstFrames[action.ActionCharge] = 92 // Q -> CA
	burstFrames[action.ActionDash] = 91   // Q -> D
	burstFrames[action.ActionWalk] = 92   // Q -> Walk
	burstFrames[action.ActionSwap] = 94   // Q -> Swap
}

// Summons forth countless lovely dreams and nightmares that pull in nearby objects and opponents,
// dealing AoE Anemo DMG and summoning a Mini Baku.
func (c *char) Burst(p map[string]int) (action.Info, error) {

	// Activation dmg
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       burstActivationDmgName,
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: burstDurability,
		PoiseDMG:   burstPoise,
		Mult:       burstDMG[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, burstRadius), burstHitmark, burstHitmark)

	// might be useful for checking
	c.AddStatus(burstKey, burstDuration, false)

	if c.Base.Cons >= 4 {
		c.c4EnergyGenerationsRemaining = c4EnergyGenerations
	}

	c.queueSnacks()

	c.ConsumeEnergy(burstEnergyDrainDelay)
	c.SetCDWithDelay(action.ActionBurst, burstCd, burstCdDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}

func (c *char) queueSnacks() {
	snackAttack := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       snackDmgName,
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: snackDurability,
		PoiseDMG:   snackPoise,
		Mult:       snackDMG[c.TalentLvlBurst()],
	}

	snackFunc := func() {
		heal := false
		dmg := false

		// C4 triggers both DMG and Heal
		if c.Base.Cons >= 4 {
			dmg = true
			heal = true
		} else {
			// Heals active char if is bellow 70% hp otherwise deals DMG
			dmg = c.Core.Player.ActiveChar().CurrentHP() > (c.Core.Player.ActiveChar().MaxHP() * snackHealTriggerHpRatio)
			heal = !dmg
		}

		// Not sure if it snapshots on creation (here) or on hitmark
		snapshot := c.Snapshot(&snackAttack)

		c.Core.Tasks.Add(func() {
			if dmg {
				// Snacks spawn on nearby enemy when enemies exist, otherwise on active character
				// Maybe implement it as gadget?
				c.Core.QueueAttackWithSnap(snackAttack, snapshot, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, snackRadius), 0)
			}

			if heal {
				// Heals double the amount on Mizuki
				healMultiplier := 1.0
				if c.Core.Player.Active() == c.Index {
					healMultiplier = 2.0
				}
				c.Core.Player.Heal(info.HealInfo{
					Caller:  c.Index,
					Target:  c.Core.Player.Active(),
					Message: snackHealName,
					Src:     ((c.Stat(attributes.EM) * snackHealEM[c.TalentLvlBurst()]) + snackHealFlat[c.TalentLvlBurst()]) * healMultiplier,
					Bonus:   c.Stat(attributes.Heal),
				})
			}

			// C4 restores 5 energy to mizuki up to 4 times
			if c.Base.Cons >= 4 {
				if c.c4EnergyGenerationsRemaining > 0 {
					c.c4EnergyGenerationsRemaining--
					c.AddEnergy(c4Key, c4Energy)
				}
			}
		}, snackHitmark)
	}

	spawnTime := snackHitmark
	for i := int(snackInterval); i <= burstDuration; i += snackInterval {
		spawnTime += snackInterval
		c.Core.Tasks.Add(snackFunc, spawnTime)
	}
}
