package kazuha

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// If Chihayaburu comes into contact with Hydro/Pyro/Cryo/Electro when cast, this Chihayaburu will absorb that element and if Plunging Attack: Midare Ranzan is used before the effect expires, it will deal an additional 200% ATK of the absorbed elemental type as DMG. This will be considered Plunging Attack DMG.
// Elemental Absorption may only occur once per use of Chihayaburu.
//
// - checks for ascension level in skill.go to avoid queuing this up only to fail the ascension level check
func (c *char) absorbCheckA1(src, count, maxcount int) func() {
	return func() {
		if count == maxcount {
			return
		}
		c.a1Absorb = c.Core.Combat.AbsorbCheck(c.Index, c.a1AbsorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

		if c.a1Absorb != attributes.NoElement {
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"kazuha a1 absorbed ", c.a1Absorb.String(),
			)
			return
		}
		// otherwise queue up
		c.Core.Tasks.Add(c.absorbCheckA1(src, count+1, maxcount), 6)
	}
}

// Upon triggering a Swirl reaction, Kaedehara Kazuha will grant all party members a 0.04%
// Elemental DMG Bonus to the element absorbed by Swirl for every point of Elemental Mastery
// he has for 8s. Bonuses for different elements obtained through this method can co-exist.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)

	swirlfunc := func(ele attributes.Stat, key string) func(args ...interface{}) bool {
		icd := -1
		return func(args ...interface{}) bool {
			if _, ok := args[0].(*enemy.Enemy); !ok {
				return false
			}

			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			// do not overwrite mod if same frame
			if c.Core.F < icd {
				return false
			}
			icd = c.Core.F + 1

			// recalc em
			dmg := 0.0004 * c.NonExtraStat(attributes.EM)

			for _, char := range c.Core.Player.Chars() {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("kazuha-a4-"+key, 60*8),
					AffectedStat: ele,
					Extra:        true,
					Amount: func() ([]float64, bool) {
						clear(m)
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
