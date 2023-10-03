package hydro

import (
	"github.com/genshinsim/gcsim/internal/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	a4Bonus float64
	gender  int
}

func NewChar(gender int) core.NewCharacterFunc {
	return func(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
		c := char{
			gender: gender,
		}
		c.Character = tmpl.NewWithWrapper(s, w)

		c.Base.Element = attributes.Hydro
		c.EnergyMax = 80
		c.BurstCon = 5
		c.SkillCon = 3
		c.NormalHitNum = normalHitNum

		w.Character = &c

		return nil
	}
}

func (c *char) Init() error {
	return nil
}

func (c *char) getSourcewaterDroplets() []combat.Gadget {
	droplets := make([]combat.Gadget, 0)
	for _, g := range c.Core.Combat.Gadgets() {
		_, ok := g.(*common.SourcewaterDroplet)
		if !ok {
			continue
		}
		droplets = append(droplets, g)
	}
	return droplets
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "droplets":
		return len(c.getSourcewaterDroplets()), nil
	default:
		return c.Character.Condition(fields)
	}
}
