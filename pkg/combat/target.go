package combat

import (
	"github.com/genshinsim/gsim/pkg/def"
)

func (s *Sim) Targets() []def.Target {
	return s.targets
}
