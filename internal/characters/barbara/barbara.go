package barbara

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
	core.RegisterCharFunc(keys.Barbara, NewChar)
}

type char struct {
	*tmpl.Character
	c6icd      int
	skillInitF int
	// burstBuffExpiry   int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.Base.Element = attributes.Hydro
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 4
	c.CharZone = character.ZoneMondstadt

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.a1()

	if c.Base.Cons >= 1 {
		c.c1(1)
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}

func (c *char) a1() {
	c.Core.AddStamMod(func(a action.Action) (float64, bool) { // @srl does this activate for the active char?
		if c.Core.Status.Duration("barbskill") >= 0 {
			return -0.12, false
		}
		return 0, false
	}, "barb-a1-stam")
}
