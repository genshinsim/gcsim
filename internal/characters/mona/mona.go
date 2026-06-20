package mona

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	bubbleKey = "mona-bubble"
	omenKey   = "omen-debuff"
)

func init() {
	core.RegisterCharFunc(keys.Mona, NewChar)
}

type char struct {
	*tmpl.Character
	a4Stats          []float64
	c6Src            int
	c6Stacks         int
	astralGlowStacks int
	astralGlowSrc    int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	hex, ok := p.Params["hexerei"]
	if !ok {
		// default hexerei is enabled
		hex = 1
	}
	c.IsHexerei = (hex != 0)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstHook()
	c.burstDamageBonus()
	c.a4()
	c.hexInit()
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons >= 6 {
		c.c6Init()
		c.c6ChargeAttackInit()
	}
	return nil
}

func (c *char) omenIsNearby() bool {
	// TODO: check range of this in DM
	for _, e := range c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 15), nil) {
		if e.StatusIsActive(omenKey) || e.StatusIsActive(bubbleKey) {
			return true
		}
	}
	return false
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	// TODO: you can do the CA after the N4 resets into idle
	if c.Core.Player.LastAction.Type == action.ActionAttack && c.NormalCounter == 0 {
		return player.ErrInvalidChargeAction
	}
	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "astral-glow":
		return c.astralGlowStacks, nil
	case "c6-stacks":
		return c.c6Stacks, nil
	default:
		return c.Character.Condition(fields)
	}
}
