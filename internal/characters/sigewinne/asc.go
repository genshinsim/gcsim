package sigewinne

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	baseConvalescenceHP       = 30000
	convalescenceHPFraction   = 1000
	flatconvalescenceIncrease = 80
	flatconvalescenceCap      = 2800
	a1DmgBuff                 = 0.08
	a4HpDebtHealingBonusRatio = 0.03 / 1000
	convalescenceKey          = "sigewinne-convalescence"
)

func (c *char) a1() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagElementalArt || atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}

		if atk.Info.ActorIndex == c.Index {
			return false
		}

		active := c.Core.Player.ActiveChar()
		if active.Index == atk.Info.ActorIndex {
			return false
		}

		if c.Tags[convalescenceKey] > 0 {
			c.Tags[convalescenceKey]--
			amt := 0.0
			if c.Base.Cons >= 1 {
				amt = max(C1flatconvalescenceCap, max(c.MaxHP()-baseConvalescenceHP, 0)/convalescenceHPFraction*C1flatconvalescenceIncrease)
			} else {
				amt = max(flatconvalescenceCap, max(c.MaxHP()-baseConvalescenceHP, 0)/convalescenceHPFraction*flatconvalescenceIncrease)
			}

			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Shenhe Quill proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt).
					Write("effect_ends_at", c.StatusExpiry(convalescenceKey)).
					Write("quill_left", c.Tags[convalescenceKey])
			}

			atk.Info.FlatDmg += amt
		}

		return false
	}, "sigewinne-convalescence-hook")
}

func (c *char) a4() {
	c.AddHealBonusMod(character.HealBonusMod{
		Base: modifier.NewBase("sigewinne-a4-healing-bonus", -1),
		Amount: func() (float64, bool) {
			totalHpDebt := 0.
			for _, other := range c.Core.Player.Chars() {
				totalHpDebt += other.CurrentHPDebt()
			}
			return totalHpDebt * a4HpDebtHealingBonusRatio, true
		},
	})
}
