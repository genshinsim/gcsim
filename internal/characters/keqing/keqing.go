package keqing

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
	core.RegisterCharFunc(keys.Keqing, NewChar)
}

type char struct {
	*tmpl.Character
	c2ICD int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassSword
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) bool {
	// check if stiletto is on-field
	if a == action.ActionSkill && c.Core.Status.Duration(stilettoKey) > 0 {
		return true
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		return 25
	}
	return c.Character.ActionStam(a, p)
}
