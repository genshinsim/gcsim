package qiqi

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key    = "qiqi-c1"
	c1ICDKey = "qiqi-c1-icd"
	c2Key    = "qiqi-c2"
	c6Key    = "qiqi-c6"
)

// Qiqi regenerates 2 Energy when the Herald of Frost hits an opponent marked by a
// Fortune-Preserving Talisman.
func (c *char) c1CB(a info.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}

	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}

	if !e.StatusIsActive(talismanKey) {
		return
	}

	c.AddEnergy(c1Key, 2)
	c.Core.Log.NewEvent("Qiqi C1 Activation - Adding 2 energy", glog.LogCharacterEvent, c.Index()).
		Write("target", a.Target.Key())
}

// Radiance: Stellar Glimmer: Qiqi also regenerates 6 Elemental Energy when the Herald of Frost's
// attacks and coordinated attacks hit opponents. This effect can trigger once every 6s.
func (c *char) c1RevelationCB(a info.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}

	if !c.revelation {
		return
	}

	if c.getRadiance() == radianceNone {
		return
	}

	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(c1ICDKey) {
		return
	}

	c.AddStatus(c1ICDKey, 6*60, true)

	c.AddEnergy(c1Key, 6)
}

// Qiqi's Normal and Charge Attack DMG against opponents affected by Cryo is increased by 15%.
// Radiance: Stellar Glimmer: Qiqi's ATK is increased by 50%.
func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = .15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c2Key, -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil
			}

			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil
			}
			if !e.AuraContains(attributes.Cryo, attributes.Frozen) {
				return nil
			}

			return m
		},
	})

	if !c.revelation {
		return
	}

	mSSC := make([]float64, attributes.EndStatType)
	mSSC[attributes.ATKP] = .5
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c2Key+"-radiance", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			if c.getRadiance() == radianceNone {
				return nil
			}
			return mSSC
		},
	})
}

// Targets marked by a Fortune-Preserving Talisman have their ATK reduced by 20%.
// Additionally, when the healing effect of Fortune-Preserving Talisman triggers, the nearby party
// member with the lowest HP percentage remaining will also receive an additional instance of
// healing for 180% of Qiqi's ATK.
func (c *char) c4OnHeal() {
	if c.Base.Cons < 4 {
		return
	}

	if !c.revelation {
		return
	}

	lowestHPInd := c.Index()
	lowestHP := c.CurrentHPRatio()
	for i, char := range c.Core.Player.Chars() {
		if char.CurrentHPRatio() < lowestHP {
			lowestHPInd = i
			lowestHP = char.CurrentHPRatio()
		}
	}

	healAmt := 1.8 * c.TotalAtk()
	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  lowestHPInd,
		Message: "Fortune-Preserving Talisman (C4)",
		Src:     healAmt,
		Bonus:   c.Stat(attributes.Heal),
	})
}

// Using Adeptus Art: Preserver of Fortune revives all fallen party members nearby and regenerates
// 50% of their HP. This effect can only occur once every 15 mins.
//
// Additionally, after using Adeptus Art: Preserver of Fortune, Qiqi will also gain 4 stacks of
// Glimpse of Mystery, and when your current active character (other than Qiqi) deals
// Stellar-Conduct or Stellar Swirl reaction DMG with their Talent or skill, 1 stack of Glimpse of
// Mystery is consumed, and DMG dealt is increased by 600% of Qiqi's ATK. The Glimpse of Mystery
// effect lasts up to 12s, and unleashing Adeptus Art: Preserver of Fortune refreshes the stack
// count.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	if !c.revelation {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if c.c6Stacks <= 0 {
			return
		}

		// Qiqi can't activate
		if atk.Info.ActorIndex == c.Index() {
			return
		}

		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		if !c.StatusIsActive(c6Key) {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		default:
			return
		}

		c.c6Stacks -= 1
		amt := 6 * c.TotalAtk()
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Qiqi C6 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("c6 stacks left", c.c6Stacks)
		}
		atk.Info.FlatDmg += amt
	}, "qiqi-c6-dmg")
}

func (c *char) c6OnBurst() {
	if c.Base.Cons < 6 {
		return
	}

	if !c.revelation {
		return
	}

	c.AddStatus(c6Key, 12*60, true)
	c.c6Stacks = 4
}
