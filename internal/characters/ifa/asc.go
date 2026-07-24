package ifa

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("ifa-a1", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !c.nightsoulState.HasBlessing() {
					return 0
				}

				m := c.getTeamNightsoul()

				m += c.c2BonusPoints(m)

				m = min(m, 150.0+c.c2CapIncrease())

				c.Core.Log.NewEvent("ifa a1 stacks", glog.LogCharacterEvent, char.Index()).
					Write("stacks", m)

				switch ai.AttackTag {
				case attacks.AttackTagSwirlPyro,
					attacks.AttackTagSwirlHydro,
					attacks.AttackTagSwirlElectro,
					attacks.AttackTagSwirlCryo,
					attacks.AttackTagECDamage:
					return 0.015 * m
				case attacks.AttackTagReactionLunarCharge,
					attacks.AttackTagDirectLunarCharged:
					return 0.002 * m
				}
				return 0
			},
		})
	}
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 80

	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...any) {
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("ifa-a4", 10*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})
	}, "ifa-a4-on-nightsoul-burst")
}
