package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

//Upon triggering a Swirl reaction, Kaedehara Kazuha will grant all party members a 0.04%
//Elemental DMG Bonus to the element absorbed by Swirl for every point of Elemental Mastery
//he has for 8s. Bonuses for different elements obtained through this method can co-exist.
//
//this ignores any EM he gets from Sucrose A4, which is: When Astable Anemohypostasis Creation
//- 6308 or Forbidden Creation - Isomer 75 / Type II hits an opponent, increases all party
//members' (excluding Sucrose) Elemental Mastery by an amount equal to 20% of Sucrose's
//Elemental Mastery for 8s.
//
//he still benefits from sucrose em but just cannot share it
func (c *char) a4() {
	m := make([]float64, attributes.EndStatType)

	swirlfunc := func(ele attributes.Stat, key string) func(args ...interface{}) bool {
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

			//recalc em
			dmg := 0.0004 * c.Stat(attributes.EM)

			for _, char := range c.Core.Player.Chars() {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("kazuha-a4-"+key, 60*8),
					AffectedStat: attributes.NoStat,
					Amount: func() ([]float64, bool) {
						m[attributes.CryoP] = 0
						m[attributes.ElectroP] = 0
						m[attributes.HydroP] = 0
						m[attributes.PyroP] = 0
						m[ele] = dmg
						return m, true
					},
				})
			}

			c.Core.Log.NewEvent("kazuha a4 proc", glog.LogCharacterEvent, c.Index).
				Write("reaction", ele.String())

			return false
		}
	}

	c.Core.Events.Subscribe(event.OnSwirlCryo, swirlfunc(attributes.CryoP, "cryo"), "kazuha-a4-cryo")
	c.Core.Events.Subscribe(event.OnSwirlElectro, swirlfunc(attributes.ElectroP, "electro"), "kazuha-a4-electro")
	c.Core.Events.Subscribe(event.OnSwirlHydro, swirlfunc(attributes.HydroP, "hydro"), "kazuha-a4-hydro")
	c.Core.Events.Subscribe(event.OnSwirlPyro, swirlfunc(attributes.PyroP, "pyro"), "kazuha-a4-pyro")
}

func (c *char) absorbCheckA1(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}
		c.a1Ele = c.Core.Combat.AbsorbCheck(c.infuseCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

		if c.a1Ele != attributes.NoElement {
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"kazuha a1 infused ", c.a1Ele.String(),
			)
			return
		}
		//otherwise queue up
		c.Core.Tasks.Add(c.absorbCheckA1(src, count+1, max), 6)
	}
}
