package sigewinne

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	baseConvalescenceHP       = 30000
	flatConvalescenceIncrease = 80.0 / 1000
	flatConvalescenceCap      = 2800
	a1DmgBuff                 = 0.08
	a4HpDebtHealingBonusRatio = 0.03 / 1000
	a4HealingBonusCap         = 0.3
	convalescenceKey          = "sigewinne-convalescence"
)

func (c *char) a1() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		default:
			return false
		}

		if atk.Info.ActorIndex == c.Index {
			return false
		}

		active := c.Core.Player.ActiveChar()
		if active.Index == atk.Info.ActorIndex {
			return false
		}

		if !c.StatusIsActive(convalescenceKey) || c.Tag(convalescenceKey) == 0 {
			return false
		}
		c.SetTag(convalescenceKey, c.Tag(convalescenceKey)-1)

		hp := max(c.MaxHP()-baseConvalescenceHP, 0)
		var amt float64
		if c.Base.Cons >= 1 {
			amt = min(c1FlatConvalescenceCap, hp*c1FlatConvalescenceIncrease)
		} else {
			amt = min(flatConvalescenceCap, hp*flatConvalescenceIncrease)
		}

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Sigewinne A1 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("effect_ends_at", c.StatusExpiry(convalescenceKey)).
				Write("quill_left", c.Tag(convalescenceKey))
		}

		atk.Info.FlatDmg += amt

		return false
	}, "sigewinne-convalescence-hook")
}

func (c *char) a1Self() {
	c.AddStatus(convalescenceKey, skillCD*60, true)
	c.SetTag(convalescenceKey, 10)

	buff := make([]float64, attributes.EndStatType)
	buff[attributes.HydroP] = a1DmgBuff
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("sigewinne-a1", skillCD*60),
		Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			return buff, true
		},
	})
}

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("sigewinne-a4", -1),
		AffectedStat: attributes.Heal,
		Amount: func() ([]float64, bool) {
			totalHpDebt := 0.
			for _, other := range c.Core.Player.Chars() {
				totalHpDebt += other.CurrentHPDebt()
			}
			heal := min(a4HealingBonusCap, totalHpDebt*a4HpDebtHealingBonusRatio)
			m[attributes.Heal] = heal
			return m, true
		},
	})
}
