package simulation

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (s *Simulation) initDetailLog() {
	var sb strings.Builder
	s.stats.ReactionsTriggered = make(map[core.ReactionType]int)
	//add new targets
	s.C.Events.Subscribe(core.OnTargetAdded, func(args ...interface{}) bool {
		t := args[0].(core.Target)

		s.C.Log.NewEvent("Target Added", core.LogSimEvent, -1, "target_type", t.Type())
		// s.C.Log.Debugw("Target Added", "frame", s.C.F, core.LogSimEvent, "target_type", t.Type())

		s.stats.ElementUptime = append(s.stats.ElementUptime, make(map[core.EleType]int))

		return false
	}, "sim-new-target-stats")
	//add call backs to track details
	s.C.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)

		// No need to pull damage stats for non-enemies
		if t.Type() != core.TargettableEnemy {
			return false
		}
		atk := args[1].(*core.AttackEvent)

		//skip if do not log
		if atk.Info.DoNotLog {
			return false
		}

		dmg := args[2].(float64)
		sb.Reset()
		sb.WriteString(atk.Info.Abil)
		if atk.Info.Amped {
			if atk.Info.AmpMult == 1.5 {
				sb.WriteString(" [amp: 1.5]")
			} else if atk.Info.AmpMult == 2 {
				sb.WriteString(" [amp: 2.0]")
			}
		}
		s.stats.DamageByChar[atk.Info.ActorIndex][sb.String()] += dmg
		if dmg > 0 {
			s.stats.DamageInstancesByChar[atk.Info.ActorIndex][sb.String()] += 1
		}
		s.stats.DamageByCharByTargets[atk.Info.ActorIndex][t.Index()] += dmg

		// Want to capture information in 0.25s intervals - allows more flexibility in bucketizing
		frameBucket := int(s.C.F/15) * 15
		// details := DamageDetails{
		// 	FrameBucket: frameBucket,
		// 	Char:        atk.Info.ActorIndex,
		// 	Target:      t.Index(),
		// }
		// Go defaults to 0 for map values that don't exist
		s.stats.DamageDetailByTime[frameBucket] += dmg
		return false
	}, "dmg-log")

	eventSubFunc := func(t core.ReactionType) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[t]++
			return false
		}
	}

	var reactions = map[core.EventType]core.ReactionType{
		core.OnOverload:           core.Overload,
		core.OnSuperconduct:       core.Superconduct,
		core.OnMelt:               core.Melt,
		core.OnVaporize:           core.Vaporize,
		core.OnFrozen:             core.Freeze,
		core.OnElectroCharged:     core.ElectroCharged,
		core.OnSwirlHydro:         core.SwirlHydro,
		core.OnSwirlCryo:          core.SwirlCryo,
		core.OnSwirlElectro:       core.SwirlElectro,
		core.OnSwirlPyro:          core.SwirlPyro,
		core.OnCrystallizeCryo:    core.CrystallizeCryo,
		core.OnCrystallizeElectro: core.CrystallizeElectro,
		core.OnCrystallizeHydro:   core.CrystallizeHydro,
		core.OnCrystallizePyro:    core.CrystallizePyro,
	}

	for k, v := range reactions {
		s.C.Events.Subscribe(k, eventSubFunc(v), "reaction-log")
	}

	s.C.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		p := args[0].(core.Particle)
		s.stats.ParticleCount[p.Source] += p.Num
		return false
	}, "particles-log")

	s.C.Events.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		activeChar := s.C.Chars[s.C.ActiveChar]
		s.stats.EnergyWhenBurst[s.C.ActiveChar] = append(s.stats.EnergyWhenBurst[s.C.ActiveChar], activeChar.CurrentEnergy())
		return false
	}, "energy-calc-log")

}
