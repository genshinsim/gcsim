package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (s *Simulation) initChars() error {
	dup := make(map[core.CharKey]bool)
	res := make(map[core.EleType]int)

	count := len(s.cfg.Characters.Profile)

	if count > 4 {
		return fmt.Errorf("more than 4 characters in a team detected")
	}

	// if s.opts.LogDetails {
	s.stats.CharNames = make([]string, count)
	s.stats.CharDetails = make([]CharDetail, 0, count)
	s.stats.DamageByChar = make([]map[string]float64, count)
	s.stats.DamageInstancesByChar = make([]map[string]int, count)
	s.stats.DamageByCharByTargets = make([]map[int]float64, count)
	// s.stats.DamageDetailByTime = make(map[DamageDetails]float64)
	s.stats.DamageDetailByTime = make(map[int]float64)
	s.stats.CharActiveTime = make([]int, count)
	s.stats.AbilUsageCountByChar = make([]map[string]int, count)
	s.stats.ParticleCount = make(map[string]int)
	s.stats.EnergyDetail = make([]map[string][4]float64, count)
	s.stats.EnergyWhenBurst = make([][]float64, count)
	// }

	s.C.ActiveChar = -1
	for i, v := range s.cfg.Characters.Profile {
		//call new char function
		char, err := s.C.AddChar(v)
		if err != nil {
			return err
		}

		if v.Base.Key == s.cfg.Characters.Initial {
			s.C.ActiveChar = i
		}

		if _, ok := dup[v.Base.Key]; ok {
			return fmt.Errorf("duplicated character %v", v.Base.Key)
		}
		dup[v.Base.Key] = true

		//track resonance
		res[char.Ele()]++

		//setup maps
		// if s.opts.LogDetails {
		s.stats.DamageByChar[i] = make(map[string]float64)
		s.stats.DamageInstancesByChar[i] = make(map[string]int)
		s.stats.DamageByCharByTargets[i] = make(map[int]float64)
		s.stats.AbilUsageCountByChar[i] = make(map[string]int)
		s.stats.CharNames[i] = v.Base.Key.String()
		s.stats.EnergyDetail[i] = make(map[string][4]float64)
		s.stats.EnergyWhenBurst[i] = make([]float64, 0, s.cfg.Settings.Duration/12+2)

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
			Sets: v.Sets,
		})

		// }

	}

	if s.C.ActiveChar == -1 {
		return errors.New("no active char set")
	}

	if count == 4 {
		s.initResonance(res)
	}

	return nil
}
