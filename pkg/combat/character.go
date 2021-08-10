package combat

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func (s *Sim) CharByPos(ind int) (core.Character, bool) {
	if ind < 0 || ind >= len(s.chars) {
		return nil, false
	}
	return s.chars[ind], true
}

func (s *Sim) CharByName(name string) (core.Character, bool) {
	ind, ok := s.charPos[name]
	if !ok {
		return nil, false
	}
	return s.chars[ind], true
}
func (s *Sim) ActiveCharIndex() int { return s.active }

func (s *Sim) DistributeParticle(p core.Particle) {
	for i, c := range s.chars {
		c.ReceiveParticle(p, s.active == i, len(s.chars))
	}
	s.executeEventHook(core.PostParticleHook)
}

func (s *Sim) Characters() []core.Character {
	return s.chars
}

func (s *Sim) ActiveDuration() int {
	return s.charActiveDuration
}
