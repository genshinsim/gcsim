package pyro

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const a4OnReactICD = "travelerpyro-a4-icd"

func (c *Traveler) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	fReactionHook := func(args ...any) bool {
		// Attack must be from active character
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if !c.nightsoulState.HasBlessing() {
			return false
		}
		if c.StatusIsActive(a4OnReactICD) {
			return false
		}

		c.AddStatus(a4OnReactICD, 12*60, true)
		c.AddEnergy("travelerpyro-a4-energy", 5)
		return false
	}

	fNSHook := func(args ...any) bool {
		c.AddEnergy("travelerpyro-a4-energy", 4)
		return false
	}

	c.Core.Events.Subscribe(event.OnBurning, fReactionHook, "travelerpyro-a4-onburning")
	c.Core.Events.Subscribe(event.OnVaporize, fReactionHook, "travelerpyro-a4-onvaporize")
	c.Core.Events.Subscribe(event.OnMelt, fReactionHook, "travelerpyro-a4-onmelt")
	c.Core.Events.Subscribe(event.OnOverload, fReactionHook, "travelerpyro-a4-onoverload")
	c.Core.Events.Subscribe(event.OnBurgeon, fReactionHook, "travelerpyro-a4-onburgeon")
	c.Core.Events.Subscribe(event.OnSwirlPyro, fReactionHook, "travelerpyro-a4-onswirlpyro")
	c.Core.Events.Subscribe(event.OnCrystallizePyro, fReactionHook, "travelerpyro-a4-oncrystallizepyro")
	c.Core.Events.Subscribe(event.OnNightsoulBurst, fNSHook, "travelerpyro-a4-onnsburst")
}
