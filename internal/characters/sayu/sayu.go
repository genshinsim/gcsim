package sayu

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Sayu, NewChar)
}

type char struct {
	*tmpl.Character
	eDuration           int
	eAbsorb             attributes.Element
	eAbsorbTag          combat.ICDTag
	absorbCheckLocation combat.AttackPattern
	qTickRadius         float64
	c2Bonus             float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.eDuration = -1
	c.eAbsorb = attributes.NoElement
	c.qTickRadius = 1
	c.c2Bonus = .0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4()
	c.rollAbsorb()
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}
