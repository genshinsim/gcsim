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
	c.a1buff = make([]float64, attributes.EndStatType)
	c.a1buff[attributes.EM] = 50

	swirlfunc := func(ele attributes.Element) func(evt event.EventPayload) bool {
		icd := -1
		return func(evt event.EventPayload) bool {
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
						return c.a1buff, true
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

func (c *char) a4() {
	c.a4buff = make([]float64, attributes.EndStatType)
	c.a4buff[attributes.EM] = c.Stat(attributes.EM) * .20

	for i, char := range c.Core.Player.Chars() {
		if i == c.Index {
			continue //nothing for sucrose
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("sucrose-a4", 480), //8 s
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return c.a4buff, true
			},
		})
	}

	c.Core.Log.NewEvent("sucrose a4 triggered", glog.LogCharacterEvent, c.Index).
		Write("em snapshot", c.a4buff[attributes.EM]).
		Write("expiry", c.Core.F+480)
}
