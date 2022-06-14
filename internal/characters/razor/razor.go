package razor

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Razor, NewChar)
}

type char struct {
	*tmpl.Character
	sigils         int
	sigilsDuration int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneMondstadt

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// skill
	c.energySigil()

	// burst
	c.speedBurst()
	c.wolfBurst()
	c.onSwapClearBurst()

	c.a4()

	c.c1()
	c.c2()
	c.c6()

	return nil
}
