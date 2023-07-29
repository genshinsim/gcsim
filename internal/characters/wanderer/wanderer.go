package wanderer

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Wanderer, NewChar)
}

type char struct {
	*tmpl.Character
	skydwellerPoints    int
	maxSkydwellerPoints int
	a1ValidBuffs        []attributes.Element
	a1MaxAbsorb         int
	a4Prob              float64
	c6Count             int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {

	c.maxSkydwellerPoints = 100
	c.a4Prob = 0.16
	c.a1ValidBuffs = []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo}

	c.a1MaxAbsorb = 2
	if c.Base.Cons >= 4 {
		c.a1MaxAbsorb = 1
	}

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if c.StatusIsActive(skillKey) {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

// Overwriting of remaining actions to account for falling state

func (c *char) Walk(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	ai := c.Character.Walk(p)
	ai.Frames = func(next action.Action) int { return delay + ai.Frames(next) }
	ai.AnimationLength = delay + ai.AnimationLength
	ai.CanQueueAfter = delay + ai.CanQueueAfter

	return ai
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "skydweller-points":
		if c.skydwellerPoints <= 0 {
			return 0, nil
		}
		return c.skydwellerPoints, nil
	default:
		return c.Character.Condition(fields)
	}
}
