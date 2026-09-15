package iansan

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

type char struct {
	*tmpl.Character
	nightsoulState *nightsoul.State

	nightsoulSrc      int
	particleGenerated bool
	burstSrc          int
	burstBuff         []float64
	burstRestoreNS    int
	pointsOverflow    float64

	a1Buff     []float64
	a1Increase bool

	c1Points    float64
	c2Buff      []float64
	c4Generated bool
	c4Stacks    int
	c6Buff      []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.SkillCon = 3
	c.BurstCon = 5
	c.NormalHitNum = normalHitNum

	c.nightsoulState = nightsoul.New(s, w)
	c.nightsoulState.MaxPoints = 54

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstInit()
	c.a1Init()
	c.a4Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()

	c.Core.Events.Subscribe(event.OnActionExec, c.burstMovementRestore, burstBuffStatus)
	// TODO: subscribe to a player moved event
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge && c.StatusIsActive(fastSkill) {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 10
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	if a == action.ActionCharge && c.StatusIsActive(fastSkill) {
		// allow skill, charge
		return nil
	}
	return c.Character.NextQueueItemIsValid(k, a, p)
}
