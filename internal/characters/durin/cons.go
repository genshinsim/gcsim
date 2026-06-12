package durin

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1Key = "durin-c1"

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalBurst:
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return
		}

		if atk.Info.ActorIndex == c.Index() {
			return
		}

		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)

		if char.Tags[c1Key] <= 0 {
			return
		}

		if !c.StatusIsActive(burstKeyWhite) {
			return
		}

		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		mult := 0.6
		consume := 1

		consume *= c.c4c1ConsumeMult()
		char.Tags[c1Key] -= consume

		amt := mult * c.TotalAtk()

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Durin C1 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("effect_ends_at", c.StatusExpiry(c1Key)).
				Write("stacks_left", char.Tags[c1Key])
		}

		atk.Info.FlatDmg += amt
	}, "durin-c1-white-hook")

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if c.Tags[c1Key] == 0 {
			return
		}

		mult := 1.5
		consume := 2

		if !c.StatusIsActive(burstKeyBlack) {
			return
		}

		consume *= c.c4c1ConsumeMult()
		c.Tags[c1Key] -= consume

		amt := mult * c.TotalAtk()

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Durin C1 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("effect_ends_at", c.StatusExpiry(c1Key)).
				Write("stacks_left", c.Tags[c1Key])
		}

		atk.Info.FlatDmg += amt
	}, "durin-c1-black-hook")
}

func (c *char) c1OnBurst(isWhite bool) {
	if c.Base.Cons < 1 {
		return
	}

	if isWhite {
		c.RemoveTag(c1Key)

		for _, char := range c.Core.Player.Chars() {
			if c.Index() == char.Index() {
				continue
			}
			char.SetTag(c1Key, 20)
		}
		return
	}

	c.SetTag(c1Key, 20)
	for _, char := range c.Core.Player.Chars() {
		if c.Index() == char.Index() {
			continue
		}
		char.RemoveTag(c1Key)
	}
}

var c2ReactToElements = map[event.Event][]attributes.Element{
	event.OnOverload:        {attributes.Electro, attributes.Pyro},
	event.OnSwirlPyro:       {attributes.Anemo, attributes.Pyro},
	event.OnCrystallizePyro: {attributes.Geo, attributes.Pyro},
	event.OnBurning:         {attributes.Dendro, attributes.Pyro},
	event.OnVaporize:        {attributes.Hydro, attributes.Pyro},
	event.OnMelt:            {attributes.Cryo, attributes.Pyro},
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2Buff = make([]float64, attributes.EndStatType)

	makeBuff := func(elements []attributes.Element) func(args ...any) {
		return func(args ...any) {
			_, ok := args[0].(*enemy.Enemy)
			if !ok {
				return
			}
			c.c2MakeBuff(elements)
		}
	}

	for event, elements := range c2ReactToElements {
		c.Core.Events.Subscribe(event, makeBuff(elements), fmt.Sprintf("durin-c2-hook-%v", event))
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		t, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if !t.IsBurning() {
			return
		}
		switch atk.Info.Element {
		case attributes.Dendro:
		case attributes.Pyro:
		default:
			return
		}

		if !c.StatusIsActive(burstKeyWhite) && !c.StatusIsActive(burstKeyBlack) {
			return
		}

		c.c2MakeBuff([]attributes.Element{attributes.Pyro, attributes.Dendro})
	}, "durin-c2-hook-on-dmg")
}

func (c *char) c2MakeBuff(elements []attributes.Element) {
	for _, elem := range elements {
		for _, char := range c.Core.Player.Chars() {
			stat := attributes.EleToDmgP(elem)
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("durin-c2-"+elem.String(), 6*60),
				AffectedStat: stat,
				Amount: func() []float64 {
					for i := range c.c2Buff {
						c.c2Buff[i] = 0
					}
					c.c2Buff[stat] = 0.5
					return c.c2Buff
				},
			})
		}
	}
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.4

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("durin-c4", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil
			}
			return m
		},
	})
}

func (c *char) c4c1ConsumeMult() int {
	if c.Base.Cons < 4 {
		return 1
	}
	if c.Core.Rand.Float64() < 0.3 {
		return 0
	}

	return 1
}

func (c *char) c6DefIgnore(isWhite bool) float64 {
	if c.Base.Cons < 6 {
		return 0
	}

	if isWhite {
		return 0.3
	}

	return 0.7
}

func (c *char) c6WhiteMakeCB() func(a info.AttackCB) {
	if c.Base.Cons < 6 {
		return nil
	}

	return func(a info.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddDefMod(info.DefMod{
			Base:  modifier.NewBaseWithHitlag("durin-c6", 6*60),
			Value: -0.3,
		})
	}
}
