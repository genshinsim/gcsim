package mizuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
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
	snackSpawnRadius        = 4 // assumption
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

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, burstRadius), burstHitmark, burstHitmark)

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
	snackFunc := func() {

		pos := c.calculateSnackSpawnLocation()

		// randomize a bit the spawn location
		pos.Y += c.Core.Rand.Float64() * 0.2
		pos.X += c.Core.Rand.Float64() * 0.2

		newSnack(c, pos)
	}

	spawnTime := burstHitmark
	for i := int(snackInterval); i <= burstDuration; i += snackInterval {
		spawnTime += snackInterval
		c.Core.Tasks.Add(snackFunc, spawnTime)
	}
}

func (c *char) calculateSnackSpawnLocation() geometry.Point {
	playerPos := c.Core.Combat.Player().Pos()
	finalPosition := playerPos

	// find the closest enemy
	target := c.Core.Combat.ClosestEnemyWithinArea(combat.AttackPattern{
		Shape: geometry.NewCircle(playerPos, snackSpawnRadius, geometry.DefaultDirection(), 360),
	}, nil)

	// if no enemy found, it spawns on active character
	if target != nil {

		// spawns around the enemy bounds, not inside
		targetShape := target.Shape()
		direction := geometry.Point{
			X: playerPos.X - targetShape.Pos().X,
			Y: playerPos.Y - targetShape.Pos().Y,
		}
		pos := target.Pos()
		if v, ok := targetShape.(*geometry.Circle); ok {
			direction = direction.Normalize()
			pos.X += v.Radius() * direction.X
			pos.Y += v.Radius() * direction.Y
		} else if _, ok := targetShape.(*geometry.Rectangle); ok {
			// currently we cannot reliably get an edge to spawn the snack on rectangle.
			// place it somewhere around the middle
			pos.X += direction.X / 2
			pos.Y += direction.Y / 2
		}
	}
	return finalPosition
}
