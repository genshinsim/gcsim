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
	burstActivationDmgName     = "Anraku Secret Spring Therapy"
	burstKey                   = "mizuki-burst"
	burstHitmark               = 93
	burstDurability            = 50
	burstPoise                 = 100
	burstDuration              = 12 * 60
	burstCdDelay               = 1
	burstEnergyDrainDelay      = 4
	burstCd                    = 15 * 60
	burstRadius                = 8
	snackInterval              = 1.5 * 60
	snackSpawnOnEnemyRadius    = 6
	snackSpawnLocationVariance = 1.0
)

func init() {
	burstFrames = frames.InitAbilSlice(94) // Q -> Swap
	burstFrames[action.ActionAttack] = 93
	burstFrames[action.ActionCharge] = 92
	burstFrames[action.ActionSkill] = 93
	burstFrames[action.ActionDash] = 91
	burstFrames[action.ActionJump] = 93
	burstFrames[action.ActionWalk] = 92
}

// Summons forth countless lovely dreams and nightmares that pull in nearby objects and opponents,
// dealing AoE Anemo DMG and summoning a Mini Baku.
func (c *char) Burst(p map[string]int) (action.Info, error) {
	// Activation dmg
	ai := info.AttackInfo{
		ActorIndex:   c.Index(),
		Abil:         burstActivationDmgName,
		AttackTag:    attacks.AttackTagElementalBurst,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypeDefault,
		Element:      attributes.Anemo,
		Durability:   burstDurability,
		PoiseDMG:     burstPoise,
		Mult:         burstDMG[c.TalentLvlBurst()],
		HitlagFactor: 0.05,
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, burstRadius), burstHitmark, burstHitmark)

	// might be useful for checking in scripts
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
		CanQueueAfter:   burstFrames[action.ActionDash],
		State:           action.BurstState,
	}, nil
}

func (c *char) queueSnacks() {
	randomSign := func() float64 {
		rnd := c.Core.Rand.Float64()
		if 0.5 < rnd {
			return -1
		}
		return 1
	}
	snackFunc := func() {
		pos := c.calculateSnackSpawnLocation()

		// randomize a bit the spawn location
		pos.Y += c.Core.Rand.Float64() * snackSpawnLocationVariance * randomSign()
		pos.X += c.Core.Rand.Float64() * snackSpawnLocationVariance * randomSign()

		newSnack(c, pos)
	}

	// Spawn timer starts at burst hitmark
	spawnTime := burstHitmark
	for i := int(snackInterval); i <= burstDuration; i += snackInterval {
		c.Core.Tasks.Add(snackFunc, spawnTime+i)
	}
}

func (c *char) calculateSnackSpawnLocation() info.Point {
	// According to testing, snacks appear within a small range (1m) in front of the target/player.
	// However since the enemy direction is not set by default towards the player, we calculate
	// a position relative to the player/enemy
	playerPos := c.Core.Combat.Player().Pos()
	finalPosition := playerPos

	// find the closest enemy
	target := c.Core.Combat.ClosestEnemyWithinArea(
		combat.NewCircleHitOnTarget(playerPos, nil, snackSpawnOnEnemyRadius),
		nil,
	)

	// if enemy is found use this, otherwise use the player position
	if target != nil {
		targetShape := target.Shape()
		finalPosition = targetShape.Pos()
		direction := info.Point{
			X: playerPos.X - finalPosition.X,
			Y: playerPos.Y - finalPosition.Y,
		}
		if v, ok := targetShape.(*info.Circle); ok {
			if finalPosition != playerPos {
				direction = direction.Normalize()
				finalPosition.X += v.Radius() * direction.X
				finalPosition.Y += v.Radius() * direction.Y
			}
		} else if _, ok := targetShape.(*info.Rectangle); ok {
			// currently we cannot reliably get an edge to spawn the snack on rectangle.
			// place it somewhere around the middle
			finalPosition.X += direction.X / 2
			finalPosition.Y += direction.Y / 2
		}
	}
	return finalPosition
}
