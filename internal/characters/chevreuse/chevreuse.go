package chevreuse

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Chevreuse, NewChar)
}

func (c *char) Init() error {
	c.c6StackCounts = [4]int{0, 0, 0, 0}

	// setup overcharged ball
	c.overchargedBallEventSub()

	// make sure to use the same key everywhere so that these passives don't stack
	c.Core.Player.AddStamPercentMod("utility-dash", -1, func(a action.Action) (float64, bool) {
		if a == action.ActionDash && c.CurrentHPRatio() > 0 {
			return -0.2, false
		}
		return 0, false
	})

	// start subscribing for a1/c1
	c.a1()
	c.c1()
	return nil
}

type char struct {
	*tmpl.Character
	onlyPyroElectro bool
	overChargedBall bool
	c4ShotsLeft     int
	c6HealQueued    bool
	c6StackCounts   [4]int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum
	w.Character = &c

	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "overcharged-ball":
		return c.overChargedBall, nil
	default:
		return c.Character.Condition(fields)
	}
}
