package aloy

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Aloy, NewChar)
}

type char struct {
	*tmpl.Character
	coilICDExpiry int
	lastFieldExit int
	// coil related
	coils int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	c.coilICDExpiry = 0
	c.lastFieldExit = 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.coilMod()
	c.onExitField()

	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "coil":
		return c.coils, nil
	default:
		return c.Character.Condition(fields)
	}
}
