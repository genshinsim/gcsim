package dummy

import (
	"math/rand"

	"github.com/genshinsim/gsim/pkg/def"
)

type Sim struct {
	F          int
	R          *rand.Rand
	OnDamage   func(ds *def.Snapshot)
	OnShielded func(shd def.Shield)
	OnParticle func(p def.Particle)
	Chars      []def.Character
	Targs      []def.Target
	status     map[string]int

	onAttackLanded []attackLandedHook
	eventHooks     [][]eHook
}

func NewSim(cfg ...func(*Sim)) *Sim {
	s := &Sim{}
	for _, f := range cfg {
		f(s)
	}
	s.onAttackLanded = make([]attackLandedHook, 0, 10)
	s.eventHooks = make([][]eHook, def.EndEventHook)
	s.status = make(map[string]int)
	return s
}

func (s *Sim) Skip(frames int) {
	for i := 0; i < frames; i++ {
		s.F++
		//tick auras and shields firsts
		for _, v := range s.Targs {
			v.AuraTick()
		}

		//then tick each character
		for _, c := range s.Chars {
			c.Tick()
		}

		//then tick each target again
		for _, v := range s.Targs {
			v.Tick()
		}
	}
}

func (s *Sim) ActiveCharIndex() int                             { return 0 }
func (s *Sim) SwapCD() int                                      { return 0 }
func (s *Sim) Stam() float64                                    { return 0 }
func (s *Sim) Frame() int                                       { return s.F }
func (s *Sim) Flags() def.Flags                                 { return def.Flags{} }
func (s *Sim) SetCustomFlag(key string, val int)                {}
func (s *Sim) CharByName(name string) (def.Character, bool)     { return nil, false }
func (s *Sim) TargetHasDebuff(debuff string, param int) bool    { return false }
func (s *Sim) TargetHasElement(ele def.EleType, param int) bool { return false }
func (s *Sim) Targets() []def.Target                            { return s.Targs }
func (s *Sim) ReactionBonus() float64                           { return 0 }
func (s *Sim) HealActive(hp float64)                            {}
func (s *Sim) HealAll(hp float64)                               {}
func (s *Sim) HealIndex(index int, hp float64)                  {}
func (s *Sim) AddIncHealBonus(f func() float64)                 {}
func (s *Sim) AddOnHurt(f func(s def.Sim))                      {}

func (s *Sim) IsShielded() bool                      { return false }
func (s *Sim) GetShield(t def.ShieldType) def.Shield { return nil }
func (s *Sim) Rand() *rand.Rand                      { return s.R }

func (s *Sim) CharByPos(ind int) (def.Character, bool) {
	if ind < 0 || ind >= len(s.Chars) {
		return nil, false
	}
	return s.Chars[ind], true
}

func (s *Sim) Characters() []def.Character {
	return s.Chars
}

func (s *Sim) ApplyDamage(ds *def.Snapshot) {
	if s.OnDamage != nil {
		s.OnDamage(ds)
	}
}

func (s *Sim) AddShield(shd def.Shield) {
	if s.OnShielded != nil {
		s.OnShielded(shd)
	}
}

func (s *Sim) DistributeParticle(p def.Particle) {
	if s.OnParticle != nil {
		s.OnParticle(p)
	}
}

func (s *Sim) AddStatus(key string, dur int) {
	s.status[key] = s.F + dur
}

func (s *Sim) DeleteStatus(key string) {
	delete(s.status, key)
}

func (s *Sim) Status(key string) int {
	f, ok := s.status[key]
	if !ok {
		return 0
	}
	if f > s.F {
		return f - s.F
	}
	return 0
}

type eHook struct {
	f   func(s def.Sim) bool
	key string
	src int
}

//AddHook adds a hook to sim. Hook will be called based on the type of hook
func (s *Sim) AddEventHook(f func(s def.Sim) bool, key string, hook def.EventHookType) {

	a := s.eventHooks[hook]

	//check if override first
	ind := len(a)
	for i, v := range a {
		if v.key == key {
			ind = i
		}
	}
	if ind != 0 && ind != len(a) {
		a[ind] = eHook{
			f:   f,
			key: key,
			src: s.F,
		}
	} else {
		a = append(a, eHook{
			f:   f,
			key: key,
			src: s.F,
		})
	}
	s.eventHooks[hook] = a
}

func (s *Sim) ExecuteEventHook(t def.EventHookType) {
	n := 0
	for _, v := range s.eventHooks[t] {
		if !v.f(s) {
			s.eventHooks[t][n] = v
			n++
		}
	}
	s.eventHooks[t] = s.eventHooks[t][:n]
}

type attackLandedHook struct {
	f   func(t def.Target, ds *def.Snapshot)
	key string
	src int
}

func (s *Sim) OnAttackLanded(t def.Target, ds *def.Snapshot) {
	for _, v := range s.onAttackLanded {
		v.f(t, ds)
	}
}

func (s *Sim) AddOnAttackLanded(f func(t def.Target, ds *def.Snapshot), key string) {

	//check if override first
	ind := -1
	for i, v := range s.onAttackLanded {
		if v.key == key {
			ind = i
		}
	}
	if ind != -1 {
		s.onAttackLanded[ind] = attackLandedHook{
			f:   f,
			key: key,
			src: s.F,
		}
	} else {
		s.onAttackLanded = append(s.onAttackLanded, attackLandedHook{
			f:   f,
			key: key,
			src: s.F,
		})
	}
}
