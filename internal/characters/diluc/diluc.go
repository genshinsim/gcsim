package diluc

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
	core.RegisterCharFunc(keys.Diluc, NewChar)
}

type char struct {
	*tmpl.Character
	eCounter int
	eWindow  int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Pyro
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum

	c.eCounter = 0
	c.eWindow = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 1 && c.Core.Combat.DamageMode {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) bool {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.Core.F < c.eWindow {
		return true
	}
	return c.Character.ActionReady(a, p)
}
