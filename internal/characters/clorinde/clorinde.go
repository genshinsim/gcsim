package clorinde

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Clorinde, NewChar)
}

type char struct {
	*tmpl.Character

	normalSCounter int
	a1stacks       *stacks.MultipleRefreshNoRemove
	a1BuffPercent  float64
	a1Cap          float64
	a4stacks       *stacks.MultipleRefreshNoRemove
	a4bonus        []float64
	c6Stacks       int

	// track bol manually skip template
	hpDebt float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = base.SkillDetails.BurstEnergyCost
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.HasArkhe = true

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4Init()
	c.c1()
	c.c4()
	return nil
}

func (c *char) ResetNormalCounter() {
	c.normalSCounter = 0
	c.Character.ResetNormalCounter()
}

func (c *char) AdvanceNormalIndex() {
	if c.StatusIsActive(skillStateKey) {
		c.normalSCounter++
		if c.normalSCounter == skillHitNum {
			c.normalSCounter = 0
		}
		return
	}
	c.Character.AdvanceNormalIndex()
}

func (c *char) NextNormalCounter() int {
	if c.StatusIsActive(skillStateKey) {
		return c.normalSCounter + 1
	}
	return c.Character.NextNormalCounter()
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if a1 window is active is on-field
	if a == action.ActionSkill && c.StatusIsActive(skillStateKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if c.StatusIsActive(skillStateKey) {
		switch k {
		case model.AnimationXingqiuN0StartDelay:
			return 10
		case model.AnimationYelanN0StartDelay:
			return 5
		default:
			return c.Character.AnimationStartDelay(k)
		}
	}

	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 14
	case model.AnimationYelanN0StartDelay:
		return 4
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
