package venti

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	qPos                info.Point
	qAbsorb             attributes.Element
	absorbCheckLocation info.AttackPattern
	aiAbsorb            info.AttackInfo
	snapAbsorb          info.Snapshot
	burstSrc            int
	c4bonus             []float64
	burstHexMult        float64
	burstHexSrc         int
	c6HexBonus          []float64
	c6ResistModRange    []string
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	hex, ok := p.Params["hexerei"]
	if !ok {
		// default hexerei is enabled
		hex = 1
	}
	c.IsHexerei = (hex != 0)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstHexMult = 1
	// C4:
	// When Venti picks up an Elemental Orb or Particle, he receives a 25% Anemo DMG Bonus for 10s.
	c.c4()
	c.c6HexereiInit()
	c.hexInit()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 9
	}
	return c.Character.AnimationStartDelay(k)
}
