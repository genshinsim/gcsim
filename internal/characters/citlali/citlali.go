package citlali

// Citlali's Frostfall Storm applies once every 1.5s
// Initial E hit has no ICD.
// Initial burst has no ICD.
// Spiritvessel Skull ICD is default.
// C4 has no ICD.
// NA is default.
// CA has no ICD.

// 5 particles on initial E hit

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Citlali, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState   *nightsoul.State
	itzpapaSrc       int
	skillShield      *shd
	numStellarBlades int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = 3
	c.SkillCon = 3
	c.BurstCon = 5
	c.HasArkhe = false

	w.Character = &c

	c.nightsoulState = nightsoul.New(s, w)
	c.nightsoulState.MaxPoints = 80 // TODO: the REAL one

	c.itzpapaSrc = -1

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4()

	c.c1()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	default:
		return c.Character.Condition(fields)
	}
}
