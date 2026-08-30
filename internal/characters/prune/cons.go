package prune

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1ICDKey        = "prune-c1-icd"
	c4ICDKey        = "prune-c4-icd"
	c6WindowKey     = "prune-c6-window"
	c6ActiveBuffKey = "prune-c6-active-buff"
)

func (c *char) makeC1CB(a info.AttackCB) {
	if c.Base.Cons < 1 || a.Target.Type() != info.TargettableEnemy || c.StatusIsActive(c1ICDKey) {
		return
	}

	c.AddStatus(c1ICDKey, 1.8*60, false)
	c.AddEnergy("prune-c1", 2)
}

func (c *char) c2(duration int) {
	if c.Base.Cons < 2 {
		return
	}

	c.c2Buff[attributes.ATKP] = 0.10
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("prune-c2", duration),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return c.c2Buff
		},
	})
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2Buff = make([]float64, attributes.EndStatType)
}

func (c *char) makeC2CB(a info.AttackCB) {
	if c.Base.Cons < 2 || a.Target.Type() != info.TargettableEnemy || !c.StatusIsActive(burstKey) {
		return
	}

	c.c2Buff[attributes.ATKP] = math.Min(c.c2Buff[attributes.ATKP]+0.05, 0.40)
}

func (c *char) c4Ricochet(ele attributes.Element, tag attacks.AttackTag) info.AttackCBFunc {
	return func(a info.AttackCB) {
		if c.Base.Cons < 4 {
			return
		}

		// has to be on the enemy hit by the attack that call this
		if a.Target.Type() != info.TargettableEnemy {
			return
		}

		// C4 has 0.1s ICD
		if c.StatusIsActive(c4ICDKey) {
			return
		}
		c.AddStatus(c4ICDKey, 6, false)

		// default to ricocheting onto the enemy that the hammer hit
		target := a.Target

		// prefer a different nearby enemy when one exists, the exact radius of detection is unknown
		next := c.Core.Combat.RandomEnemyWithinArea(combat.NewCircleHitOnTarget(a.Target, nil, 8), func(t info.Enemy) bool {
			return t.Key() != a.Target.Key()
		},
		)

		if next != nil {
			target = next
		}

		// queue the attack
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Banehunter Oathhammer Ricochet (C4)",
			AttackTag:  tag,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   15,
			Element:    ele,
			Durability: 0,
			Mult:       0.8,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(target, nil, 1.5),
			0,
			63,
			c.makeA4CB,
			c.makeC1CB,
			c.makeC2CB,
		)
	}
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	chars := c.Core.Player.Chars()

	c.c6Buff = make([]float64, attributes.EndStatType)
	c.c6Buff[attributes.ATK] = 350

	buff := func(args ...any) {
		// reaction target must be an enemy
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		// reaction must have been triggered by a character attack
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		// triggerer must be affected by tolling rally
		triggerer := chars[atk.Info.ActorIndex]
		if !triggerer.StatusIsActive(a4Key) {
			return
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("prune-c6-self-buff", (5+0.5)*60),
			AffectedStat: attributes.ATK,
			Amount: func() []float64 {
				return c.c6Buff
			},
		})

		// the status last for 5.5s, with a check every 0.5s
		c.AddStatus(c6WindowKey, (5+0.5)*60, true)

		c.c6TickSrc++
		c.c6Tick(c.c6TickSrc)
	}

	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Core.Events.Subscribe(i, buff, "prune-c6-buff")
	}
}

func (c *char) c6BurstBonusDur() int {
	if c.Base.Cons >= 6 {
		return 0
	}
	return 4 * 60
}

func (c *char) c6Tick(src int) {
	if src != c.c6TickSrc {
		return
	}

	active := c.Core.Player.Active()

	for _, char := range c.Core.Player.Chars() {
		// not apply to prune
		if char.Index() == c.Index() {
			char.DeleteStatMod(c6ActiveBuffKey)
			continue
		}

		// has to be the active character with tolling rally
		if char.Index() == active && c.StatusIsActive(c6WindowKey) && char.StatusIsActive(a4Key) {

			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(c6ActiveBuffKey, -1),
				AffectedStat: attributes.ATK,
				Amount: func() []float64 {
					return c.c6Buff
				},
			})
			continue
		}

		char.DeleteStatMod(c6ActiveBuffKey)
	}

	if !c.StatusIsActive(c6WindowKey) {
		return
	}

	c.Core.Tasks.Add(func() { c.c6Tick(src) }, 30)
}
