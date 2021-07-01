package combat

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/def"
)

func (s *Sim) newArtifactSet(name string, c def.Character, count int) error {
	switch name {

	default:
		return fmt.Errorf("invalid artifact set: %v", name)
	}
}
