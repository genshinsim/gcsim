package hutao

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Hutao, NewChar)
}

type char struct {
	*tmpl.Character
	paraParticleICD int
	ppBonus         float64
	tickActive      bool
	applyA1         bool
	c6icd           int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Pyro
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassSpear
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.ppHook()
	c.onExitField()
	c.a4()
	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		if c.Core.Status.Duration("paramita") > 0 && c.Base.Cons >= 1 {
			return 0
		}
		return 25
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.Core.Status.Duration("paramita") > 0 {
		switch ai.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = attributes.Pyro
	}
	return ds
}
