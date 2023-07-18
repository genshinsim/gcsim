package anemo

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

type char struct {
	*tmpl.Character
	qAbsorb              attributes.Element
	qICDTag              attacks.ICDTag
	qAbsorbCheckLocation combat.AttackPattern
	eAbsorb              attributes.Element
	eICDTag              attacks.ICDTag
	eAbsorbCheckLocation combat.AttackPattern
	gender               int
}

func NewChar(gender int) core.NewCharacterFunc {
	return func(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
		c := char{
			gender: gender,
		}
		c.Character = tmpl.NewWithWrapper(s, w)

		c.Base.Element = attributes.Anemo
		c.EnergyMax = 60
		c.BurstCon = 3
		c.SkillCon = 5
		c.NormalHitNum = normalHitNum

		w.Character = &c

		return nil
	}
}

func (c *char) Init() error {
	c.a4()
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}
