package combat

import "github.com/genshinsim/gsim/pkg/core"

func (s *Sim) IsShielded() bool { return len(s.shields) > 0 }

func (s *Sim) GetShield(t core.ShieldType) core.Shield {
	for _, v := range s.shields {
		if v.Type() == t {
			return v
		}
	}
	return nil
}

func (s *Sim) AddShieldBonus(f func() float64) {
	s.ShieldBonusFunc = append(s.ShieldBonusFunc, f)
}

func (s *Sim) AddShield(shd core.Shield) {
	//we always assume over write of the same type
	ind := len(s.shields)
	for i, v := range s.shields {
		if v.Type() == shd.Type() {
			ind = i
		}
	}
	if ind != 0 && ind != len(s.shields) {
		s.log.Debugw("shield added", "frame", s.f, "event", core.LogShieldEvent, "frame", s.f, "overwrite", true, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
		s.shields[ind].OnOverwrite()
		s.shields[ind] = shd
	} else {
		s.shields = append(s.shields, shd)
		s.log.Debugw("shield added", "frame", s.f, "event", core.LogShieldEvent, "frame", s.f, "overwrite", false, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
	}

	s.executeEventHook(core.PostShieldHook)
}

func (s *Sim) DamageShields(dmg float64, ele core.EleType) float64 {
	var bonus float64
	//find shield bonuses
	for _, f := range s.ShieldBonusFunc {
		bonus += f()
	}
	min := dmg //min of damage taken
	n := 0
	for _, v := range s.shields {
		taken, ok := v.OnDamage(dmg, ele, bonus)
		if taken < min {
			min = taken
		}
		if ok {
			s.shields[n] = v
			n++
		}
	}
	s.shields = s.shields[:n]
	return min
}

func (s *Sim) tickShields() {
	n := 0
	for _, v := range s.shields {
		if v.Expiry() == s.f {
			v.OnExpire()
			s.log.Debugw("shield expired", "frame", s.f, "event", core.LogShieldEvent, "frame", s.f, "name", v.Desc(), "hp", v.CurrentHP())
		} else {
			s.shields[n] = v
			n++
		}
	}
	s.shields = s.shields[:n]
}
