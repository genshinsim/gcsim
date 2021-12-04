package gcsim

import "github.com/genshinsim/gcsim/pkg/core"

func (s *Simulation) initResonance(count map[core.EleType]int) {
	for k, v := range count {
		if v >= 2 {
			switch k {
			case core.Pyro:
				s.C.Log.Debugw("adding pyro resonance", "frame", s.C.F, "event", core.LogSimEvent)
				for _, c := range s.C.Chars {
					val := make([]float64, core.EndStatType)
					val[core.ATKP] = 0.25
					c.AddMod(core.CharStatMod{
						Key: "pyro-res",
						Amount: func(a core.AttackTag) ([]float64, bool) {
							return val, true
						},
						Expiry: -1,
					})
				}
			case core.Hydro:
				//heal not implemented yet
				s.C.Log.Debugw("adding hydro resonance", "frame", s.C.F, "event", core.LogSimEvent)
				s.C.Log.Warnw("hydro resonance not implemented", "event", core.LogSimEvent)
			case core.Cryo:
				s.C.Log.Debugw("adding cryo resonance", "frame", s.C.F, "event", core.LogSimEvent)
				s.C.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
					t := args[0].(core.Target)
					ds := args[1].(*core.Snapshot)
					if t.AuraType() == core.Cryo {
						ds.Stats[core.CR] += .15
						s.C.Log.Debugw("cryo resonance + 15% crit pre damage (cryo)", "frame", s.C.F, "event", core.LogCalc, "char", ds.ActorIndex, "next", ds.Stats[core.CR])
					}
					if t.AuraType() == core.Frozen {
						ds.Stats[core.CR] += .15
						s.C.Log.Debugw("cryo resonance + 15% crit pre damage  (frozen)", "frame", s.C.F, "event", core.LogCalc, "char", ds.ActorIndex, "next", ds.Stats[core.CR])
					}
					return false
				}, "cryo res")
			case core.Electro:
				s.C.Log.Debugw("adding electro resonance", "frame", s.C.F, "event", core.LogSimEvent)
				last := 0
				s.C.Events.Subscribe(core.OnReactionOccured, func(args ...interface{}) bool {
					ds := args[1].(*core.Snapshot)
					switch ds.ReactionType {
					case core.Overload, core.Superconduct, core.ElectroCharged:
					default:
						return false
					}
					if s.C.F-last < 300 && last != 0 { // every 5 seconds
						return false
					}
					s.C.Energy.DistributeParticle(core.Particle{
						Source: "electro res",
						Num:    1,
						Ele:    core.Electro,
					})
					last = s.C.F
					return false
				}, "electro res")
			case core.Geo:
				s.C.Log.Debugw("adding geo resonance", "frame", s.C.F, "event", core.LogSimEvent)
				//Increases shield strength by 15%. Additionally, characters protected by a shield will have the
				//following special characteristics:
				//	DMG dealt increased by 15%, dealing DMG to enemies will decrease their Geo RES by 20% for 15s.
				s.C.Shields.AddBonus(func() float64 {
					return 0.15 //shield bonus always active
				})
				s.C.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
					t := args[0].(core.Target)
					ds := args[1].(*core.Snapshot)
					if s.C.Shields.IsShielded() && s.C.ActiveChar == ds.ActorIndex {
						t.AddResMod("geo res", core.ResistMod{
							Duration: 15 * 60,
							Ele:      core.Geo,
							Value:    -0.2,
						})
					}
					return false
				}, "geo res")

				s.C.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
					ds := args[1].(*core.Snapshot)
					if s.C.Shields.IsShielded() && s.C.ActiveChar == ds.ActorIndex {
						ds.Stats[core.DmgP] += .15
						s.C.Log.Debugw("geo resonance + 15% DMG pre damage  (shielded)", "frame", s.C.F, "event", core.LogCalc, "char", ds.ActorIndex, "next", ds.Stats[core.DmgP])
					}
					return false
				}, "geo res")

			case core.Anemo:
				s.C.Log.Debugw("adding anemo resonance", "frame", s.C.F, "event", core.LogSimEvent)
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
