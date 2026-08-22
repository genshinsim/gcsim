package prune

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1ICDKey    = "prune-c1-icd"
	c4ICDKey    = "prune-c4-icd"
	c6WindowKey = "prune-c6-window"
)

func (c *char) makeC1CB(a info.AttackCB) {
	if c.Base.Cons < 1 || a.Target.Type() != info.TargettableEnemy || c.StatusIsActive(c1ICDKey) {
		return
	}

	c.AddStatus(c1ICDKey, 108, false)
	c.AddEnergy("prune-c1", 2)
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

	// prune always gets the C6 ATK buff while the C6 window is active
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("prune-c6-self-buff", -1),
		AffectedStat: attributes.ATK,
		Amount: func() []float64 {
			if !c.StatusIsActive(c6WindowKey) {
				return nil
			}
			return c.c6Buff
		},
	})

	// give the currently active character 350 ATK
	for _, ch := range chars {
		char := ch

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("prune-c6-active-buff", -1),
			AffectedStat: attributes.ATK,
			Amount: func() []float64 {
				// C6 must currently be active
				if !c.StatusIsActive(c6WindowKey) {
					return nil
				}

				// this character must currently be on field
				if c.Core.Player.Active() != char.Index() {
					return nil
				}

				// current active character must have Tolling Rally
				if !char.StatusIsActive(a4Key) {
					return nil
				}

				return c.c6Buff
			},
		})
	}

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

		c.AddStatus(c6WindowKey, 5*60, true)

		c.Core.Log.NewEvent("prune c6 triggered", glog.LogCharacterEvent, c.Index()).
			Write("atk", c.c6Buff[attributes.ATK]).
			Write("expiry", c.Core.F+5*60)
	}

	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Core.Events.Subscribe(i, buff, "prune-c6-buff")
	}
}
