package odette

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	danceDoubleSrc int
	a1StacksSelf   int
	a1StacksOthers int
	a1Src          int
	c2Src          int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.stellarInit()
	c.a1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 12
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "a1-stacks-self":
		if !c.StatusIsActive(danceDoubleKey) {
			return 0, nil
		}
		return c.a1StacksSelf, nil
	case "a1-stacks-other":
		if !c.StatusIsActive(danceDoubleKey) {
			return 0, nil
		}
		return c.a1StacksOthers, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) useSpecialSkill() bool {
	return c.StatusIsActive(danceDoubleKey) && !c.StatusIsActive(danceDoubleUpgradeKey) && c.StatusIsActive(skillRecastKey)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.useSpecialSkill() {
		if c.Character.Charges(action.ActionSpecialSkill) > 0 {
			return true, action.NoFailure
		}
		return false, action.SkillCD
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) Charges(a action.Action) int {
	if a == action.ActionSkill && c.useSpecialSkill() {
		return c.Character.Charges(action.ActionSpecialSkill)
	}
	return c.Character.Charges(a)
}

func (c *char) Cooldown(a action.Action) int {
	if a == action.ActionSkill && c.useSpecialSkill() {
		return c.Character.Cooldown(action.ActionSpecialSkill)
	}
	return c.Character.Cooldown(a)
}

func (c *char) ResetActionCooldown(a action.Action) {
	if a == action.ActionSkill && c.useSpecialSkill() {
		c.Character.ResetActionCooldown(action.ActionSpecialSkill)
		return
	}
	c.Character.ResetActionCooldown(a)
}

func (c *char) ReduceActionCooldown(a action.Action, v int) {
	if a == action.ActionSkill && c.useSpecialSkill() {
		c.Character.ReduceActionCooldown(action.ActionSpecialSkill, v)
		return
	}
	c.Character.ReduceActionCooldown(a, v)
}
