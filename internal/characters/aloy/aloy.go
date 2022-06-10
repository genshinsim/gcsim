package aloy

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type char struct {
	*tmpl.Character
	coilICDExpiry int
	lastFieldExit int
}

func init() {
	core.RegisterCharFunc(keys.Aloy, NewChar)
	initCancelFrames()
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.Base.Element = attributes.Cryo
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = 4

	c.coilICDExpiry = 0
	c.lastFieldExit = 0

	c.Tags["coil_stacks"] = 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.coilMod()
	c.onExitField()

	return nil
}

func initCancelFrames() {
	initAimedFrames()
	initAttackFrames()
	initBurstFrames()
	initSkillFrames()
}
