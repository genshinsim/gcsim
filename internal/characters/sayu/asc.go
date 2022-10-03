package sayu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const a1ICDKey = "sayu-a1-icd"

// A1:
// When Sayu triggers a Swirl reaction while active, she heals all your
// characters and nearby allies for 300 HP. She will also heal an additional 1.2
// HP for every point of Elemental Mastery she has.  This effect can be
// triggered once every 2s.
func (c *char) a1() {
	swirlfunc := func(ele attributes.Element) func(evt event.EventPayload) bool {
		return func(evt event.EventPayload) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			if c.Core.Player.Active() != c.Index {
				return false
			}
			if c.StatusIsActive(a1ICDKey) {
				return false
			}
			c.AddStatus(a1ICDKey, 120, true) //2s

			if c.Base.Cons >= 4 {
				c.AddEnergy("sayu-c4", 1.2)
			}

			heal := 300 + c.Stat(attributes.EM)*1.2
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Someone More Capable",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			})

			return false
		}
	}

	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc(attributes.Cryo), "sayu-a1-cryo")
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc(attributes.Electro), "sayu-a1-electro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc(attributes.Hydro), "sayu-a1-hydro")
	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc(attributes.Pyro), "sayu-a1-pyro")
}
