package shenhe

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Shenhe, NewChar)
}

type char struct {
	*tmpl.Character
	quillcount []int
	c4count    int
	c4expiry   int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSpear
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = character.ZoneLiyue
	c.Base.Element = attributes.Cryo

	c.c4count = 0
	c.c4expiry = 0

	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.quillcount = make([]int, len(c.Core.Player.Chars()))
	c.quillDamageMod()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}
