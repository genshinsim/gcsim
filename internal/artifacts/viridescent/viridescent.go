package viridescent

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("viridescent venerer", New)
	core.RegisterSetFunc("viridescentvenerer", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.AnemoP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "vv-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		//add +0.4 reaction damage
		c.AddReactBonusMod(core.ReactionBonusMod{
			Key:    "vv",
			Expiry: -1,
			Amount: func(ai core.AttackInfo) (float64, bool) {
				//check to make sure this is not an amped swirl
				if ai.Amped {
					return 0, false
				}
				switch ai.AttackTag {
				case core.AttackTagSwirlCryo:
				case core.AttackTagSwirlElectro:
				case core.AttackTagSwirlHydro:
				case core.AttackTagSwirlPyro:
				default:
					return 0, false
				}
				return 0.6, false
			},
		})

		vvfunc := func(ele core.EleType, key string) func(args ...interface{}) bool {
			return func(args ...interface{}) bool {
				atk := args[1].(*core.AttackEvent)
				t := args[0].(core.Target)
				if atk.Info.ActorIndex != c.CharIndex() {
					return false
				}

				//ignore if character not on field
				if s.ActiveChar != c.CharIndex() {
					return false
				}

				t.AddResMod(key, core.ResistMod{
					Duration: 600, //10 seconds
					Ele:      ele,
					Value:    -0.4,
				})

				s.Log.NewEvent("vv 4pc proc", core.LogArtifactEvent, c.CharIndex(), "reaction", key, "char", c.CharIndex())

				return false
			}
		}
		s.Events.Subscribe(core.OnSwirlCryo, vvfunc(core.Cryo, "vvcryo"), "vv4pc-"+c.Name())
		s.Events.Subscribe(core.OnSwirlElectro, vvfunc(core.Electro, "vvelectro"), "vv4pc-"+c.Name())
		s.Events.Subscribe(core.OnSwirlHydro, vvfunc(core.Hydro, "vvhydro"), "vv4pc-"+c.Name())
		s.Events.Subscribe(core.OnSwirlPyro, vvfunc(core.Pyro, "vvpyro"), "vv4pc-"+c.Name())

		// Additional event for on damage proc on secondary targets
		// Got some very unexpected results when trying to modify the above vvfunc to allow for this, so I'm just copying it separately here
		// Possibly closure related? Not sure
		s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			atk := args[1].(*core.AttackEvent)
			t := args[0].(core.Target)
			if atk.Info.ActorIndex != c.CharIndex() {
				return false
			}

			//ignore if character not on field
			if s.ActiveChar != c.CharIndex() {
				return false
			}

			ele := atk.Info.Element
			key := "vv" + ele.String()
			switch atk.Info.AttackTag {
			case core.AttackTagSwirlCryo:
			case core.AttackTagSwirlElectro:
			case core.AttackTagSwirlHydro:
			case core.AttackTagSwirlPyro:
			default:
				return false
			}

			t.AddResMod(key, core.ResistMod{
				Duration: 600, //10 seconds
				Ele:      ele,
				Value:    -0.4,
			})

			s.Log.NewEvent("vv 4pc proc", core.LogArtifactEvent, c.CharIndex(), "reaction", key, "char", c.CharIndex())

			return false
		}, fmt.Sprintf("vv4pc-secondary-%v", c.Name()))

	}
	//add flat stat to char
}
