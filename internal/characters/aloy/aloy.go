package aloy

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Aloy, NewChar)
}

type char struct {
	*tmpl.Character
	coilICDExpiry int
	lastFieldExit int
	//coil related
	coils int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum

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

func (c *char) Condition(k string) int64 {
	switch k {
	case "coil":
		return int64(c.coils)
	default:
		return 0
	}
}
