package simulation

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (s *Simulation) initDetailLog() {
	var sb strings.Builder
	s.stats.ReactionsTriggered = make(map[core.ReactionType]int)
	//add new targets
	s.C.Subscribe(core.OnTargetAdded, func(args ...interface{}) bool {
		t := args[0].(coretype.Target)

		s.C.Log.NewEvent("Target Added", coretype.LogSimEvent, -1, "target_type", t.Type())
		// s.C.Log.Debugw("Target Added", "frame", s.C.F, coretype.LogSimEvent, "target_type", t.Type())

		s.stats.ElementUptime = append(s.stats.ElementUptime, make(map[coretype.EleType]int))

		return false
	}, "sim-new-target-stats")
	//add call backs to track details
	s.C.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		t := args[0].(coretype.Target)

		// No need to pull damage stats for non-enemies
		if t.Type() != coretype.TargettableEnemy {
			return false
		}
		atk := args[1].(*coretype.AttackEvent)

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
		frameBucket := int(s.C.Frame/15) * 15
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
		coretype.OnSwirlHydro:     core.SwirlHydro,
		coretype.OnSwirlCryo:      core.SwirlCryo,
		coretype.OnSwirlElectro:   core.SwirlElectro,
		coretype.OnSwirlPyro:      core.SwirlPyro,
		core.OnCrystallizeCryo:    core.CrystallizeCryo,
		core.OnCrystallizeElectro: core.CrystallizeElectro,
		core.OnCrystallizeHydro:   core.CrystallizeHydro,
		core.OnCrystallizePyro:    core.CrystallizePyro,
	}

	for k, v := range reactions {
		s.C.Subscribe(k, eventSubFunc(v), "reaction-log")
	}

	s.C.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		p := args[0].(core.Particle)
		s.stats.ParticleCount[p.Source] += p.Num
		return false
	}, "particles-log")

	s.C.Subscribe(core.OnEnergyChange, func(args ...interface{}) bool {
		char := args[0].(coretype.Character)
		preEnergy := args[1].(float64)
		amt := args[2].(float64)
		src := args[3].(string)

		temp, ok := s.stats.EnergyDetail[char.Index()][src]
		if !ok {
			temp = [4]float64{0, 0, 0, 0}
		}

		idxToAdd := 0
		if s.C.ActiveChar != char.Index() {
			idxToAdd = 1
		}
		// Total energy gained either on/off-field
		temp[idxToAdd] += char.CurrentEnergy() - preEnergy
		// Total energy wasted (changed into a positive number)
		temp[2+idxToAdd] += -(char.CurrentEnergy() - preEnergy - amt)
		s.stats.EnergyDetail[char.Index()][src] = temp
		return false
	}, "energy-change-log")

	s.C.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		activeChar := s.C.Chars[s.C.ActiveChar]
		s.stats.EnergyWhenBurst[s.C.ActiveChar] = append(s.stats.EnergyWhenBurst[s.C.ActiveChar], activeChar.CurrentEnergy())
		return false
	}, "energy-calc-log")

}
