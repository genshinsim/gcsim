package ifa

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Ifa, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState      *nightsoul.State
	nightsoulSrc        int
	skillParticleICD    bool
	teamNightsoulPoints []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.nightsoulState = nightsoul.New(c.Core, c.CharWrapper)
	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.nightsoulTeamTracking()
	c.a1Init()
	c.a4Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 7
	case info.AnimationYelanN0StartDelay:
		return 4
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) nightsoulTeamTracking() {
	for range c.Core.Player.Chars() {
		c.teamNightsoulPoints = append(c.teamNightsoulPoints, 0.0)
	}

	c.Core.Events.Subscribe(event.OnNightsoulRemove, func(args ...any) {
		amount := args[1].(float64)
		if amount < 0.0000001 {
			return
		}
		idx := args[0].(int)

		c.teamNightsoulPoints[idx] -= amount
	}, "ifa-subtract-team-nightsoul")

	c.Core.Events.Subscribe(event.OnNightsoulAdd, func(args ...any) {
		amount := args[1].(float64)
		if amount < 0.0000001 {
			return
		}
		idx := args[0].(int)

		c.teamNightsoulPoints[idx] += amount
	}, "ifa-add-team-nightsoul")
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if !c.nightsoulState.HasBlessing() {
		return c.Character.ActionStam(a, p)
	}

	return 0
}
