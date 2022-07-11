package beidou

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Beidou, NewChar)
}

type char struct {
	*tmpl.Character
	burstAtk *combat.AttackEvent
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstProc()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}
