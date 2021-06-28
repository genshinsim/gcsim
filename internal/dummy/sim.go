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
	Chars      []def.Character
}

func NewSim(cfg ...func(*Sim)) *Sim {
	s := &Sim{}
	for _, f := range cfg {
		f(s)
	}
	return s
}

func (s *Sim) SwapCD() int                                      { return 0 }
func (s *Sim) Stam() float64                                    { return 0 }
func (s *Sim) Frame() int                                       { return s.F }
func (s *Sim) Flags() def.Flags                                 { return def.Flags{} }
func (s *Sim) CharByName(name string) (def.Character, bool)     { return nil, false }
func (s *Sim) DistributeParticle(p def.Particle)                {}
func (s *Sim) TargetHasDebuff(debuff string, param int) bool    { return false }
func (s *Sim) TargetHasElement(ele def.EleType, param int) bool { return false }
func (s *Sim) OnAttackLanded(t def.Target)                      {}
func (s *Sim) ReactionBonus() float64                           { return 0 }
func (s *Sim) Status(key string) int                            { return 0 }

func (s *Sim) IsShielded() bool                      { return false }
func (s *Sim) GetShield(t def.ShieldType) def.Shield { return nil }
func (s *Sim) Rand() *rand.Rand                      { return s.R }

func (s *Sim) CharByPos(ind int) (def.Character, bool) {
	if ind < 0 || ind >= len(s.Chars) {
		return nil, false
	}
	return s.Chars[ind], true
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
