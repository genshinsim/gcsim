package venti

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
	core.RegisterCharFunc(keys.Venti, NewChar)
}

type char struct {
	*tmpl.Character
	qPos                combat.Point
	qAbsorb             attributes.Element
	absorbCheckLocation combat.AttackPattern
	aiAbsorb            combat.AttackInfo
	snapAbsorb          combat.Snapshot
	c4bonus             []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// C4:
	// When Venti picks up an Elemental Orb or Particle, he receives a 25% Anemo DMG Bonus for 10s.
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}
