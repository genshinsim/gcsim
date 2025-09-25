package skirk

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	a1Dur    = 1054
	a1Key    = "skirk-a1"
	a1IcdKey = "skirk-a1-icd"
	a4Key    = "deaths-crossing"
	a4Dur    = 20 * 60
)

var (
	a4MultAttack = []float64{1, 1.1, 1.2, 1.7}
	a4MultBurst  = []float64{1, 1.05, 1.15, 1.60}
)

func (c *char) a1Init() {
	a1Hook := func(args ...any) bool {
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
	c.voidRifts = NewRingQueue[int](3)
	c.Core.Events.Subscribe(event.OnFrozen, a1Hook, a1Key+"frozen")
	c.Core.Events.Subscribe(event.OnSuperconduct, a1Hook, a1Key+"superconduct")
	c.Core.Events.Subscribe(event.OnSwirlCryo, a1Hook, a1Key+"cryo-swirl")
	c.Core.Events.Subscribe(event.OnCrystallizeCryo, a1Hook, a1Key+"cryo-crystallize")
}

func (c *char) absorbVoidRiftCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	c.absorbVoidRifts()
}

func (c *char) absorbVoidRifts() int {
	filter := func(src int) bool {
		return src+a1Dur >= c.Core.F
	}
	count := c.voidRifts.Count(filter)
	c.voidRifts.Clear()

	c.onVoidAbsorb(count)
	return count
}

func (c *char) onVoidAbsorb(count int) {
	if count <= 0 {
		return
	}

	c.AddSerpentsSubtlety("a1-void-rifts", float64(count)*8.0)

	for range count {
		c.c1()
		c.c6OnVoidAbsorb()
	}
}

func (c *char) createVoidRift() {
	// absorb the rift immediately if currently in the hE/ or E-Burst animation
	if c.StatusIsActive(skillAbsorbRiftAnimKey) {
		c.onVoidAbsorb(1)
	}
	if c.StatusIsActive(burstAbsorbRiftAnimKey) {
		c.onVoidAbsorb(1)
		c.burstVoids = min(c.burstVoids+1, 3)
	}
	c.voidRifts.PushOverwrite(c.Core.F)
}

func (c *char) a4Init() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		charElem := c.Core.Player.Chars()[atk.Info.ActorIndex].Base.Element
		if atk.Info.ActorIndex == c.Index() {
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
		c.AddStatus(getA4StackName(atk.Info.ActorIndex), a4Dur, true)

		return false
	}, a4Key+"-hook")
}

func getA4StackName(index int) string {
	return fmt.Sprintf("%s-char%d", a4Key, index)
}

func (c *char) getA4Stacks() int {
	count := 0
	for index := range c.Core.Player.Chars() {
		if c.StatusIsActive(getA4StackName(index)) {
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
