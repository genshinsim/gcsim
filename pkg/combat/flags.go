package combat

import "github.com/genshinsim/gsim/pkg/core"

func (s *Sim) Flags() core.Flags { return s.flags }

func (s *Sim) SetCustomFlag(key string, val int) {
	s.flags.Custom[key] = val
}

func (s *Sim) GetCustomFlag(key string) (int, bool) {
	val, ok := s.flags.Custom[key]
	return val, ok
}
