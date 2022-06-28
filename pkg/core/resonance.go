package core

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (s *Core) SetupResonance() {
	chars := s.Player.Chars()
	if len(chars) < 4 {
		return //no resonance if less than 4 chars
	}
	//count number of ele first
	count := make(map[attributes.Element]int)
	for _, c := range chars {
		count[c.Base.Element]++
	}

	for k, v := range count {
		if v >= 2 {
			switch k {
			case attributes.Pyro:
				val := make([]float64, attributes.EndStatType)
				val[attributes.ATKP] = 0.25
				f := func() ([]float64, bool) {
					return val, true
				}
				for _, c := range chars {
					c.AddStatMod(character.StatMod{Base: modifier.NewBase("pyro-res", -1), AffectedStat: attributes.NoStat, Amount: f})
				}
			case attributes.Hydro:
				//heal not implemented yet
				f := func() (float64, bool) {
					return 0.3, true
				}
				for _, c := range chars {
					c.AddHealBonusMod(character.HealBonusMod{Base: modifier.NewBase("hydro-res", -1), Amount: f})
				}
			case attributes.Cryo:
				val := make([]float64, attributes.EndStatType)
				val[attributes.CR] = .15
				f := func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					r, ok := t.(Reactable)
					if !ok {
						return nil, false
					}
					if r.AuraContains(attributes.Cryo) || r.AuraContains(attributes.Frozen) {
						return val, true
					}
					return nil, false
				}
				for _, c := range chars {
					c.AddAttackMod(character.AttackMod{Base: modifier.NewBase("cyro-res", -1), Amount: f})
				}
			case attributes.Electro:
				last := 0
				recover := func(args ...interface{}) bool {
					if s.F-last < 300 && last != 0 { // every 5 seconds
						return false
					}
					s.Player.DistributeParticle(character.Particle{
						Source: "electro-res",
						Num:    1,
						Ele:    attributes.Electro,
					})
					last = s.F
					return false
				}
				s.Events.Subscribe(event.OnOverload, recover, "electro-res")
				s.Events.Subscribe(event.OnSuperconduct, recover, "electro-res")
				s.Events.Subscribe(event.OnElectroCharged, recover, "electro-res")

			case attributes.Geo:
				//Increases shield strength by 15%. Additionally, characters protected by a shield will have the
				//following special characteristics:

				//	DMG dealt increased by 15%, dealing DMG to enemies will decrease their Geo RES by 20% for 15s.
				f := func() (float64, bool) { return 0.15, true }
				s.Player.Shields.AddShieldBonusMod("geo-res", -1, f)

				//shred geo res of target
				s.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
					t, ok := args[0].(Enemy)
					if !ok {
						return false
					}
					atk := args[1].(*combat.AttackEvent)
					if s.Player.Shields.PlayerIsShielded() && s.Player.Active() == atk.Info.ActorIndex {
						t.AddResistMod("geo-res", 15*60, attributes.Geo, -0.2)
					}
					return false
				}, "geo res")

				val := make([]float64, attributes.EndStatType)
				val[attributes.DmgP] = .15
				atkf := func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					if s.Player.Shields.PlayerIsShielded() && s.Player.Active() == ae.Info.ActorIndex {
						return val, true
					}
					return nil, false
				}
				for _, c := range chars {
					c.AddAttackMod(character.AttackMod{Base: modifier.NewBase("geo-res", -1), Amount: atkf})
				}

			case attributes.Anemo:
				for _, c := range chars {
					c.AddCooldownMod(character.CooldownMod{Base: modifier.NewBase("anemo-res", -1), Amount: func(a action.Action) float64 { return -0.05 }})
				}
			}
		}
	}
}
