package mona

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const (
	bubbleKey = "mona-bubble"
	omenKey   = "omen-debuff"
)

func init() {
	core.RegisterCharFunc(keys.Mona, NewChar)
}

type char struct {
	*tmpl.Character
	c2icd int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Hydro
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.c2icd = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstHook()
	c.a4()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}
