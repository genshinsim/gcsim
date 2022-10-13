package cyno

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Cyno, NewChar)
}

type char struct {
	*tmpl.Character
	burstExtension float64
	c1buff         []float64
	c2counter      int
	c4counter      int
	c6stacks       int
	a1Extended     bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onSwapClearBurst()
	c.a1Extension()

	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Cons >= 6 {
		c.c6()
	}

	// make sure to use the same key everywhere so that these passives don't stack
	c.Core.Player.AddStamPercentMod("utility-dash", -1, func(a action.Action) (float64, bool) {
		if a == action.ActionDash && c.HPCurrent > 0 {
			return -0.2, false
		}
		return 0, false
	})

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		return 25
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "endseer":
		return c.StatusIsActive(a1Key), nil
	default:
		return c.Character.Condition(fields)
	}
}
