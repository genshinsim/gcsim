package combat

import "github.com/genshinsim/gsim/pkg/def"

func (s *Sim) Flags() def.Flags { return def.Flags{} }

func (s *Sim) SetCustomFlag(key string, val int) {
	s.flags.Custom[key] = val
}
