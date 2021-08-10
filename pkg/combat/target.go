package combat

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func (s *Sim) Targets() []core.Target {
	return s.targets
}
