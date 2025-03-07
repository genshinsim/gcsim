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

const (
	c2key      = "xilonen-c2"
	c2BuffKey  = "xilonen-c2-buff"
	c4key      = "xilonen-c4"
	c6key      = "xilonen-c6"
	c6IcdKey   = "xilonen-c6-icd"
	c6StamKey  = "xilonen-c6-stam"
	c2Interval = 0.3 * 60
	c6Duration = 5 * 60
)

var c2Buffs = map[attributes.Element][]float64{
	attributes.Geo:   make([]float64, attributes.EndStatType),
	attributes.Pyro:  make([]float64, attributes.EndStatType),
	attributes.Hydro: make([]float64, attributes.EndStatType),
	attributes.Cryo:  make([]float64, attributes.EndStatType),
}

func init() {
	c2Buffs[attributes.Geo][attributes.DmgP] = 0.5
	c2Buffs[attributes.Pyro][attributes.ATKP] = 0.45
	c2Buffs[attributes.Hydro][attributes.HPP] = 0.45
	c2Buffs[attributes.Cryo][attributes.CD] = 0.60
}

func (c *char) nightsoulDurationMul() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1.0 + 0.45
}

func (c *char) nightsoulConsumptionMul() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1.0 - 0.3
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	if c.samplersConverted >= 3 {
		return
	}

	c.activeGeoSampler(-1)()
	for _, ch := range c.Core.Player.Chars() {
		if ch.Base.Element != attributes.Geo {
			continue
		}
		ch.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c2BuffKey, -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				return c2Buffs[attributes.Geo], true
			},
		})
	}
}

func (c *char) applyC2Buff(src int, other *character.CharWrapper) func() {
	return func() {
		if c.c2Src != src {
			return
		}
		if !c.StatusIsActive(c2key) {
			return
		}
		other.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(c2BuffKey, 60),
			Amount: func() ([]float64, bool) {
				return c2Buffs[other.Base.Element], true
			},
		})
		c.QueueCharTask(c.applyC2Buff(src, other), c2Interval)
	}
}

func (c *char) c2activate() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2Src = c.Core.F
	for _, other := range c.Core.Player.Chars() {
		switch other.Base.Element {
		case attributes.Geo:
			// skip, because it's already applied
		case attributes.Pyro, attributes.Hydro, attributes.Cryo:
			c.QueueCharTask(c.applyC2Buff(c.c2Src, other), c2Interval)
		case attributes.Electro:
			other.AddEnergy(c2key, 25)
			other.ReduceActionCooldown(action.ActionBurst, 6*60)
		default:
			continue
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
		if !char.StatusIsActive(c4key) || char.Tag(c4key) == 0 {
			return false
		}

		amt := 0.65 * c.TotalDef(false)
		char.SetTag(c4key, char.Tag(c4key)-1)

		c.Core.Log.NewEvent("xilonen c4 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
			Write("before", atk.Info.FlatDmg).
			Write("addition", amt).
			Write("effect_ends_at", c.StatusExpiry(c4key)).
			Write("c4_left", char.Tag(c4key))

		atk.Info.FlatDmg += amt
		return false
	}, fmt.Sprintf("%s-hook", c4key))
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	onAction := func(...interface{}) bool {
		if c.Core.Player.Active() == c.Index && c.nightsoulState.HasBlessing() {
			c.applyC6()
		}
		return false
	}

	c.Core.Events.Subscribe(event.OnAttack, onAction, "xilonen-c6-on-attack")
	c.Core.Events.Subscribe(event.OnDash, onAction, "xilonen-c6-on-dash")
	c.Core.Events.Subscribe(event.OnPlunge, onAction, "xilonen-c6-on-plunge")
}

func (c *char) applyC6() {
	if c.StatusIsActive(c6IcdKey) {
		return
	}
	c.AddStatus(c6IcdKey, 15*60, true)
	c.c6FlatDmg() // sets c6 key

	// "pause" Nightsoul's Blessing time limit countdown
	duration := c.StatusDuration(skillMaxDurKey) + c6Duration
	c.setNightsoulExitTimer(duration)

	for i := 1; i <= 3; i++ {
		c.Core.Tasks.Add(func() {
			if !c.StatusIsActive(c6key) {
				return
			}
			hpplus := c.Stat(attributes.Heal)
			heal := c.TotalDef(false) * 1.2
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Imperishable Night Carnival",
				Src:     heal,
				Bonus:   hpplus,
			})
		}, i*90)
	}
}

func (c *char) c6FlatDmg() {
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(c6key, c6Duration),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagNormal, attacks.AttackTagPlunge:
			default:
				return nil, false
			}
			if !slices.Contains(atk.Info.AdditionalTags, attacks.AdditionalTagNightsoul) {
				return nil, false
			}

			amt := c.TotalDef(false) * 3.0
			c.Core.Log.NewEvent("c6 proc dmg add", glog.LogPreDamageMod, c.Index).
				Write("amt", amt)

			atk.Info.FlatDmg += amt
			return nil, true
		},
	})
}
