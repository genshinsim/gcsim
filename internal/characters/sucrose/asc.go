package sucrose

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When Sucrose triggers a Swirl reaction, all characters in the party with the matching element (excluding Sucrose)
// have their Elemental Mastery increased by 50 for 8s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	c.a1Buff = make([]float64, attributes.EndStatType)
	c.a1Buff[attributes.EM] = 50
	swirlfunc := func(ele attributes.Element) func(args ...interface{}) bool {
		icd := -1
		return func(args ...interface{}) bool {
			if _, ok := args[0].(*gadget.Gadget); ok {
				return false
			}

			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			// do not overwrite mod if same frame
			//TODO: this probably isn't needed?
			if c.Core.F < icd {
				return false
			}
			icd = c.Core.F + 1

			for _, char := range c.Core.Player.Chars() {
				this := char
				if this.Base.Element != ele {
					continue
				}
				this.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("sucrose-a1", 480), //8s
					AffectedStat: attributes.EM,
					Amount: func() ([]float64, bool) {
						return c.a1Buff, true
					},
				})
			}

			c.Core.Log.NewEvent("sucrose a1 triggered", glog.LogCharacterEvent, c.Index).
				Write("reaction", "swirl-"+ele.String()).
				Write("expiry", c.Core.F+480)
			return false
		}
	}

	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc(attributes.Cryo), "sucrose-a1-cryo")
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc(attributes.Electro), "sucrose-a1-electro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc(attributes.Hydro), "sucrose-a1-hydro")
	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc(attributes.Pyro), "sucrose-a1-pyro")
}

// When Astable Anemohypostasis Creation - 6308 or Forbidden Creation - Isomer 75 / Type II hits an opponent,
// increases all party members' (excluding Sucrose) Elemental Mastery by an amount equal to 20% of Sucrose's Elemental Mastery for 8s.
//
// - called inside of an attack callback in skill.go and burst.go
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4Buff[attributes.EM] = c.NonExtraStat(attributes.EM) * .20
	for i, char := range c.Core.Player.Chars() {
		if i == c.Index {
			continue //nothing for sucrose
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("sucrose-a4", 480), //8 s
			AffectedStat: attributes.EM,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return c.a4Buff, true
			},
		})
	}

	c.Core.Log.NewEvent("sucrose a4 triggered", glog.LogCharacterEvent, c.Index).
		Write("em snapshot", c.a4Buff[attributes.EM]).
		Write("expiry", c.Core.F+480)
}
