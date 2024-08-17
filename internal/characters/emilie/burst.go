package emilie

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const (
	burstMarkKey = "emilie-burst-mark"

	burstSpawn          = 96
	burstResetLumidouce = 306

	burstRadius       = 12
	burstHitmark      = 12
	burstTickInterval = 0.3 * 60
	burstMarkDuration = 0.7 * 60
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(111) // Q -> E
	burstFrames[action.ActionAttack] = 108
	burstFrames[action.ActionDash] = 97
	burstFrames[action.ActionJump] = 98
	burstFrames[action.ActionWalk] = 96
	burstFrames[action.ActionSwap] = 105
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	var ok bool
	c.caseTravel, ok = p["travel"]
	if !ok {
		c.caseTravel = lumidouceAttackTravel
	}

	c.QueueCharTask(func() {
		c.prevLumidouceLvl = 1
		if c.StatusIsActive(lumidouceStatus) {
			c.prevLumidouceLvl = c.Tag(lumidouceLevel)
		}

		c.spawnBurstLumidouceCase()
		c.c6()
	}, burstSpawn)
	c.QueueCharTask(func() {
		c.spawnLumidouceCase(c.prevLumidouceLvl, c.lumidoucePos)
	}, burstResetLumidouce)

	duration := int(burstCD[c.TalentLvlBurst()] * 60)
	c.burstMarkDuration = burstMarkDuration
	if c.Base.Cons >= 4 {
		duration += 2 * 60
		c.burstMarkDuration -= 0.3 * 60
	}

	c.ConsumeEnergy(107)
	c.SetCDWithDelay(action.ActionBurst, duration, 97)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionWalk], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) spawnBurstLumidouceCase() {
	player := c.Core.Combat.Player()

	c.lumidouceSrc = c.Core.F
	c.lumidoucePos = geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 2.1}, player.Direction())
	c.SetTag(lumidouceLevel, 3)
	c.SetTag(lumidouceScent, 0)
	c.AddStatus(lumidouceStatus, int(burstDuration[c.TalentLvlBurst()]*60), true)
	c.QueueCharTask(c.lumidouceBurstAttack(c.lumidouceSrc), burstTickInterval)
}

func (c *char) lumidouceBurstAttack(src int) func() {
	return func() {
		if c.lumidouceSrc != src {
			return
		}
		if !c.StatusIsActive(lumidouceStatus) {
			return
		}

		burstArea := combat.NewCircleHitOnTarget(c.lumidoucePos, nil, burstRadius)
		pos := c.getRandomEnemyPosition(burstArea)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Lumidouce Case (Level 3)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Dendro,
			Durability: 25,
			Mult:       burstDMG[c.TalentLvlBurst()],
		}
		ap := combat.NewCircleHitOnTarget(pos, nil, 2.5)
		c.Core.QueueAttack(ai, ap, burstHitmark, burstHitmark, c.particleCB, c.c2)

		c.QueueCharTask(c.lumidouceBurstAttack(src), burstTickInterval)
	}
}

func (c *char) getRandomEnemyPosition(area combat.AttackPattern) geometry.Point {
	enemy := c.Core.Combat.RandomEnemyWithinArea(
		area,
		func(e combat.Enemy) bool {
			return !e.StatusIsActive(burstMarkKey)
		},
	)
	var pos geometry.Point
	if enemy != nil {
		pos = enemy.Pos()
		enemy.AddStatus(burstMarkKey, c.burstMarkDuration, true) // same enemy can't be targeted again for 0.7s
	} else {
		pos = geometry.CalcRandomPointFromCenter(area.Shape.Pos(), 0.5, burstRadius, c.Core.Rand)
	}
	return pos
}
