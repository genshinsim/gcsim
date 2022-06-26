package simulation

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func (s *Simulation) initDetailLog() {
	var sb strings.Builder
	s.stats.ReactionsTriggered = make(map[combat.ReactionType]int)
	s.stats.ElementUptime = make([]map[attributes.Element]int, len(s.C.Combat.Targets()))
	for i := range s.stats.ElementUptime {
		s.stats.ElementUptime[i] = make(map[attributes.Element]int)
	}
	//add call back to track actions executed
	s.C.Events.Subscribe(event.OnActionExec, func(args ...interface{}) bool {

		active := args[0].(int)
		action := args[1].(action.Action)
		s.stats.AbilUsageCountByChar[active][action.String()]++
		return false
	}, "sim-abil-usage")
	//add new targets
	s.C.Events.Subscribe(event.OnTargetAdded, func(args ...interface{}) bool {
		t := args[0].(combat.Target)

		s.C.Log.NewEvent("Target Added", glog.LogSimEvent, -1, "target_type", t.Type())
		// s.C.Log.Debugw("Target Added", "frame", s.C.F, core.LogSimEvent, "target_type", t.Type())

		s.stats.ElementUptime = append(s.stats.ElementUptime, make(map[attributes.Element]int))

		return false
	}, "sim-new-target-stats")
	//add call backs to track details
	s.C.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		t := args[0].(combat.Target)

		// No need to pull damage stats for non-enemies
		if t.Type() != combat.TargettableEnemy {
			return false
		}
		atk := args[1].(*combat.AttackEvent)

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

	eventSubFunc := func(t combat.ReactionType) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[t]++
			return false
		}
	}

	var reactions = map[event.Event]combat.ReactionType{
		event.OnOverload:           combat.Overload,
		event.OnSuperconduct:       combat.Superconduct,
		event.OnMelt:               combat.Melt,
		event.OnVaporize:           combat.Vaporize,
		event.OnFrozen:             combat.Freeze,
		event.OnElectroCharged:     combat.ElectroCharged,
		event.OnSwirlHydro:         combat.SwirlHydro,
		event.OnSwirlCryo:          combat.SwirlCryo,
		event.OnSwirlElectro:       combat.SwirlElectro,
		event.OnSwirlPyro:          combat.SwirlPyro,
		event.OnCrystallizeCryo:    combat.CrystallizeCryo,
		event.OnCrystallizeElectro: combat.CrystallizeElectro,
		event.OnCrystallizeHydro:   combat.CrystallizeHydro,
		event.OnCrystallizePyro:    combat.CrystallizePyro,
	}

	for k, v := range reactions {
		s.C.Events.Subscribe(k, eventSubFunc(v), "reaction-log")
	}

	s.C.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		p := args[0].(character.Particle)
		s.stats.ParticleCount[p.Source] += p.Num
		return false
	}, "particles-log")

	s.C.Events.Subscribe(event.OnEnergyChange, func(args ...interface{}) bool {
		char := args[0].(*character.CharWrapper)
		preEnergy := args[1].(float64)
		amt := args[2].(float64)
		src := args[3].(string)

		temp, ok := s.stats.EnergyDetail[char.Index][src]
		if !ok {
			temp = [4]float64{0, 0, 0, 0}
		}

		idxToAdd := 0
		if s.C.Player.Active() != char.Index {
			idxToAdd = 1
		}
		// Total energy gained either on/off-field
		temp[idxToAdd] += char.Energy - preEnergy
		// Total energy wasted (changed into a positive number)
		temp[2+idxToAdd] += -(char.Energy - preEnergy - amt)
		s.stats.EnergyDetail[char.Index][src] = temp
		return false
	}, "energy-change-log")

	s.C.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		activeChar := s.C.Player.ActiveChar()
		s.stats.EnergyWhenBurst[activeChar.Index] = append(s.stats.EnergyWhenBurst[activeChar.Index], activeChar.Energy)
		return false
	}, "energy-calc-log")

}

func (s *Simulation) initTeamStats() {
	count := len(s.cfg.Characters)
	s.stats.CharNames = make([]string, count)
	s.stats.CharDetails = make([]CharDetail, 0, count)
	s.stats.DamageByChar = make([]map[string]float64, count)
	s.stats.DamageInstancesByChar = make([]map[string]int, count)
	s.stats.DamageByCharByTargets = make([]map[int]float64, count)
	// s.stats.DamageDetailByTime = make(map[DamageDetails]float64)
	s.stats.DamageDetailByTime = make(map[int]float64)
	s.stats.CharActiveTime = make([]int, count)
	s.stats.AbilUsageCountByChar = make([]map[string]int, count)
	s.stats.ParticleCount = make(map[string]float64)
	s.stats.EnergyDetail = make([]map[string][4]float64, count)
	s.stats.EnergyWhenBurst = make([][]float64, count)

	for i, v := range s.cfg.Characters {
		// if s.opts.LogDetails {
		s.stats.DamageByChar[i] = make(map[string]float64)
		s.stats.DamageInstancesByChar[i] = make(map[string]int)
		s.stats.DamageByCharByTargets[i] = make(map[int]float64)
		s.stats.AbilUsageCountByChar[i] = make(map[string]int)
		s.stats.CharNames[i] = v.Base.Key.String()
		s.stats.EnergyDetail[i] = make(map[string][4]float64)
		s.stats.EnergyWhenBurst[i] = make([]float64, 0, int(s.cfg.Settings.Duration/12+2))

		//convert set to string
		m := make(map[string]int)
		for k, v := range v.Sets {
			m[k.String()] = v
		}

		//log the character data
		s.stats.CharDetails = append(s.stats.CharDetails, CharDetail{
			Name:     v.Base.Key.String(),
			Level:    v.Base.Level,
			MaxLevel: v.Base.MaxLevel,
			Cons:     v.Base.Cons,
			Weapon: WeaponDetail{
				Refine:   v.Weapon.Refine,
				Level:    v.Weapon.Level,
				MaxLevel: v.Weapon.MaxLevel,
			},
			Talents: TalentDetail{
				Attack: v.Talents.Attack,
				Skill:  v.Talents.Skill,
				Burst:  v.Talents.Burst,
			},
			Sets: m,
		})

	}
}
