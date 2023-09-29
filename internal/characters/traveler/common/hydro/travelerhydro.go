package hydro

import (
	"github.com/genshinsim/gcsim/internal/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	droplets []*common.SourcewaterDroplet
	a4Bonus  float64
	gender   int
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
	c.droplets = make([]*common.SourcewaterDroplet, 0)

	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "droplets":
		return len(c.droplets), nil
	default:
		return c.Character.Condition(fields)
	}
}
