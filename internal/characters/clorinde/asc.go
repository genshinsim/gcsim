package clorinde

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	clorindeA1BuffKey      = `clorinde-a1-buff`
	clordineA1BuffDuration = int(a1Duration * 60)
	clorindeA4BuffKey      = `clorinde-a4-buff`
	clordineA4BuffDuration = int(a4Duration * 60)
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.a1BuffPercent = a1PercentBuff
	c.a1Cap = a1FlatDmg
	if c.Base.Cons >= 2 {
		c.a1BuffPercent = c2A1PercentBuff
		c.a1Cap = c2A1FlatDmg
	}
	c.a1stacks = stacks.NewMultipleRefreshNoRemove(3, c.QueueCharTask, &c.Core.F)
	// on electro reaction, add buff; 3 stacks independent
	c.Core.Events.Subscribe(event.OnElectroCharged, c.a1CB, "clorinde-a1-ec")
	c.Core.Events.Subscribe(event.OnSuperconduct, c.a1CB, "clorinde-a1-superconduct")
	c.Core.Events.Subscribe(event.OnAggravate, c.a1CB, "clorinde-a1-aggravate")
	c.Core.Events.Subscribe(event.OnQuicken, c.a1CB, "clorinde-a1-quicken")
	c.Core.Events.Subscribe(event.OnHyperbloom, c.a1CB, "clorinde-a1-hyperbloom")
	c.Core.Events.Subscribe(event.OnOverload, c.a1CB, "clorinde-a1-overload")
	c.Core.Events.Subscribe(event.OnSwirlElectro, c.a1CB, "clorinde-a1-swirl-electro")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, c.a1CB, "clorinde-a1-crystallize-electro")
}

func (c *char) a1CB(args ...interface{}) bool {
	// no requirement who triggers other than that it must be against an enemy
	if _, ok := args[0].(*enemy.Enemy); !ok {
		return false
	}
	// add a stack and refresh the mod for 15s
	c.a1stacks.Add(clordineA1BuffDuration)
	c.AddAttackMod(character.AttackMod{
		Base:   modifier.NewBaseWithHitlag(clorindeA1BuffKey, clordineA1BuffDuration),
		Amount: c.a1Amount,
	})

	return false
}

func (c *char) a1Amount(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
	var amt float64
	switch atk.Info.AttackTag {
	case attacks.AttackTagNormal:
		if atk.Info.Element != attributes.Electro {
			// only app
			return nil, false
		}
	case attacks.AttackTagElementalBurst:
	default:
		return nil, false
	}
	totalAtk := atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[attributes.ATKP]) + atk.Snapshot.Stats[attributes.ATK]
	amt = min(totalAtk*c.a1BuffPercent*float64(c.a1stacks.Count()), c.a1Cap)
	atk.Info.FlatDmg += amt
	c.Core.Log.NewEvent("a1 adding flat dmg", glog.LogCharacterEvent, c.Index).
		Write("amt", amt).
		Write("c2_applied", c.Base.Cons >= 2)
	// we don't actually change any stats here..
	return nil, true
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4stacks = stacks.NewMultipleRefreshNoRemove(2, c.QueueCharTask, &c.Core.F)
	c.a4bonus = make([]float64, attributes.EndStatType)
}

func (c *char) a4(change float64) {
	// if BOL > 100%, then on BOL change, add 15s buff; 2 stacks each tracked independently
	if c.Base.Ascension < 4 {
		return
	}
	if c.currentHPDebtRatio() < 1 {
		return
	}
	if math.Abs(change) < tolerance {
		return
	}
	c.a4stacks.Add(clordineA4BuffDuration)
	c.AddStatMod(character.StatMod{
		Base:   modifier.NewBaseWithHitlag(clorindeA4BuffKey, clordineA4BuffDuration),
		Amount: c.a4Amount,
	})
	c.Core.Log.NewEvent("a4 triggered", glog.LogCharacterEvent, c.Index).
		Write("stacks", c.a4stacks.Count())
}

func (c *char) a4Amount() ([]float64, bool) {
	c.a4bonus[attributes.CR] = float64(c.a4stacks.Count()) * a4CritBuff
	return c.a4bonus, true
}
