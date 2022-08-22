package heizou

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Heizou, NewChar)
}

type char struct {
	*tmpl.Character
	decStack            int
	infuseCheckLocation combat.AttackPattern
	c1buff              []float64
	a4buff              []float64
	burstTaggedCount    int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	c.infuseCheckLocation = combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy)

	w.Character = &c

	return nil
}

func (c *char) Init() error {

	c.a1()

	c.a4buff = make([]float64, attributes.EndStatType)
	c.a4buff[attributes.EM] = 80

	// make sure to use the same key everywhere so that these passives don't stack
	c.Core.Player.AddStamPercentMod("utility-dash", -1, func(a action.Action) (float64, bool) {
		if a == action.ActionDash && c.HPCurrent > 0 {
			return -0.2, false
		}
		return 0, false
	})

	if c.Base.Cons >= 1 {
		c.c1buff = make([]float64, attributes.EndStatType)
		c.c1buff[attributes.AtkSpd] = .15
		c.c1()
	}

	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		return 25
	}
	return c.Character.ActionStam(a, p)
}
