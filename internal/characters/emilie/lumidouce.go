package emilie

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	lumidouceStatus        = "lumidouce-case"
	lumidouceLevel         = "lumidouce-level"
	lumidouceScent         = "lumidouce-scent"
	lumidouceScentCDKey    = "lumidouce-scent-cd"
	lumidouceScentResetKey = "lumidouce-scent-reset"

	lumidouceAttackTravel  = 5
	lumidouceHitmark       = 5 - lumidouceAttackTravel
	lumidouceHitmarkLevel2 = 23 - lumidouceAttackTravel

	lumidouceTickInterval       = 1.5 * 60
	lumidouceScentCD            = 2 * 60
	lumidouceScentResetInterval = 8 * 60
	lumidouceScentInterval      = 0.5 * 60
)

func (c *char) spawnLumidouceCase(level int, pos geometry.Point) {
	c.lumidouceSrc = c.Core.F
	c.lumidoucePos = pos
	c.SetTag(lumidouceLevel, level)
	c.SetTag(lumidouceScent, 0)
	c.AddStatus(lumidouceStatus, int(skillDuration[c.TalentLvlSkill()]*60), true)
	c.AddStatus(lumidouceScentResetKey, lumidouceScentResetInterval, true)
	c.QueueCharTask(c.lumidouceAttack(c.lumidouceSrc), lumidouceTickInterval)
	c.QueueCharTask(c.lumidouceOnBurning(c.lumidouceSrc), lumidouceScentInterval)
	c.QueueCharTask(c.lumidouceScentCollect(c.lumidouceSrc), lumidouceScentInterval)
}

func (c *char) lumidouceAttack(src int) func() {
	return func() {
		if c.lumidouceSrc != src {
			return
		}
		if !c.StatusIsActive(lumidouceStatus) {
			return
		}

		level := c.Tag(lumidouceLevel)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Lumidouce Case (Level %v)", level),
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagEmilieLumidouce,
			ICDGroup:   attacks.ICDGroupEmilieLumidouce,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Dendro,
			Durability: 25,
			Mult:       skillLumidouce[level-1][c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHit(c.lumidoucePos, c.Core.Combat.PrimaryTarget(), nil, 1)
		c.Core.QueueAttack(ai, ap, lumidouceHitmark+c.caseTravel, lumidouceHitmark+c.caseTravel, c.particleCB)
		if level == 2 {
			c.Core.QueueAttack(ai, ap, lumidouceHitmarkLevel2+c.caseTravel, lumidouceHitmarkLevel2+c.caseTravel, c.particleCB)
		}

		c.QueueCharTask(c.lumidouceAttack(src), lumidouceTickInterval)
	}
}

func (c *char) lumidouceOnBurning(src int) func() {
	return func() {
		if c.lumidouceSrc != src {
			return
		}
		if !c.StatusIsActive(lumidouceStatus) {
			return
		}

		if c.StatusIsActive(lumidouceScentCDKey) {
			c.QueueCharTask(c.lumidouceOnBurning(src), lumidouceScentInterval)
			return
		}

		generate := false
		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.lumidoucePos, nil, 20), nil)
		for _, v := range enemies {
			e, ok := v.(*enemy.Enemy)
			if !ok {
				continue
			}
			if e.IsBurning() {
				generate = true
				break
			}
		}
		if generate {
			c.generateScent()
			c.AddStatus(lumidouceScentCDKey, lumidouceScentCD, true)

			c.AddStatus(lumidouceScentResetKey, lumidouceScentResetInterval, true)
		}

		c.QueueCharTask(c.lumidouceOnBurning(src), lumidouceScentInterval)
	}
}

func (c *char) lumidouceScentCollect(src int) func() {
	return func() {
		if c.lumidouceSrc != src {
			return
		}
		if !c.StatusIsActive(lumidouceStatus) {
			return
		}

		if !c.StatusIsActive(lumidouceScentResetKey) && (c.Tag(lumidouceScent) > 0 || c.Tag(lumidouceLevel) > 1) {
			c.SetTag(lumidouceLevel, 1)
			c.Core.Log.NewEvent("scent reset", glog.LogCharacterEvent, c.Index)
		}

		if c.Tag(lumidouceScent) == 2 {
			if c.Tag(lumidouceLevel) < 2 {
				c.SetTag(lumidouceLevel, c.Tag(lumidouceLevel)+1)
				c.SetTag(lumidouceScent, 0)
			} else {
				c.a1()
			}
		}

		c.QueueCharTask(c.lumidouceScentCollect(src), lumidouceScentInterval)
	}
}

func (c *char) generateScent() {
	c.SetTag(lumidouceScent, c.Tag(lumidouceScent)+1)

	c.Core.Log.NewEvent("scent generated", glog.LogCharacterEvent, c.Index).
		Write("level", c.Tag(lumidouceLevel)).
		Write("scent", c.Tag(lumidouceScent))
}
