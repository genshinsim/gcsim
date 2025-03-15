package lanyan

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var elementAbsorbPriority = []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo}

func (c *char) absorbA1(e *enemy.Enemy) attributes.Element {
	if c.Base.Ascension < 1 {
		return attributes.Anemo
	}
	for _, eleCheck := range elementAbsorbPriority {
		if e.AuraContains(eleCheck) {
			return eleCheck
		}
	}
	return attributes.Anemo
}
