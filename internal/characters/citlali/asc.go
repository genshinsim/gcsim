package citlali

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	nightSoulGenerationIcd = "a1-ns-icd-key"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnMelt, c.a1Hook, "citlali-a1-onmelt")
	c.Core.Events.Subscribe(event.OnFrozen, c.a1Hook, "citlali-a1-onfrozen")
}

func (c *char) a1Hook(args ...interface{}) bool {
	t, ok := args[0].(*enemy.Enemy)
	if !ok {
		return false
	}
	if !c.nightsoulState.HasBlessing() {
		return false
	}

	if !c.StatusIsActive(nightSoulGenerationIcd) {
		c.AddStatus(nightSoulGenerationIcd, 8*60, true)
		c.generateNightsoulPoints(16)
		if c.Base.Cons >= 1 {
			c.numStellarBlades += 3
		}
	}

	amt := -0.2
	if c.Base.Cons >= 2 {
		amt = -0.4
	}

	t.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("citlali-a1-hydro", 12*60),
		Ele:   attributes.Hydro,
		Value: amt,
	})
	t.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag("citlali-a1-pyro", 12*60),
		Ele:   attributes.Pyro,
		Value: amt,
	})

	return false
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		c.generateNightsoulPoints(4)
		return false
	}, "citlali-a4-ns-gain")
}

func (c *char) a4Dmg(abil string) float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	em := c.NonExtraStat(attributes.EM)
	if abil == iceStormAbil {
		return 12 * em
	}
	if abil == frostFallAbil {
		return 0.9 * em
	}
	return 0
}
