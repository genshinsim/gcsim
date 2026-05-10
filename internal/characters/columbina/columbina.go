package columbina

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Columbina, NewChar)
}

type char struct {
	*tmpl.Character
	skillSrc            int
	gravity             [3]float64
	gravityTask         bool
	gravityLastReaction info.ReactionType
	burstSrc            int
	burstArea           info.AttackPattern
	a1Stacks            int
	a1Buff              []float64
	a4MoondewCount      int
	c2Buff              []float64
	c2LCBuff            []float64
	c2LBBuff            []float64
	c2LCrBuff           []float64
	c6Buff              []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.Moonsign = 1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillInit()
	c.a1Init()
	c.a4Init()
	c.moonsignInit()
	c.consElevationInit()
	c.c2Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 10
	case info.AnimationYelanN0StartDelay:
		return 10
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge && c.Core.Player.Dew() > 0 {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "gravity":
		return c.totalGravity(), nil
	case "a1-stacks":
		if !c.StatModIsActive(a1Key) {
			return 0, nil
		}
		return c.a1Stacks, nil
	default:
		return c.Character.Condition(fields)
	}
}
