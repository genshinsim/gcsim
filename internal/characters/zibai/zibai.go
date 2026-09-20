package zibai

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	radiance   float64
	skillSrc   int
	skillsUsed int
	c6Elev     float64
}

const lunarCrystallizeAbil = " (Lunar-Crystallize)"

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum

	w.Character = &c

	c.Moonsign = 1

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.skillInit()
	c.a1Init()
	c.a4Init()
	c.moonsignInit()
	c.c1Init()
	c.c2Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 18
	}

	if k == info.AnimationYelanN0StartDelay {
		return 7
	}

	return c.Character.AnimationStartDelay(k)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "radiance":
		if c.StatusIsActive(skillKey) {
			return c.radiance, nil
		}
		return 0, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if c.StatusIsActive(skillKey) && a == action.ActionSkill {
		if c.radiance < 70 {
			return false, action.InsufficientStamina
		}
		return true, action.NoFailure
	}

	return c.Character.ActionReady(a, p)
}

func (c *char) ResetNormalCounter() {
	c.c4ResetNormalCount()
}
