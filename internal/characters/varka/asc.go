package varka

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var a1Multiplier = []float64{1.0, 1.4, 2.2}

const (
	a4Key    = "varka-a4-stacks"
	a4ICDKey = "varka-a4-icd"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	dmgp := attributes.EleToDmgP(c.conversionElem)

	m := make([]float64, attributes.EndStatType)

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("varka-a1-anemo", -1),
		AffectedStat: attributes.AnemoP,
		Amount: func() []float64 {
			m[dmgp] = 0
			m[attributes.AnemoP] = c.a1Buff
			return m
		},
	})

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("varka-a1-"+c.conversionElem.String(), -1),
		AffectedStat: dmgp,
		Amount: func() []float64 {
			m[dmgp] = c.a1Buff
			m[attributes.AnemoP] = 0
			return m
		},
	})

	// the buff updates every 1s
	c.a1UpdateBuff()

	elementCount := make([]int, attributes.EndEleType)
	for _, char := range c.Core.Player.Chars() {
		elementCount[char.Base.Element] += 1
	}

	twoAnemo := elementCount[attributes.Anemo] >= 2
	twoPHEC := elementCount[attributes.Pyro] >= 2 || elementCount[attributes.Hydro] >= 2 || elementCount[attributes.Electro] >= 2 || elementCount[attributes.Cryo] >= 2

	if twoAnemo && twoPHEC {
		c.a1Multiplier = a1Multiplier[2]
	} else if twoAnemo || twoPHEC {
		c.a1Multiplier = a1Multiplier[1]
	} else {
		c.a1Multiplier = a1Multiplier[0]
	}
}

func (c *char) a1UpdateBuff() {
	stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
	c.a1Buff = min(stats.TotalATK()*0.1/1000, 0.25)
	c.QueueCharTask(c.a1UpdateBuff, 60)
}

func (c *char) a1SkillMulti() float64 {
	if c.Base.Ascension < 1 {
		return 1
	}
	return c.a1Multiplier
}

func (c *char) a4Init() {
	a4Hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		char := c.Core.Player.Chars()[atk.Info.ActorIndex]
		if char.StatusIsActive(a4ICDKey) {
			return
		}
		char.AddStatus(a4ICDKey, 1*60, true)

		if c.StatusIsActive(a4Key) {
			c.a4Stacks = min(c.a4Stacks+1, 4)
		} else {
			c.a4Stacks = 1
		}
		c.AddStatus(a4Key, 8*60, true)
	}

	c.Core.Events.Subscribe(event.OnSwirlPyro, a4Hook, "varka-a4-pyro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, a4Hook, "varka-a4-hydro")
	c.Core.Events.Subscribe(event.OnSwirlElectro, a4Hook, "varka-a4-electro")
	c.Core.Events.Subscribe(event.OnSwirlCryo, a4Hook, "varka-a4-cryo")

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("varka-a4", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra && !slices.Contains(atk.Info.AdditionalTags, attacks.AdditionalTagVarkaSpecial) {
				return nil
			}
			if !c.StatusIsActive(a4Key) {
				return nil
			}
			m[attributes.DmgP] = float64(c.a4Stacks) * 0.075
			return m
		},
	})
}

func (c *char) hexSkillCDReduction() int {
	if !c.IsHexerei {
		return 0.5 * 60
	}
	return 1.0 * 60
}
