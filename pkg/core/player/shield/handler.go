package shield

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

//Handler keeps track of all the active shields
//Note that there's no need to distinguish between characters here since the shields are shared
//However we do care which is active since only active is shielded
type Handler struct {
	shields []Shield
	log     glog.Logger
	events  event.Eventter
	f       *int

	shieldBonusMods []shieldBonusMod
}

func New(f *int, log glog.Logger, events event.Eventter) *Handler {
	h := &Handler{
		shields:         make([]Shield, 0, EndShieldType),
		log:             log,
		events:          events,
		f:               f,
		shieldBonusMods: make([]shieldBonusMod, 0, 5),
	}
	return h
}

func (s *Handler) Count() int { return len(s.shields) }

func (h *Handler) PlayerIsShielded() bool {
	return len(h.shields) > 0
}

// func (s *Handler) IsShielded(char int) bool {
// 	return len(s.shields) > 0 && char == s.core.ActiveChar
// }

func (s *Handler) Get(t ShieldType) Shield {
	for _, v := range s.shields {
		if v.Type() == t {
			return v
		}
	}
	return nil
}

func (s *Handler) Add(shd Shield) {
	//we always assume over write of the same type
	ind := -1
	for i, v := range s.shields {
		if v.Type() == shd.Type() {
			ind = i
		}
	}
	if ind > -1 {
		s.log.NewEvent("shield overridden", glog.LogShieldEvent, -1, "overwrite", true, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
		s.shields[ind].OnOverwrite()
		s.shields[ind] = shd
	} else {
		s.shields = append(s.shields, shd)
		s.log.NewEvent("shield added", glog.LogShieldEvent, -1, "overwrite", false, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
	}
	s.events.Emit(event.OnShielded, shd)
}

func (s *Handler) OnDamage(char int, dmg float64, ele attributes.Element) float64 {
	//find shield bonuses
	bonus := s.ShieldBonus()
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

func (s *Handler) Tick() {
	n := 0
	for _, v := range s.shields {
		if v.Expiry() == *s.f {
			v.OnExpire()
			s.log.NewEvent("shield expired", glog.LogShieldEvent, -1, "name", v.Desc(), "hp", v.CurrentHP())
		} else {
			s.shields[n] = v
			n++
		}
	}
	s.shields = s.shields[:n]
}
