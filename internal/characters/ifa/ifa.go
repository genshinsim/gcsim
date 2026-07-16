package ifa

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Ifa, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState   *nightsoul.State
	nightsoulSrc     int
	skillParticleICD bool
	skillLastStamF   int
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

func (c *char) getTeamNightsoul() float64 {
	sum := 0.0
	for _, char := range c.Core.Player.Chars() {
		points, err := char.Condition([]string{"nightsoul", "points"})
		if err != nil {
			continue
		}
		p, ok := points.(float64)
		if !ok {
			continue
		}
		sum += p
	}
	return sum
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if !c.nightsoulState.HasBlessing() {
		return c.Character.ActionStam(a, p)
	}

	return 0
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	if a == action.ActionSwap && c.nightsoulState.Points() > 0 {
		return false, action.SwapCD
	}

	return c.Character.ActionReady(a, p)
}
