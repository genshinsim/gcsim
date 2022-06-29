package sucrose

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 50

	swirlfunc := func(ele attributes.Element) func(args ...interface{}) bool {
		icd := -1
		return func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			// do not overwrite mod if same frame
			if c.Core.F < icd {
				return false
			}
			icd = c.Core.F + 1

			dur := 60 * 8
			for _, char := range c.Core.Player.Chars() {
				this := char
				if this.Base.Element != ele {
					continue
				}
				this.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("sucrose-a1", dur),
					AffectedStat: attributes.EM,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}

			c.Core.Log.NewEvent("sucrose a1 triggered", glog.LogCharacterEvent, c.Index, "reaction", "swirl-"+ele.String(), "expiry", c.Core.F+dur)
			return false
		}
	}

	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc(attributes.Cryo), "sucrose-a1-cryo")
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc(attributes.Electro), "sucrose-a1-electro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc(attributes.Hydro), "sucrose-a1-hydro")
	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc(attributes.Pyro), "sucrose-a1-pyro")
}

func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = c.Stat(attributes.EM) * .20

	dur := 60 * 8
	c.Core.Status.Add("sucrosea4", dur)
	for i, char := range c.Core.Player.Chars() {
		if i == c.Index {
			continue //nothing for sucrose
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("sucrose-a4", dur),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	c.Core.Log.NewEvent("sucrose a4 triggered", glog.LogCharacterEvent, c.Index, "em snapshot", m[attributes.EM], "expiry", c.Core.F+dur)
}
