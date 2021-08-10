package dummy

import (
	"math/rand"

	"github.com/genshinsim/gsim/pkg/core"
)

type Sim struct {
	F          int
	R          *rand.Rand
	OnDamage   func(ds *core.Snapshot)
	OnShielded func(shd core.Shield)
	OnParticle func(p core.Particle)
	Chars      []core.Character
	Targs      []core.Target
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
	s.eventHooks = make([][]eHook, core.EndEventHook)
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

func (s *Sim) ActiveCharIndex() int                                                  { return 0 }
func (s *Sim) SwapCD() int                                                           { return 0 }
func (s *Sim) RestoreStam(v float64)                                                 {}
func (s *Sim) Stam() float64                                                         { return 0 }
func (s *Sim) Frame() int                                                              { return s.F }
func (s *Sim) Flags() core.Flags                                                       { return core.Flags{} }
func (s *Sim) SetCustomFlag(key string, val int)                                       {}
func (s *Sim) GetCustomFlag(key string) (int, bool)                                    { return 0, false }
func (s *Sim) CharByName(name string) (core.Character, bool)                           { return nil, false }
func (s *Sim) TargetHasResMod(debuff string, param int) bool                           { return false }
func (s *Sim) TargetHasDefMod(debuff string, param int) bool                           { return false }
func (s *Sim) TargetHasElement(ele core.EleType, param int) bool                       { return false }
func (s *Sim) Targets() []core.Target                                                  { return s.Targs }
func (s *Sim) AddOnAmpReaction(f func(t core.Target, ds *core.Snapshot), key string)   {}
func (s *Sim) OnAmpReaction(t core.Target, ds *core.Snapshot)                          {}
func (s *Sim) AddOnTransReaction(f func(t core.Target, ds *core.Snapshot), key string) {}
func (s *Sim) OnTransReaction(t core.Target, ds *core.Snapshot)                        {}
func (s *Sim) AddOnReaction(f func(t core.Target, ds *core.Snapshot), key string)      {}
func (s *Sim) OnReaction(t core.Target, ds *core.Snapshot)                             {}
func (s *Sim) HealActive(hp float64)                                                   {}
func (s *Sim) HealAll(hp float64)                                                    {}
func (s *Sim) HealAllPercent(percent float64)                                        {}
func (s *Sim) HealIndex(index int, hp float64)                                       {}
func (s *Sim) AddIncHealBonus(f func() float64)                                        {}
func (s *Sim) AddOnHurt(f func(s core.Sim))                                            {}
func (s *Sim) IsShielded() bool                                                        { return false }
func (s *Sim) GetShield(t core.ShieldType) core.Shield                                 { return nil }
func (s *Sim) AddShieldBonus(f func() float64)                                         {}
func (s *Sim) Rand() *rand.Rand                                                      { return s.R }
func (s *Sim) AddInitHook(f func())                                                    {}
func (s *Sim) OnTargetDefeated(t core.Target)                                          {}
func (s *Sim) AddOnTargetDefeated(f func(t core.Target), key string)                   {}
func (s *Sim) ActiveDuration() int                                                     { return 0 }
func (s *Sim) NewConstruct(c core.Construct, refresh bool)                             {}
func (s *Sim) NewNoLimitCons(c core.Construct, refresh bool)                           {}
func (s *Sim) ConstructCount() int                                                     { return 0 }
func (s *Sim) ConstructCountType(t core.GeoConstructType) int                          { return 0 }
func (s *Sim) HasConstruct(key int) bool                                               { return false }
func (s *Sim) Destroy(key int) bool                                                    { return false }
func (s *Sim) AddStamMod(f func(a core.ActionType) float64)                            {}

func (s *Sim) CharByPos(ind int) (core.Character, bool) {
	if ind < 0 || ind >= len(s.Chars) {
		return nil, false
	}
	return s.Chars[ind], true
}

func (s *Sim) Characters() []core.Character {
	return s.Chars
}

func (s *Sim) ApplyDamage(ds *core.Snapshot) {
	if s.OnDamage != nil {
		s.OnDamage(ds)
	}
}

func (s *Sim) AddShield(shd core.Shield) {
	if s.OnShielded != nil {
		s.OnShielded(shd)
	}
}

func (s *Sim) DistributeParticle(p core.Particle) {
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
	f   func(s core.Sim) bool
	key string
	src int
}

//AddHook adds a hook to sim. Hook will be called based on the type of hook
func (s *Sim) AddEventHook(f func(s core.Sim) bool, key string, hook core.EventHookType) {

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

func (s *Sim) ExecuteEventHook(t core.EventHookType) {
	n := 0
	for _, v := range s.eventHooks[t] {
		if !v.f(s) {
			s.eventHooks[t][n] = v
			n++
		}
	}
	s.eventHooks[t] = s.eventHooks[t][:n]
}

func (s *Sim) AddOnAttackWillLand(f func(t core.Target, ds *core.Snapshot), key string) {}
func (s *Sim) OnAttackWillLand(t core.Target, ds *core.Snapshot)                        {}

type attackLandedHook struct {
	f   func(t core.Target, ds *core.Snapshot, dmg float64, crit bool)
	key string
	src int
}

func (s *Sim) OnAttackLanded(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
	for _, v := range s.onAttackLanded {
		v.f(t, ds, dmg, crit)
	}
}

func (s *Sim) AddOnAttackLanded(f func(t core.Target, ds *core.Snapshot, dmg float64, crit bool), key string) {

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
