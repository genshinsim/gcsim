package arlecchino

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var a1Directive = []float64{0.0, 0.65, 1.3}

func (c *char) passive() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.4
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("arlecchino-passive", -1),
		AffectedStat: attributes.PyroP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

func (c *char) a1OnKill() {
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		e, ok := args[0].(*enemy.Enemy)
		// ignore if not an enemy
		if !ok {
			return false
		}

		if !e.StatusIsActive(directiveKey) {
			return false
		}
		// always max level debt
		newDebt := a1Directive[len(a1Directive)-1] * c.MaxHP()
		if c.StatusIsActive(directiveLimitKey) {
			newDebt = min(c.skillDebtMax-c.skillDebt, newDebt)
		}

		if newDebt > 0 {
			c.skillDebt += newDebt
			c.ModifyHPDebtByAmount(newDebt)
		}
		e.RemoveTag(directiveKey)
		e.RemoveTag(directiveSrcKey)
		e.DeleteStatus(directiveKey)
		return false
	}, "arlechinno-a1-onkill")
}

func (c *char) a1Upgrade(e combat.Enemy, src int) {
	if c.Base.Ascension < 1 {
		return
	}
	e.QueueEnemyTask(func() {
		level := e.GetTag(directiveKey)
		if level == 0 {
			return
		}
		if level >= 2 {
			return
		}
		if e.GetTag(directiveSrcKey) != src {
			return
		}
		e.SetTag(directiveKey, level+1)
		c.Core.Log.NewEvent("Directive upgraded", glog.LogCharacterEvent, c.Index).
			Write("new_level", level+1).
			Write("src", src)
	}, 5*60)
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	// Resistances are not implemented
}
