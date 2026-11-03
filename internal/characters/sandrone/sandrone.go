package sandrone

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const normalHitNum = 3

type char struct {
	*tmpl.Character

	// Fagio / Decoding Power state
	decodingPower float64
	lostPower     float64
	fagioMode     fagioMode
	tickerActive  bool
	overdriveSrc  int

	// Charged Attack state
	chargeSrc int
	beamCount int

	// A1 Refined Tactics
	tactics int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = burstEnergy
	c.NormalHitNum = normalHitNum
	// C3 raises Formule Phenomenale: Self-Evident Proposition (Normal Attack), C5 raises Q.E.D. (Burst)
	c.NormalCon = 3
	c.BurstCon = 5

	c.chargeSrc = -1
	c.overdriveSrc = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.stellarInit()
	c.a4Init()
	c.c1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "decoding-power":
		return int(c.decodingPower), nil
	case "fagio-mode":
		return int(c.fagioMode), nil
	case "refined-tactics":
		return c.refinedTactics(), nil
	case "hovermech":
		return c.StatusIsActive(hovermechKey), nil
	default:
		return c.Character.Condition(fields)
	}
}
