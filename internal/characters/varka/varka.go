package varka

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	conversionElem attributes.Element

	// when readyFrame is -1, it means that the skill was used
	fourWindsCharge1ReadyF int
	// when readyFrame is -1, it means that the skill was used
	fourWindsCharge2ReadyF int
	fourWindsCDStacks      int

	a1Buff       float64
	a1Multiplier float64
	a4Stacks     int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 0
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.conversionElem = c.getConversionElem(attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

	c.a1Init()
	c.a4Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	// TODO: Adjust this value based on if windup happened for the NA
	if k == info.AnimationXingqiuN0StartDelay {
		return 12
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	if a == action.ActionSkill && c.StatusIsActive(skillKey) {
		if c.fourWindsCharges() <= 0 {
			return false, action.SkillCD
		}
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		if c.StatusIsActive(skillKey) {
			return 0
		}
		return 50
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Charges(a action.Action) int {
	if a == action.ActionSkill {
		fourWinds := c.fourWindsCharges()
		if fourWinds >= 0 {
			return fourWinds
		}
	}
	return c.Character.Charges(a)
}

func (c *char) Cooldown(a action.Action) int {
	if a == action.ActionSkill && c.StatusIsActive(skillKey) {
		if cd := c.fourWindsCD(); cd >= 0 {
			return cd
		}
	}
	return c.Character.Cooldown(a)
}

func (c *char) ReduceActionCooldown(a action.Action, v int) {
	if a == action.ActionSkill && c.StatusIsActive(skillKey) {
		c.reduceFourWindsCD(v)
	}
	c.Character.ReduceActionCooldown(a, v)
}

func (c *char) ResetActionCooldown(a action.Action) {
	if a == action.ActionSkill && c.StatusIsActive(skillKey) {
		c.resetFourWindsCD()
	}
	c.Character.ResetActionCooldown(a)
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	// can use charge without attack beforehand unlike most of the other sword users
	// if a == action.ActionCharge {
	// 	return nil
	// }
	return c.Character.NextQueueItemIsValid(k, a, p)
}
