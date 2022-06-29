package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func SetupTargetsInCore(core *core.Core, p core.Coord, targets []enemy.EnemyProfile) error {

	// s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
	// s.stats.ElementUptime[0] = make(map[core.EleType]int)

	if p.R == 0 {
		return errors.New("player cannot have 0 radius")
	}
	player := avatar.New(core, p.X, p.Y, p.R)
	core.Combat.AddTarget(player)

	// add targets
	for i, v := range targets {
		if v.Pos.R == 0 {
			return fmt.Errorf("target cannot have 0 radius (index %v): %v", i, v)
		}
		e := enemy.New(core, v)
		core.Combat.AddTarget(e)
		//s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
	}

	return nil
}

func SetupCharactersInCore(core *core.Core, chars []character.CharacterProfile, initial keys.Char) error {
	if len(chars) > 4 {
		return errors.New("cannot have more than 4 characters per team")
	}
	dup := make(map[keys.Char]bool)

	active := -1
	for _, v := range chars {
		i, err := core.AddChar(v)
		if err != nil {
			return err
		}

		if v.Base.Key == initial {
			core.Player.SetActive(i)
			active = i
		}

		if _, ok := dup[v.Base.Key]; ok {
			return fmt.Errorf("duplicated character %v", v.Base.Key)
		}
		dup[v.Base.Key] = true
	}

	if active == -1 {
		return errors.New("no active character set")
	}

	return nil
}

func SetupResonance(s *core.Core) {
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
					c.AddStatMod(character.StatMod{
						Base:         modifier.NewBase("pyro-res", -1),
						AffectedStat: attributes.NoStat,
						Amount:       f,
					})
				}
			case attributes.Hydro:
				//heal not implemented yet
				f := func() (float64, bool) {
					return 0.3, true
				}
				for _, c := range chars {
					c.AddHealBonusMod(character.HealBonusMod{
						Base:   modifier.NewBase("hydro-res", -1),
						Amount: f,
					})
				}
			case attributes.Cryo:
				val := make([]float64, attributes.EndStatType)
				val[attributes.CR] = .15
				f := func(ae *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					r, ok := t.(*enemy.Enemy)
					if !ok {
						return nil, false
					}
					if r.AuraContains(attributes.Cryo) || r.AuraContains(attributes.Frozen) {
						return val, true
					}
					return nil, false
				}
				for _, c := range chars {
					c.AddAttackMod(character.AttackMod{
						Base:   modifier.NewBase("cyro-res", -1),
						Amount: f,
					})
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
					t, ok := args[0].(*enemy.Enemy)
					if !ok {
						return false
					}
					atk := args[1].(*combat.AttackEvent)
					if s.Player.Shields.PlayerIsShielded() && s.Player.Active() == atk.Info.ActorIndex {
						t.AddResistMod(enemy.ResistMod{
							Base:  modifier.NewBase("geo-res", 15*60),
							Ele:   attributes.Geo,
							Value: -0.2,
						})
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
					c.AddAttackMod(character.AttackMod{
						Base:   modifier.NewBase("geo-res", -1),
						Amount: atkf,
					})
				}

			case attributes.Anemo:
				for _, c := range chars {
					c.AddCooldownMod(character.CooldownMod{
						Base:   modifier.NewBase("anemo-res", -1),
						Amount: func(a action.Action) float64 { return -0.05 },
					})
				}
			}
		}
	}
}

func (s *Simulation) randEnergy() {
	//drop energy
	s.C.Player.DistributeParticle(character.Particle{
		Source: "drop",
		Num:    float64(s.cfg.Energy.Amount),
		Ele:    attributes.NoElement,
	})

	//calculate next
	next := int(s.C.Rand.Float64()*s.cfg.Energy.Mean/5 + s.cfg.Energy.Mean)
	// next := int(-math.Log(1-s.C.Rand.Float64()) / s.cfg.Energy.Lambda)
	s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "rand energy queued - ", fmt.Sprintf("next %v", s.C.F+next)).Write("settings", s.cfg.Energy, "first", next)
	s.C.Tasks.Add(s.randEnergy, next)
}

func (s *Simulation) SetupRandEnergyDrop() {
	//do nothing if none set
	if s.cfg.Energy.Every == 0 {
		return
	}
	//every is given in seconds, so lambda (events per second) is 1 / every
	// s.cfg.Energy.Mean = 1.0 / s.cfg.Energy.Every
	//lambda is per s so we need to scale it to per frame
	// s.cfg.Energy.Mean /= 60

	//convert every to per frame; right now every is in seconds
	s.cfg.Energy.Mean = s.cfg.Energy.Every * 60
	next := int(s.C.Rand.Float64()*s.cfg.Energy.Mean/5 + s.cfg.Energy.Mean)
	// next := int(-math.Log(1-s.C.Rand.Float64()) / s.cfg.Energy.Lambda)
	s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "rand energy started - ", fmt.Sprintf("next %v", s.C.F+next)).Write("settings", s.cfg.Energy, "first", next)
	//start the first round
	s.C.Tasks.Add(s.randEnergy, next)
}
