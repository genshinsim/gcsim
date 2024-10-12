package xilonen

import (
	"fmt"
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2key = "xilonen-c2"
const c4key = "xilonen-c4"
const c6key = "xilonen-c6"
const c6IcdKey = "xilonen-c6-icd"
const c6StamKey = "xilonen-c6-stam"

func (c *char) c1DurMod() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1.45
}

func (c *char) c1ValMod() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1 / 1.45
}

var c2BuffGeo []float64
var c2BuffPyro []float64
var c2BuffHydro []float64
var c2BuffCryo []float64

func c2buffsInit() {
	c2BuffGeo = make([]float64, attributes.EndStatType)
	c2BuffGeo[attributes.DmgP] = 0.50

	c2BuffPyro = make([]float64, attributes.EndStatType)
	c2BuffPyro[attributes.ATKP] = 0.45

	c2BuffHydro = make([]float64, attributes.EndStatType)
	c2BuffHydro[attributes.HPP] = 0.45

	c2BuffCryo = make([]float64, attributes.EndStatType)
	c2BuffCryo[attributes.CD] = 0.60
}

func (c *char) c2buff() {
	for _, ch := range c.Core.Player.Chars() {
		switch ch.Base.Element {
		case attributes.Geo:
			ch.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.DmgP,
				Amount: func() ([]float64, bool) {
					// geo is always active
					return c2BuffGeo, true
				},
			})
		case attributes.Pyro:
			ch.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					if c.StatusIsActive(activeSamplerKey) {
						return c2BuffPyro, true
					}
					return nil, false
				},
			})
		case attributes.Hydro:
			ch.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.HPP,
				Amount: func() ([]float64, bool) {
					if c.StatusIsActive(activeSamplerKey) {
						return c2BuffHydro, true
					}
					return nil, false
				},
			})
		case attributes.Cryo:
			ch.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.CD,
				Amount: func() ([]float64, bool) {
					if c.StatusIsActive(activeSamplerKey) {
						return c2BuffCryo, true
					}
					return nil, false
				},
			})
		}
	}
}

func (c *char) c2GeoSampler() func() {
	return func() {
		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), nil)
		c.applySamplerShred(attributes.Geo, enemies)

		// TODO: how often does this apply?
		c.QueueCharTask(c.c2GeoSampler(), 30)
	}
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	if slices.Contains(c.shredElements, attributes.Geo) {
		c.c2GeoSampler()()
	}
	c.c2buff()
}

func (c *char) c2electro() {
	if c.Base.Cons < 2 {
		return
	}
	for _, ch := range c.Core.Player.Chars() {
		if ch.Base.Element == attributes.Electro {
			ch.AddEnergy(c2key, 25)
			ch.ReduceActionCooldown(action.ActionBurst, 6*60)
		}
	}
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.AddStatus(c4key, 15*60, true) // 15 sec duration
		char.SetTag(c4key, 6)              // 6 c4 stacks
	}
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}

		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)

		if !char.StatusIsActive(c4key) {
			return false
		}

		if char.Tags[c4key] > 0 {
			amt := 0.65 * c.TotalDef()
			char.Tags[c4key]--

			c.Core.Log.NewEvent("Xilonen c4 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("effect_ends_at", c.StatusExpiry(c4key)).
				Write("c4_left", c.Tags[c4key])

			atk.Info.FlatDmg += amt
		}

		return false
	}, fmt.Sprintf("%s-hook", c4key))
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	if c.StatusIsActive(c6IcdKey) {
		return
	}

	if !c.nightsoulState.HasBlessing() {
		return
	}

	c.AddStatus(c6key, 5*60, true)
	c.AddStatus(c6IcdKey, 15*60, true)
	c.c6activated = true

	src := c.nightsoulSrc
	cancelTime := c.StatusDuration(skillMaxDurKey) + 5*60
	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}
		c.exitNightsoul()
	}, cancelTime)
	c.AddStatus(skillMaxDurKey, cancelTime, true)

	c.QueueCharTask(func() {
		if !c.nightsoulState.HasBlessing() {
			return
		}
		if c.nightsoulState.Points() < maxNightsoulPoints {
			return
		}
	}, 5*60)

	c.QueueCharTask(c.nightsoulPointReduceFunc(c.nightsoulSrc), 5*60)

	for i := 0; i < 4; i++ {
		hpplus := c.Stat(attributes.Heal)
		heal := c.TotalDef() * 1.2
		c.Core.Tasks.Add(func() {
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Xilonen C6 Healing",
				Src:     heal,
				Bonus:   hpplus,
			})
		}, i*90)
	}
}

func (c *char) c6Stam() {
	if c.Base.Cons < 6 {
		return
	}
	c.Core.Player.AddStamPercentMod(c6StamKey, -1, func(a action.Action) (float64, bool) {
		if c.StatusIsActive(c6key) {
			return -1, false
		}
		return 0, false
	})
}

func (c *char) c6DmgMult() float64 {
	if c.Base.Cons < 6 {
		return 0.0
	}
	if !c.StatusIsActive(c6key) {
		return 0.0
	}
	return 3.0
}
