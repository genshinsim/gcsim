package viridescent

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("viridescent venerer", New)
	core.RegisterSetFunc("viridescentvenerer", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.AnemoP] = 0.15
		c.AddMod(coretype.CharStatMod{
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
				//overload dmg can't melt or vape so it's fine
				switch ai.AttackTag {
				case coretype.AttackTagSwirlCryo:
				case coretype.AttackTagSwirlElectro:
				case coretype.AttackTagSwirlHydro:
				case coretype.AttackTagSwirlPyro:
				default:
					return 0, false
				}
				return 0.6, false
			},
		})

		vvfunc := func(ele coretype.EleType, key string) func(args ...interface{}) bool {
			return func(args ...interface{}) bool {
				atk := args[1].(*coretype.AttackEvent)
				t := args[0].(coretype.Target)
				if atk.Info.ActorIndex != c.Index() {
					return false
				}

				//ignore if character not on field
				if s.Player.ActiveChar != c.Index()() {
					return false
				}

				t.AddResMod(key, core.ResistMod{
					Duration: 600, //10 seconds
					Ele:      ele,
					Value:    -0.4,
				})

				s.Log.NewEvent("vv 4pc proc", coretype.LogArtifactEvent, c.Index(), "reaction", key, "char", c.Index())

				return false
			}
		}
		s.Subscribe(coretype.OnSwirlCryo, vvfunc(coretype.Cryo, "vvcryo"), "vv4pc-"+c.Name())
		s.Subscribe(coretype.OnSwirlElectro, vvfunc(core.Electro, "vvelectro"), "vv4pc-"+c.Name())
		s.Subscribe(coretype.OnSwirlHydro, vvfunc(core.Hydro, "vvhydro"), "vv4pc-"+c.Name())
		s.Subscribe(coretype.OnSwirlPyro, vvfunc(core.Pyro, "vvpyro"), "vv4pc-"+c.Name())

		// Additional event for on damage proc on secondary targets
		// Got some very unexpected results when trying to modify the above vvfunc to allow for this, so I'm just copying it separately here
		// Possibly closure related? Not sure
		s.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
			atk := args[1].(*coretype.AttackEvent)
			t := args[0].(coretype.Target)
			if atk.Info.ActorIndex != c.Index() {
				return false
			}

			//ignore if character not on field
			if s.Player.ActiveChar != c.Index()() {
				return false
			}

			ele := atk.Info.Element
			key := "vv" + ele.String()
			switch atk.Info.AttackTag {
			case coretype.AttackTagSwirlCryo:
			case coretype.AttackTagSwirlElectro:
			case coretype.AttackTagSwirlHydro:
			case coretype.AttackTagSwirlPyro:
			default:
				return false
			}

			t.AddResMod(key, core.ResistMod{
				Duration: 600, //10 seconds
				Ele:      ele,
				Value:    -0.4,
			})

			s.Log.NewEvent("vv 4pc proc", coretype.LogArtifactEvent, c.Index(), "reaction", key, "char", c.Index())

			return false
		}, "vv4pc-secondary")

	}
	//add flat stat to char
}
