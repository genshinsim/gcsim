package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (s *Simulation) initResonance(count map[coretype.EleType]int) {
	for k, v := range count {
		if v >= 2 {
			switch k {
			case core.Pyro:
				s.C.Log.NewEvent("adding pyro resonance", coretype.LogSimEvent, -1)
				for _, c := range s.C.Chars {
					val := make([]float64, core.EndStatType)
					val[core.ATKP] = 0.25
					c.AddMod(coretype.CharStatMod{
						Key: "pyro-res",
						Amount: func() ([]float64, bool) {
							return val, true
						},
						Expiry: -1,
					})
				}
			case core.Hydro:
				//heal not implemented yet
				s.C.Log.NewEvent("adding hydro resonance", coretype.LogSimEvent, -1)
				s.C.Health.AddIncHealBonus(func(healedCharIndex int) float64 {
					return 0.3
				})
			case coretype.Cryo:
				s.C.Log.NewEvent("adding cryo resonance", coretype.LogSimEvent, -1)
				val := make([]float64, core.EndStatType)
				val[core.CR] = .15
				for _, c := range s.C.Chars {
					c.AddPreDamageMod(coretype.PreDamageMod{
						Key: "cryo-res",
						Amount: func(ae *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
							if t.AuraContains(coretype.Cryo) || t.AuraContains(coretype.Frozen) {
								return val, true
							}
							return nil, false
						},
						Expiry: -1,
					})
				}
			case core.Electro:
				s.C.Log.NewEvent("adding electro resonance", coretype.LogSimEvent, -1)
				last := 0
				recover := func(args ...interface{}) bool {
					if s.C.Frame-last < 300 && last != 0 { // every 5 seconds
						return false
					}
					s.C.Energy.DistributeParticle(core.Particle{
						Source: "electro res",
						Num:    1,
						Ele:    core.Electro,
					})
					last = s.C.Frame
					return false
				}
				s.C.Subscribe(core.OnOverload, recover, "electro-res")
				s.C.Subscribe(core.OnSuperconduct, recover, "electro-res")
				s.C.Subscribe(core.OnElectroCharged, recover, "electro-res")

			case core.Geo:
				s.C.Log.NewEvent("adding geo resonance", coretype.LogSimEvent, -1)
				//Increases shield strength by 15%. Additionally, characters protected by a shield will have the
				//following special characteristics:
				//	DMG dealt increased by 15%, dealing DMG to enemies will decrease their Geo RES by 20% for 15s.
				s.C.Player.AddShieldBonus(func() float64 {
					return 0.15 //shield bonus always active
				})
				//shred geo res of target
				s.C.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
					t := args[0].(coretype.Target)
					atk := args[1].(*coretype.AttackEvent)
					if s.C.Player.IsCharShielded(atk.Info.ActorIndex) {
						t.AddResMod("geo res", core.ResistMod{
							Duration: 15 * 60,
							Ele:      core.Geo,
							Value:    -0.2,
						})
					}
					return false
				}, "geo res")

				val := make([]float64, core.EndStatType)
				val[core.DmgP] = .15
				for _, c := range s.C.Chars {
					c.AddPreDamageMod(coretype.PreDamageMod{
						Key: "geo-res",
						Amount: func(ae *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
							if s.C.Player.IsCharShielded(ae.Info.ActorIndex) {
								return val, true
							}
							return nil, false
						},
						Expiry: -1,
					})
				}

			case core.Anemo:
				s.C.Log.NewEvent("adding anemo resonance", coretype.LogSimEvent, -1)
				for _, c := range s.C.Chars {
					c.AddCDAdjustFunc(core.CDAdjust{
						Key:    "anemo-res",
						Amount: func(a core.ActionType) float64 { return -0.05 },
						Expiry: -1,
					})
				}
			}
		}
	}
}
