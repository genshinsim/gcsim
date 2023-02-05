package razor

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Razor, NewChar)
}

type char struct {
	*tmpl.Character
	sigils          int
	skillSigilBonus []float64
	a4Bonus         []float64
	c1bonus         []float64
	c2bonus         []float64
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
	// skill
	c.skillSigilBonus = make([]float64, attributes.EndStatType)
	c.energySigil()

	// burst
	c.speedBurst()
	c.onSwapClearBurst()

	c.a4()

	// make sure to use the same key everywhere so that these passives don't stack
	c.Core.Player.AddStamPercentMod("utility-dash", -1, func(a action.Action) (float64, bool) {
		if a == action.ActionDash && c.HPCurrent > 0 {
			return -0.2, false
		}
		return 0, false
	})
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}
