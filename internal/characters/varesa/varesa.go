package varesa

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
)

func init() {
	core.RegisterCharFunc(keys.Varesa, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState *nightsoul.State

	particleGenerated bool
	freeSkill         bool
	exitNS            bool
	a4Stacks          *stacks.MultipleRefreshNoRemove
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.BurstCon = 3
	c.NormalCon = 5
	c.NormalHitNum = 3

	c.SetNumCharges(action.ActionSkill, 2)

	c.nightsoulState = nightsoul.New(s, w)
	c.nightsoulState.MaxPoints = 40

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.c4()
	c.c6()

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge && c.StatusIsActive(skillStatus) {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	if a == action.ActionSkill && c.freeSkill {
		return true, action.NoFailure
	}
	if a == action.ActionBurst && c.StatusIsActive(apexState) {
		if !c.Core.Flags.IgnoreBurstEnergy && c.Energy < kablamCost {
			return false, action.InsufficientEnergy
		}
		if c.AvailableCDCharge[a] <= 0 {
			return false, action.BurstCD
		}
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	case skillStatus:
		return c.StatusIsActive(skillStatus), nil
	case apexState:
		return c.StatusIsActive(apexState), nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) generatePlungeNightsoul() {
	c.nightsoulState.GeneratePoints(25)
	if !c.nightsoulState.HasBlessing() && c.nightsoulState.Points() == c.nightsoulState.MaxPoints {
		c.nightsoulState.EnterBlessing(c.nightsoulState.MaxPoints)
		c.freeSkill = true // TODO: duration?
	}
}

func (c *char) clearNightsoul(next action.AnimationState) {
	// ignore volcanic kablam
	if next == action.BurstState {
		return
	}
	if c.exitNS {
		c.nightsoulState.ExitBlessing()
	}
}
