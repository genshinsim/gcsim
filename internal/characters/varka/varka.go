package varka

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	conversionElem    attributes.Element
	fourWindsCDStacks int

	a1Buff       float64
	a1Multiplier float64
	a4Stacks     int

	c1Extra int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	hex, ok := p.Params["hexerei"]
	if !ok {
		hex = 1
	}
	c.IsHexerei = (hex != 0)

	w.Character = &c

	c.SetNumCharges(action.ActionSpecialSkill, 2)

	return nil
}

func (c *char) Init() error {
	c.conversionElem = c.getConversionElem(attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)
	c.onExitField()
	c.a1Init()
	c.a4Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	// TODO: Adjust this value based on if windup happened for the NA
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 21
	case info.AnimationYelanN0StartDelay:
		return 4
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) useSpecialSkill() bool {
	return c.StatusIsActive(skillKey) && c.convertToFourWinds()
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	if a == action.ActionSkill && c.useSpecialSkill() {
		if c.Charges(action.ActionSpecialSkill) > 0 {
			return true, action.NoFailure
		}
		return false, action.SkillCD
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		if c.useSpecialSkill() && c.Charges(action.ActionSpecialSkill) > 0 {
			return 0
		}
		return 50
	}
	return c.Character.ActionStam(a, p)
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

func (c *char) ReduceActionCooldown(a action.Action, v int) {
	if a == action.ActionSkill && c.useSpecialSkill() {
		c.Character.ReduceActionCooldown(action.ActionSpecialSkill, v)
	}
	c.Character.ReduceActionCooldown(a, v)
}

func (c *char) ResetActionCooldown(a action.Action) {
	if a == action.ActionSkill && c.useSpecialSkill() {
		c.Character.ResetActionCooldown(action.ActionSpecialSkill)
	}
	c.Character.ResetActionCooldown(a)
}
