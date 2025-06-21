package skirk

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const a1Key = "skirk-a1"
const a1IcdKey = "skirk-a1-icd"
const a4Key = "deaths-crossing"
const a4Dur = 20 * 60

var a4MultAttack = []float64{1, 1.1, 1.2, 1.7}
var a4MultBurst = []float64{1, 1.05, 1.15, 1.60}

func (c *char) a1Init() {
	a1Hook := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		if c.StatusIsActive(a1IcdKey) {
			return false
		}
		c.AddStatus(a1IcdKey, 2.5*60, true)
		c.createVoidRift()

		return false
	}
	c.Core.Events.Subscribe(event.OnFrozen, a1Hook, a1Key+"frozen")
	c.Core.Events.Subscribe(event.OnSuperconduct, a1Hook, a1Key+"superconduct")
	c.Core.Events.Subscribe(event.OnSwirlCryo, a1Hook, a1Key+"cryo-swirl")
	c.Core.Events.Subscribe(event.OnCrystallizeCryo, a1Hook, a1Key+"cryo-crystallize")
}

func (c *char) absorbVoidRiftCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	c.absorbVoidRift()
}

func (c *char) absorbVoidRift() {
	count := c.voidRiftCount
	if count > 3 {
		count = 3
	}
	if count <= 0 {
		return
	}
	c.AddSerpentsSubtlety("a1-void-rifts", float64(count)*8.0)

	for i := 0; i < count; i++ {
		c.c1()
		c.c6OnVoidAbsorb()
	}
	c.voidRiftCount = 0
}

func (c *char) createVoidRift() {
	c.voidRiftCount++
	// TODO: when do these time out?
}

func (c *char) a4Init() {
	c.a4Stacks = make([]int, len(c.Core.Player.Chars()))
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		charElem := c.Core.Player.Chars()[atk.Info.ActorIndex].Base.Element
		if atk.Info.ActorIndex == c.Index {
			return false
		}
		if atk.Info.Element != charElem {
			return false
		}
		switch charElem {
		case attributes.Cryo:
		case attributes.Hydro:
		default:
			return false
		}
		c.a4Stacks[atk.Info.ActorIndex] = c.TimePassed

		return false
	}, a4Key+"-hook")
}

func (c *char) getA4Stacks() int {
	count := 0
	for _, f := range c.a4Stacks {
		if f != 0 && f+a4Dur > c.TimePassed {
			count++
		}
	}
	return count
}

func (c *char) a4MultAttack() float64 {
	return a4MultAttack[c.getA4Stacks()]
}
func (c *char) a4MultBurst() float64 {
	return a4MultBurst[c.getA4Stacks()]
}
