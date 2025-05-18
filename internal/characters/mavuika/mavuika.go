package mavuika

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

type SkillState int

const (
	ring SkillState = iota
	bike
	bikeCDKey = "flamestrider-charge"
)

type char struct {
	*tmpl.Character
	fightingSpirit     float64
	nightsoulState     *nightsoul.State
	nightsoulSrc       int
	armamentState      SkillState
	ringSrc            int
	burstStacks        float64
	a4src              int
	a4stacks           int
	a4buff             []float64
	c1buff             []float64
	c6Src              int
	savedNormalCounter int
	caState            ChargeState
	canBikePlunge      bool
}

func init() {
	core.RegisterCharFunc(keys.Mavuika, NewChar)
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)

	t.CharWrapper = w
	c.Character = t

	c.EnergyMax = 0
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c
	c.nightsoulState = nightsoul.New(c.Core, c.CharWrapper)
	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.burstInit()
	c.a1()
	c.c1Init()
	c.c2Init()
	c.a4Init()

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		switch a {
		case action.ActionCharge:
			return 0
		case action.ActionDash:
			return 0
		}
	}

	if a == action.ActionCharge {
		return 50
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	switch a {
	case action.ActionBurst:
		if c.fightingSpirit < 100 {
			return false, action.InsufficientEnergy
		}
		return c.Character.ActionReady(a, p)
	case action.ActionSkill:
		if p["recast"] != 0 && !c.StatusIsActive(skillRecastCDKey) {
			return true, action.NoFailure
		}
		return c.Character.ActionReady(a, p)
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		c.DeleteStatus(burstKey)
		if c.armamentState == bike && c.nightsoulState.HasBlessing() {
			c.exitBike()
		}
		return false
	}, "mavuika-exit")
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	case "fightingspirit":
		return c.fightingSpirit, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		switch k {
		case model.AnimationXingqiuN0StartDelay:
			return 0
		case model.AnimationYelanN0StartDelay:
			return 0
		default:
			return c.Character.AnimationStartDelay(k)
		}
	}

	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 22
	case model.AnimationYelanN0StartDelay:
		return 22
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
