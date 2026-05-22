package columbina

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var elevation = []float64{0.0, 0.015, 0.085, 0.1, 0.115, 0.13, 0.20}

const (
	elevationKey = "columbina-elevation"
	c1Key        = "columbina-c1"
	c1IcdKey     = "columbina-c1-icd"
	c2Key        = "columbina-c2"
	c2LCKey      = c2Key + "-lc"
	c2LCrKey     = c2Key + "-lcr"
	c2LBKey      = c2Key + "-lb"
	c4Key        = "columbina-c4"
	c4IcdKey     = "columbina-c4-icd"
	c6Key        = "columbina-c6"
	c6LCKey      = c6Key + "-lc"
	c6LBKey      = c6Key + "-lb"
	c6LCrKey     = c6Key + "-lcr"
)

func (c *char) consElevationInit() {
	if c.Base.Cons < 1 {
		return
	}
	amt := elevation[c.Base.Cons]
	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) {
		atk := args[0].(*info.AttackEvent)
		// do not apply elevation to Reaction damage here because the elevation is already applied at the contributor level
		if attacks.DirectLunarReactionStartDelim < atk.Info.AttackTag && atk.Info.AttackTag < attacks.DirectLunarReactionEndDelim {
			atk.Info.Elevation += amt
		}
	}, elevationKey+"-direct")

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if attacks.AttackTagIsLunar(atk.Info.AttackTag) {
			atk.Info.Elevation += amt
		}
	}, elevationKey+"-reaction")
}

func (c *char) c1OnSkill() {
	if c.Base.Cons < 1 {
		return
	}
	if c.StatusIsActive(c1IcdKey) {
		return
	}
	c.gravityTick(false)
	c.AddStatus(c1IcdKey, 15*60, true)
}

func (c *char) c1OnGravityTick(maxReaction lunarReaction) {
	if c.Base.Cons < 1 {
		return
	}
	switch maxReaction {
	case LunarCharge:
		c.Core.Player.ActiveChar().AddEnergy(c1Key, 6.0)
	case LunarCrystallize:
		// add shield
		c.Core.Player.Shields.Add(&shield.Tmpl{
			ActorIndex: c.Index(),
			Target:     c.Index(),
			Src:        c.Core.F,
			ShieldType: shield.ColumbinaC1,
			Name:       "Rainsea Shield",
			HP:         0.12 * c.MaxHP(),
			Ele:        attributes.Hydro,
			Expires:    c.Core.F + 8*60, // last until hitmark
		})
	}
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2Buff = make([]float64, attributes.EndStatType)
	c.c2Buff[attributes.HPP] = 0.4

	c.c2LCBuff = make([]float64, attributes.EndStatType)
	c.c2LBBuff = make([]float64, attributes.EndStatType)
	c.c2LCrBuff = make([]float64, attributes.EndStatType)
}

func (c *char) c2GravityRate() float64 {
	if c.Base.Cons < 2 {
		return 0.0
	}

	return 0.34
}

