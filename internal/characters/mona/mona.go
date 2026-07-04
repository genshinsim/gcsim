package mona

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
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
	c6AtkMod         character.AttackMod
	c6NearbyOmen     bool
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

	c.c1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	c.c6HexInit()

	return nil
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	// TODO: you can do the CA after the N4 resets into idle
	if a == action.ActionCharge {
		if c.Core.Player.LastAction.Type == action.ActionAttack && c.NormalCounter == 0 {
			return player.ErrInvalidChargeAction
		}
	}
	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "astral-glow":
		return c.astralGlowStacks, nil
	case "c6-stacks":
		if !c.StatusIsActive(c6Key) {
			return 0, nil
		}
		return c.c6Stacks, nil
	default:
		return c.Character.Condition(fields)
	}
}
