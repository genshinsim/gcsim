package emilie

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1ModKey      = "emilie-c1"
	c1ScentICDKey = "emilie-c1-attack-icd"
	c2ModKey      = "emilie-c2"
	c6ModKey      = "emilie-c6"

	c1ScentICD = 2.9 * 60
	c2Duration = 10 * 60
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	c.c1A1()

	c.Core.Events.Subscribe(event.OnBurning, func(args ...interface{}) bool {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		c.c1Scent()
		return false
	}, "emilie-a1-on-burning")

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if !ok {
			return false
		}
		if !t.IsBurning() {
			return false
		}
		if atk.Info.Element != attributes.Dendro {
			return false
		}
		c.c1Scent()
		return false
	}, "emilie-a1-on-damage")
}

func (c *char) c1A1() {
	if c.Base.Cons < 1 || c.Base.Ascension < 1 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.2
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c1ModKey, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.Abil != "Cleardew Cologne (A1)" {
				return nil, false
			}
			return m, true
		},
	})
}

func (c *char) c1Scent() {
	if c.StatusIsActive(c1ScentICDKey) {
		return
	}
	c.AddStatus(c1ScentICDKey, c1ScentICD, true)

	c.Core.Log.NewEvent("emilie c1 proc'd", glog.LogCharacterEvent, c.Index)
	c.generateScent()
}

func (c *char) c2(a combat.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	if a.Damage == 0 {
		return
	}

	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	e.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag(c2ModKey, c2Duration),
		Ele:   attributes.Dendro,
		Value: -0.3,
	})
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	c.AddStatus(c6ModKey, 5*60, true)
}

func (c *char) applyC6Bonus(ai *combat.AttackInfo) {
	if c.Base.Cons < 6 {
		return
	}
	if !c.StatusIsActive(c6ModKey) {
		return
	}

	switch ai.AttackTag {
	case attacks.AttackTagNormal, attacks.AttackTagExtra:
	default:
		return
	}
	ai.FlatDmg += c.TotalAtk() * 3
	ai.Element = attributes.Dendro
	ai.IgnoreInfusion = true
}
