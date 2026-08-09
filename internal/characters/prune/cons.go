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
	c4ICDKey = "prune-c4-icd"
	c1ICDKey = "prune-c1-icd"
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

	c.c2Buff[attributes.ATKP] =
		math.Min(c.c2Buff[attributes.ATKP]+0.05, 0.40)
}

func (c *char) c4Ricochet(ele attributes.Element, tag attacks.AttackTag) info.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}

	return func(a info.AttackCB) {
		// has to be on the enemy hit by the attack that call this
		if a.Target.Type() != info.TargettableEnemy {
			return
		}

		// C1 has 0.1s ICD
		if c.StatusIsActive(c4ICDKey) {
			return
		}
		c.AddStatus(c4ICDKey, 6, false)

		// default to ricocheting onto the enemy that the hammer hit
		var target info.Target = a.Target

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
			Element:    ele,
			Durability: 0,
			Mult:       0.8,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(target, nil, 1.5),
			0,
			0, // need to check frame
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

	lastTriggerFrame := make([]int, len(c.Core.Player.Chars()))
	for i := range lastTriggerFrame {
		lastTriggerFrame[i] = -1
	}

	chars := c.Core.Player.Chars()

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

		// prevent multiple qualifying reactions on the same frame from triggering again
		triggererIndex := atk.Info.ActorIndex
		if lastTriggerFrame[triggererIndex] == c.Core.F {
			return
		}
		lastTriggerFrame[triggererIndex] = c.Core.F

		// triggerer must be affected by tolling rally
		triggerer := chars[atk.Info.ActorIndex]
		if !triggerer.StatusIsActive(a4Key) {
			return
		}

		// apply buff
		c.c6Buff = make([]float64, attributes.EndStatType)
		c.c6Buff[attributes.ATK] = 350

		// for prune
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("prune-c6-buff", 5*60), // 5 s
			AffectedStat: attributes.ATK,
			Amount: func() []float64 {
				return c.c6Buff
			},
		})

		// for active character
		c.Core.Player.ActiveChar().AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("prune-c6-buff", 5*60), // 5 s
			AffectedStat: attributes.ATK,
			Amount: func() []float64 {
				return c.c6Buff
			},
		})

		c.Core.Log.NewEvent("prune c6 triggered", glog.LogCharacterEvent, c.Index()).
			Write("atk", c.c6Buff[attributes.ATK]).
			Write("expiry", c.Core.F+5*60)

	}

	for i := event.ReactionEventStartDelim + 1; i < event.ReactionEventEndDelim; i++ {
		c.Core.Events.Subscribe(i, buff, "prune-c6-buff")
	}

}
