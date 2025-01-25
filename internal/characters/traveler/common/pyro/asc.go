package pyro

import "github.com/genshinsim/gcsim/pkg/core/event"

const a4OnReactICD = "travelerpyro-a4-icd"

func (c *Traveler) a4() {
	fReactionHook := func(args ...interface{}) bool {
		// if PMC is in NS Blessinng, then the active character should
		// be inside a Blazing Threshold or Scorching Threshold since it
		// always follows active character
		if !c.nightsoulState.HasBlessing() {
			return false
		}
		if c.StatusIsActive(a4OnReactICD) {
			return false
		}

		if c.Base.Cons >= 2 && c.StatusIsActive(c2StatusKey) && c.c2ActivationsPerSkill < 2 {
			c.nightsoulState.GeneratePoints(14)
		}

		c.AddStatus(a4OnReactICD, 12, false) // TODO: hitlag affected?
		c.AddEnergy("travelerpyro-a4-energy", 5)
		return false
	}

	fNSHook := func(args ...interface{}) bool {
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
