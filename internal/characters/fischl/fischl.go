package fischl

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Fischl, NewChar)
}

type char struct {
	*tmpl.Character
	//field use for calculating oz damage
	ozSnapshot    combat.AttackEvent
	ozSource      int //keep tracks of source of oz aka resets
	ozActiveUntil int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum

	c.ozSource = -1
	c.ozActiveUntil = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "oz":
		if c.ozActiveUntil <= c.Core.F {
			return false, nil
		}
		return (c.ozActiveUntil - c.Core.F), nil
	case "oz-source":
		return c.ozSource, nil
	default:
		return c.Character.Condition(fields)
	}
}
