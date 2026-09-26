package zibai

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	a1Key         = "zibai-a1"
	a4Key         = "zibai-a4"
	lunarBonusKey = "zibai-lcr-bonus"

	a4Mult = 0.6
)

// When casting the Elemental Skill Heaven and Earth Made Manifest, or when nearby party members
// trigger Moondrift Harmony, Zibai gains the Selenic Descent effect for 4s: The DMG dealt by the
// 2nd hit of Spirit Steed's Stride is increased by 60% of Zibai's DEF.
func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnMoondriftHarmony, func(args ...any) {
		c.AddStatus(a1Key, 4*60, true)
	}, "zibai-a1")
}

func (c *char) a1OnSkill() {
	if c.Base.Ascension < 1 {
		return
	}

	c.AddStatus(a1Key, 4*60, true)
}

func (c *char) a1StrideBonusDmg() float64 {
	if c.Base.Ascension < 1 {
		return 0.0
	}
	if !c.StatusIsActive(a1Key) {
		return 0.0
	}

	return (0.6 + c.c2A4Mult()) * c.TotalDef(false)
}

// Other Geo party members increase Zibai's DEF by 15% each. Hydro party members increase her
// Elemental Mastery by 60 each.
func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	hydros := 0
	geos := 0
	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		switch char.Base.Element {
		case attributes.Hydro:
			hydros += 1
		case attributes.Geo:
			geos += 1
		}
	}
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBase(a4Key, -1),
		Amount: func() []float64 {
			m[attributes.DEFP] = 0.15 * float64(geos)
			m[attributes.EM] = 60.0 * float64(hydros)
			return m
		},
	})
}

// When a party member triggers a Hydro Crystallize reaction, it will be converted into the
// Lunar-Crystallize reaction, with every 100 DEF that Zibai has increasing Lunar-Crystallize's
// Base DMG by 0.7%, up to a maximum of 14%.
func (c *char) moonsignInit() {
	c.Core.Flags.Custom[reactable.LunarCrystallizeEnableKey] = 1
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagDirectLunarCrystallize {
			return
		}

		bonus := min(c.TotalDef(false)/100.0*0.007, 0.14)

		c.Core.Log.NewEvent("zibai adding lunar crystallize base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)

		atk.Info.BaseDmgBonus += bonus
	}, lunarBonusKey)

	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagReactionLunarCrystallize {
			return
		}

		bonus := min(c.TotalDef(false)/100.0*0.007, 0.14)

		c.Core.Log.NewEvent("zibai adding lunar crystallize base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)

		atk.Info.BaseDmgBonus += bonus
	}, lunarBonusKey+"-reaction")
}
