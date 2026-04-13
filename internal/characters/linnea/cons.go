package linnea

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	c1Key     = "linnea-c1"
	c2Key     = "linnea-c2"
	c4Key     = "linnea-c4"
	c4KeySelf = "linnea-c4-self"
	c6Key     = "linnea-c6"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnMoondriftHarmony, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		if !c.StatusIsActive(c1Key) {
			c.c1Stacks = 0
		}
		c.AddStatus(c1Key, 10*60, true)
		stacks := c.c6C1Stacks()
		c.c1Stacks = min(c.c1Stacks+stacks, 18)
	}, c1Key)

	f := func(atk *info.AttackEvent) {
		if !c.StatusIsActive(c1Key) {
			return
		}

		maxStacks := 1
		scaling := 0.75
		if atk.Info.ActorIndex == c.Index() && atk.Info.Abil == skillMillionAbil {
			maxStacks = 5
			scaling = 1.5
		}

		c6stacks, c6scale := c.c6C1Mult()
		maxStacks *= c6stacks
		scaling *= c6scale

		if c.c1Stacks > 0 {
			def := c.TotalDef(false)
			stacks := min(c.c1Stacks, maxStacks)
			amt := def * scaling * float64(stacks)
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Linnea C1 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("addition", amt).
					Write("Field Catalog stacks left", c.c1Stacks)
			}
			atk.Info.FlatDmg += amt
			c.c1Stacks -= stacks
		}
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCrystallize:
		default:
			return
		}
		f(atk)
	}, c1Key)
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagReactionLunarCrystallize:
		default:
			return
		}
		f(atk)
	}, c1Key)
}

func (c *char) c1OnSkill() {
	if c.Base.Cons < 1 {
		return
	}
	if !c.StatusIsActive(c1Key) {
		c.c1Stacks = 0
	}
	c.AddStatus(c1Key, 10*60, true)
	c.c1Stacks = min(c.c1Stacks+6, 18)
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.CD] = 0.4
	c.Core.Events.Subscribe(event.OnMoondriftHarmony, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		if !c.StatusIsActive(c1Key) {
			c.c1Stacks = 0
		}
		for _, char := range c.Core.Player.Chars() {
			switch char.Base.Element {
			case attributes.Geo:
			case attributes.Hydro:
			default:
				return
			}
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2Key, 8*60),
				AffectedStat: attributes.CD,
				Amount: func() []float64 {
					return m
				},
			})
		}
	}, c2Key)
}

func (c *char) c2MillionTonCDBonus(snap *info.Snapshot) {
	if c.Base.Cons < 2 {
		return
	}
	snap.Stats[attributes.CD] += 1.5
}

func (c *char) c2TriggerMoonDrift() {
	if c.Base.Cons < 2 {
		return
	}
	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}
	var contribMap [info.MaxChars]bool
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Hydro:
		case attributes.Geo:
		default:
			continue
		}
		contribMap[char.Index()] = true
	}
	reactable.DoLCrAttackWithContrib(contribMap, c.Core.Combat.PrimaryTarget(), c.Core, c.Index())
	c.Core.Log.NewEvent("Linnea C2 Lunar Crystallize attack triggered", glog.LogElementEvent, c.Index())
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DEFP] = 0.25
	c.Core.Events.Subscribe(event.OnMoondriftHarmony, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c4Key, 5*60),
				AffectedStat: attributes.DEFP,
				Amount: func() []float64 {
					if c.Core.Player.Active() == char.Index() {
						return m
					}
					return nil
				},
			})
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c4KeySelf, 5*60),
			AffectedStat: attributes.DEFP,
			Amount: func() []float64 {
				return m
			},
		})
	}, c4Key)
}

func (c *char) c6C1Stacks() int {
	if c.Base.Cons < 6 {
		return 6
	}
	return 18
}

func (c *char) c6C1Mult() (int, float64) {
	if c.Base.Cons < 6 {
		return 1, 1.0
	}
	return 2, 1.5
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}

	amt := 0.25
	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) {
		atk := args[0].(*info.AttackEvent)
		if atk.Info.AttackTag == attacks.AttackTagDirectLunarCrystallize {
			atk.Info.Elevation += amt
		}
	}, c6Key+"-direct")

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag == attacks.AttackTagReactionLunarCrystallize {
			atk.Info.Elevation += amt
		}
	}, c6Key+"-reaction")
}
