package venti

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	hexAttackCBICDKey = "venti-hexerei-attack-cb-icd"
)

// A1 is not implemented and will likely never be implemented:
// Holding Skyward Sonnet creates an upcurrent that lasts for 20s.

// Regenerates 15 Energy for Venti after the effects of Wind's Grand Ode end.
// If an Elemental Absorption occurred, this also restores 15 Energy to all characters of that corresponding element in the party.
//
// - checks for ascension level in burst.go to avoid queuing this up only to fail the ascension level check
func (c *char) a4() {
	c.AddEnergy("venti-a4", 15)
	if c.qAbsorb == attributes.NoElement {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		if char.Base.Element == c.qAbsorb {
			char.AddEnergy("venti-a4", 15)
		}
	}
}

func (c *char) hexInit() {
	if !c.IsHexerei {
		return
	}
	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.5

	onSwirlBuff := func(args ...any) {
		if !c.StatusIsActive(burstKey) {
			return
		}

		atk := args[1].(*info.AttackEvent)
		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)

		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("venti-hexerei-dmg", 4*60),
			Amount: func(_ *info.AttackEvent, _ info.Target) []float64 {
				return m
			},
		})

		c.burstHexMult = 1.35
		c.burstHexSrc = c.Core.F

		resetHexMult := func(src int) func() {
			return func() {
				if c.burstHexSrc != src {
					return
				}
				c.burstHexMult = 1
			}
		}

		char.QueueCharTask(resetHexMult(c.Core.F), 4*60)
	}

	c.Core.Events.Subscribe(event.OnSwirlPyro, onSwirlBuff, "venti-hexerei-on-swirl-pyro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, onSwirlBuff, "venti-hexerei-on-swirl-hydro")
	c.Core.Events.Subscribe(event.OnSwirlElectro, onSwirlBuff, "venti-hexerei-on-swirl-electro")
	c.Core.Events.Subscribe(event.OnSwirlCryo, onSwirlBuff, "venti-hexerei-on-swirl-cryo")
}

func (c *char) hexAttackCB(_ info.AttackCB) {
	if !c.IsHexerei {
		return
	}
	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}
	if !c.StatusIsActive(burstKey) {
		return
	}
	if c.StatusIsActive(hexAttackCBICDKey) {
		return
	}

	c.AddStatus(hexAttackCBICDKey, 0.1*60, true)

	c.ReduceActionCooldown(action.ActionBurst, -0.5*60)

	c.ExtendStatus(burstKey, 1*60)
}
