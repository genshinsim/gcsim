package sayu

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const a1ICDKey = "sayu-a1-icd"

// When Sayu triggers a Swirl reaction while active, she heals all your
// characters and nearby allies for 300 HP. She will also heal an additional 1.2
// HP for every point of Elemental Mastery she has.  This effect can be
// triggered once every 2s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	swirlfunc := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}

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
		c.AddStatus(a1ICDKey, 120, true) // 2s

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

	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc, "sayu-a1-cryo")
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc, "sayu-a1-electro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc, "sayu-a1-hydro")
	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc, "sayu-a1-pyro")
}

// The Muji-Muji Daruma created by Yoohoo Art: Mujina Flurry gains the following effects:
//
// - When healing a character, it will also heal characters near that healed character for 20% the amount of HP.
//   - only relevant in Co-Op
//
// - Increases the AoE of its attack against opponents
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.qTickRadius = 3.5
}
