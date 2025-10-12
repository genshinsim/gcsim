package flins

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
	a1Key               = "flins-a1"
	a4Key               = "flins-a4"
	lunarchargeBonusKey = "flins-lc-bonus"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(a1Key, -1),
		Amount: func(ai info.AttackInfo) (float64, bool) {
			if c.getAscendeantGleamLevel() < 2 {
				return 0, false
			}

			switch ai.AttackTag {
			case attacks.AttackTagDirectLunarCharged:
			case attacks.AttackTagReactionLunarCharge:
			default:
				return 0, false
			}
			return 0.2, false
		},
	})
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	scale, max := c.c4A4()

	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a4Key, -1),
		Extra:        true,
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
			m[attributes.EM] = min(stats.TotalATK()*scale, max)
			return m, true
		},
	})
}

func (c *char) lunarchargeInit() {
	c.Core.Flags.Custom[reactable.LunarChargeEnableKey] = 1

	// TODO: moonsign?

	// TODO: every 100 ATK that Ineffa has increasing Lunar-Charged's Base DMG by 0.7%, up to a maximum of 14%.
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCharged:
		case attacks.AttackTagReactionLunarCharge:
		default:
			return false
		}

		stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
		bonus := min(stats.TotalATK()/100.0*0.007, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("flins adding lunarcharged base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
		return false
	}, lunarchargeBonusKey)
}
