package anemo

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type Traveler struct {
	*tmpl.Character
	qAbsorb              attributes.Element
	qICDTag              attacks.ICDTag
	qAbsorbCheckLocation combat.AttackPattern
	eAbsorb              attributes.Element
	eICDTag              attacks.ICDTag
	eAbsorbCheckLocation combat.AttackPattern
	gender               int
}

func NewTraveler(s *core.Core, w *character.CharWrapper, p info.CharacterProfile, gender int) (*Traveler, error) {
	c := Traveler{
		gender: gender,
	}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Atk += common.TravelerBaseAtkIncrease(p)
	c.Base.Element = attributes.Anemo
	c.EnergyMax = 60
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	return &c, nil
}

func (c *Traveler) Init() error {
	c.a4()
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}
