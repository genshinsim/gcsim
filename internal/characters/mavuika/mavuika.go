package mavuika

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
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
	chargeCancel       bool
	dashFrames         []int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)

	t.CharWrapper = w
	c.Character = t

	c.EnergyMax = 0
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	fs, ok := p.Params["start_energy"]
	if !ok {
		fs = maxFightingSpirit
	}
	fs = max(min(fs, maxFightingSpirit), 0)

	c.fightingSpirit = float64(fs)

	w.Character = &c
	c.nightsoulState = nightsoul.New(c.Core, c.CharWrapper)
	c.dashFrames = frames.InitAbilSlice(24) // Dash -> Dash
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
		if !c.Core.Flags.IgnoreBurstEnergy && c.fightingSpirit < 100 {
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
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...any) {
		c.DeleteStatus(burstKey)
		if c.armamentState == bike && c.nightsoulState.HasBlessing() {
			c.exitBike()
		}
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

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if c.armamentState == bike && c.nightsoulState.HasBlessing() {
		switch k {
		case info.AnimationXingqiuN0StartDelay:
			return 0
		case info.AnimationYelanN0StartDelay:
			return 0
		default:
			return c.Character.AnimationStartDelay(k)
		}
	}

	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 22
	case info.AnimationYelanN0StartDelay:
		return 22
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) NextQueueItemIsValid(targetChar keys.Char, a action.Action, p map[string]int) error {
	if c.chargeCancel {
		if targetChar != c.Base.Key {
			return errors.New("cannot swap, Mavuika must perform a Biked charge attack after a dash cancel")
		}
		if !c.nightsoulState.HasBlessing() {
			return errors.New("nightsoul blessing expired before Mavuika could perform charge cancel")
		}
		if c.Core.Player.CurrentState() != action.DashState {
			return errors.New("cannot allow Mavuika to go into idle during charge-cancelled dash")
		}
		if a != action.ActionCharge {
			return fmt.Errorf("cannot perform action %s, Mavuika must perform a Biked charge attack after a dash cancel", a)
		}
	}
	return c.Character.NextQueueItemIsValid(targetChar, a, p)
}