func (c *char) c2OnGravityTick(reaction lunarReaction) {
	if c.Base.Cons < 2 {
		return
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c2Key, 8*60),
		AffectedStat: attributes.HPP,
		Amount: func() []float64 {
			return c.c2Buff
		},
	})

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}

	// These buffs are NOT extra buffs, despite scaling off Columbina's Max HP.
	switch reaction {
	case LunarCharge:
		for _, char := range c.Core.Player.Chars() {
			char.DeleteStatMod(c2LCrKey)
			char.DeleteStatMod(c2LBKey)
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(c2LCKey, 8*60),
				AffectedStat: attributes.ATK,
				Amount: func() []float64 {
					if c.Core.Player.Active() != char.Index() {
						return nil
					}
					c.c2LCBuff[attributes.ATK] = 0.01 * c.MaxHP()
					return c.c2LCBuff
				},
			})
		}
	case LunarCrystallize:
		for _, char := range c.Core.Player.Chars() {
			char.DeleteStatMod(c2LCKey)
			char.DeleteStatMod(c2LBKey)
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(c2LCrKey, 8*60),
				AffectedStat: attributes.DEF,
				Amount: func() []float64 {
					if c.Core.Player.Active() != char.Index() {
						return nil
					}
					c.c2LCrBuff[attributes.DEF] = 0.01 * c.MaxHP()
					return c.c2LCrBuff
				},
			})
		}
	case LunarBloom:
		for _, char := range c.Core.Player.Chars() {
			char.DeleteStatMod(c2LCKey)
			char.DeleteStatMod(c2LCrKey)
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(c2LBKey, 8*60),
				AffectedStat: attributes.EM,
				Amount: func() []float64 {
					if c.Core.Player.Active() != char.Index() {
						return nil
					}
					c.c2LBBuff[attributes.EM] = 0.0035 * c.MaxHP()
					return c.c2LBBuff
				},
			})
		}
	default:
		panic("unreachable: columbina c2 called with unknown reaction")
	}
}

func (c *char) c4OnGravityTickFlatDMG(reaction lunarReaction) float64 {
	if c.Base.Cons < 4 {
		return 0.0
	}

	c.AddEnergy(c4Key, 4)
	if c.StatusIsActive(c4IcdKey) {
		return 0.0
	}
	c.AddStatus(c4IcdKey, 15*60, true)
	switch reaction {
	case LunarCharge, LunarCrystallize:
		return 0.125 * c.MaxHP()
	case LunarBloom:
		return 0.025 * c.MaxHP()
	default:
		panic("unreachable: columbina c4 called with unknown reaction")
	}
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c.c6Buff = make([]float64, attributes.EndStatType)
	c.c6Buff[attributes.CD] = 0.8

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		ae, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)

		char := c.Core.Player.Chars()[atk.Info.ActorIndex]
		addBuff := false
		switch atk.Info.Element {
		case attributes.Electro:
			addBuff = char.StatusIsActive(c6LCKey)
		case attributes.Hydro:
			addBuff = char.StatusIsActive(c6LCKey) || char.StatusIsActive(c6LCrKey) || char.StatusIsActive(c6LBKey)
		case attributes.Geo:
			addBuff = char.StatusIsActive(c6LCrKey)
		case attributes.Dendro:
			addBuff = char.StatusIsActive(c6LBKey)
		}

		if !addBuff {
			return
		}

		ae.Snapshot.Stats[attributes.CD] += 0.8
	}, c6Key+"-reaction-attack")

	c.Core.Events.Subscribe(event.OnLunarCharged, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if !c.ReactBonusModIsActive(burstBuffKey) {
			return
		}

		for _, char := range c.Core.Player.Chars() {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag(c6LCKey, 8*60),
				Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
					switch atk.Info.Element {
					case attributes.Electro:
					case attributes.Hydro:
					default:
						return nil
					}
					return c.c6Buff
				},
			})
		}
	}, c6LCKey)
	c.Core.Events.Subscribe(event.OnLunarBloom, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if !c.ReactBonusModIsActive(burstBuffKey) {
			return
		}

		for _, char := range c.Core.Player.Chars() {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag(c6LBKey, 8*60),
				Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
					switch atk.Info.Element {
					case attributes.Dendro:
					case attributes.Hydro:
					default:
						return nil
					}
					return c.c6Buff
				},
			})
		}
	}, c6LBKey)

	c.Core.Events.Subscribe(event.OnLunarCrystallize, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if !c.ReactBonusModIsActive(burstBuffKey) {
			return
		}

		for _, char := range c.Core.Player.Chars() {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag(c6LCrKey, 8*60),
				Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
					switch atk.Info.Element {
					case attributes.Geo:
					case attributes.Hydro:
					default:
						return nil
					}
					return c.c6Buff
				},
			})
		}
	}, c6LCrKey)
}
