package heizou

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Heizou, NewChar)
}

type char struct {
	*tmpl.Character
	decStack         int
	c1buff           []float64
	a4Buff           []float64
	burstTaggedCount int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()

	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4Buff[attributes.EM] = 80

	// make sure to use the same key everywhere so that these passives don't stack
	c.Core.Player.AddStamPercentMod("utility-dash", -1, func(a action.Action) (float64, bool) {
		if a == action.ActionDash && c.CurrentHPRatio() > 0 {
			return -0.2, false
		}
		return 0, false
	})

	if c.Base.Cons >= 1 {
		c.c1buff = make([]float64, attributes.EndStatType)
		c.c1buff[attributes.AtkSpd] = .15
		c.c1()
	}

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		return 25
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "declension":
		return c.decStack, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) NextQueueItemIsValid(a action.Action, p map[string]int) error {
	// cannot use charge without attack beforehand unlike most of the other catalyst users
	if a == action.ActionCharge && c.Core.Player.LastAction.Type != action.ActionAttack {
		return player.ErrInvalidChargeAction
	}
	return c.Character.NextQueueItemIsValid(a, p)
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 10
	}
	return c.Character.AnimationStartDelay(k)
}
