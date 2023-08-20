package wriothesley

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Wriothesley, NewChar)
}

type char struct {
	*tmpl.Character
	a1ICD     int
	a1HPRatio float64
	a1Buff    []float64
	a4Stack   int
	c1Proc    bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.NormalCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Ascension >= 1 {
		c.a1()
	}
	if c.Base.Ascension >= 4 {
		c.a4()
	}

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge && c.StatModIsActive(a1Status) {
		return 0
	}

	return c.Character.ActionStam(a, p)
}
