package shield

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

type ShieldCtrl struct {
	shields   []core.Shield
	bonusFunc []func() float64
	core      *core.Core
}

func NewCtrl(c *core.Core) *ShieldCtrl {
	return &ShieldCtrl{
		shields:   make([]core.Shield, 0, core.EndShieldType),
		bonusFunc: make([]func() float64, 0, 10),
		core:      c,
	}
}

func (s *ShieldCtrl) Count() int { return len(s.shields) }

func (s *ShieldCtrl) IsShielded(char int) bool {
	return len(s.shields) > 0 && char == s.core.ActiveChar
}

func (s *ShieldCtrl) Get(t core.ShieldType) core.Shield {
	for _, v := range s.shields {
		if v.Type() == t {
			return v
		}
	}
	return nil
}

func (s *ShieldCtrl) AddBonus(f func() float64) {
	s.bonusFunc = append(s.bonusFunc, f)
}

func (s *ShieldCtrl) Add(shd core.Shield) {
	//we always assume over write of the same type
	ind := -1
	for i, v := range s.shields {
		if v.Type() == shd.Type() {
			ind = i
		}
	}
	if ind > -1 {
		s.coretype.Log.NewEvent("shield overridden", coretype.LogShieldEvent, -1, "overwrite", true, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
		s.shields[ind].OnOverwrite()
		s.shields[ind] = shd
	} else {
		s.shields = append(s.shields, shd)
		s.coretype.Log.NewEvent("shield added", coretype.LogShieldEvent, -1, "overwrite", false, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
	}
	s.core.Events.Emit(coretype.OnShielded, shd)
}

func (s *ShieldCtrl) OnDamage(dmg float64, ele coretype.EleType) float64 {
	var bonus float64
	//find shield bonuses
	for _, f := range s.bonusFunc {
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

func (s *ShieldCtrl) Tick() {
	n := 0
	for _, v := range s.shields {
		if v.Expiry() == s.core.Frame {
			v.OnExpire()
			s.coretype.Log.NewEvent("shield expired", coretype.LogShieldEvent, -1, "name", v.Desc(), "hp", v.CurrentHP())
		} else {
			s.shields[n] = v
			n++
		}
	}
	s.shields = s.shields[:n]
}
