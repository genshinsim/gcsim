package combat

import (
	"github.com/genshinsim/gsim/pkg/def"
)

func (s *Sim) CharByPos(ind int) (def.Character, bool) {
	if ind < 0 || ind >= len(s.chars) {
		return nil, false
	}
	return s.chars[ind], true
}

func (s *Sim) CharByName(name string) (def.Character, bool) {
	ind, ok := s.charPos[name]
	if !ok {
		return nil, false
	}
	return s.chars[ind], true
}

func (s *Sim) DistributeParticle(p def.Particle) {
	for i, c := range s.chars {
		c.ReceiveParticle(p, s.active == i, len(s.chars))
	}
}
