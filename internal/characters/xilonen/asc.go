package xilonen

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1IcdKey = "xilonen-a1-icd"
	a1Key    = "xilonen-a1"

	a4IcdKey = "xilonen-a4-icd"
	a4Key    = "xilonen-a4"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	if c.samplersConverted >= 2 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.30
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a1Key, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagPlunge && atk.Info.AttackTag != attacks.AttackTagNormal {
				return nil, false
			}
			if !slices.Contains(atk.Info.AdditionalTags, attacks.AdditionalTagNightsoul) {
				return nil, false
			}
			return m, true
		},
	})

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		next := args[1].(int)
		if c.StatusIsActive(activeSamplerKey) {
			ch := c.Core.Player.Chars()[next]
			c.activeSamplers(c.sampleSrc, ch)()
		}

		return false
	}, "xilonen-a1-swap")
}

func (c *char) a1cb(cb combat.AttackCB) {
	if c.Base.Ascension < 1 {
		return
	}
	if !c.nightsoulState.HasBlessing() {
		return
	}
	if c.samplersConverted < 2 || c.samplersActivated {
		return
	}
	if c.StatusIsActive(a1IcdKey) {
		return
	}

	c.AddStatus(a1IcdKey, 0.1*60, true)
	c.nightsoulState.GeneratePoints(35)
	if c.nightsoulState.Points() < c.nightsoulState.MaxPoints {
		return
	}
	c.a4MaxPoints(cb.Target, cb.AttackEvent)
	c.a1MaxPoints()
}

func (c *char) a1MaxPoints() {
	c.nightsoulState.ConsumePoints(c.nightsoulState.Points())
	if !c.StatusIsActive(c6key) {
		c.exitNightsoul()
	}

	c.AddStatus(activeSamplerKey, 15*60, true)
	c.sampleSrc = c.Core.F
	c.activeSamplers(c.sampleSrc, c.CharWrapper)()
	c.samplersActivated = true

	c.c2activate()
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DEFP] = 0.20

	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a4Key, 15*60),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, a4Key)
}

func (c *char) a4MaxPoints(t combat.Target, ae *combat.AttackEvent) {
	if c.Base.Ascension < 4 {
		return
	}
	if c.StatusIsActive(a4IcdKey) {
		return
	}
	c.AddStatus(a4IcdKey, 14*60, true)
	c.Core.Events.Emit(event.OnNightsoulBurst, t, ae)
}
