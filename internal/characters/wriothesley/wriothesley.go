package wriothesley

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	revelation           bool
	savedNormalCounter   int
	rebukeHeal           float64
	a4Stack              int
	c1N5Proc             bool
	c1SkillExtensionProc bool
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.NormalCon = 3
	c.BurstCon = 5
	c.HasArkhe = true

	revelation, ok := p.Params["revelation"]
	if !ok {
		revelation = 1
	}
	c.revelation = revelation > 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExit()
	c.stellarInit()
	c.a4()
	c.c4()

	return nil
}

func (c *char) graciousRebukeReady() bool {
	if c.Base.Ascension < 1 {
		return false
	}
	if c.Base.Cons >= 1 {
		return c.c1Ready()
	}
	return c.a1Ready()
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge && (c.graciousRebukeReady() || c.a1StellarCAReady()) {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "gracious-rebuke":
		return c.graciousRebukeReady(), nil
	case "gracious-rebuke-unveiling":
		return c.a1StellarCAReady(), nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	// cannot use charge without attack beforehand unlike most of the other catalyst users
	if a == action.ActionCharge && c.Core.Player.LastAction.Type != action.ActionAttack {
		return player.ErrInvalidChargeAction
	}
	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 12
	case info.AnimationYelanN0StartDelay:
		return 4
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
