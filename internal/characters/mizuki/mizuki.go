package mizuki

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Mizuki, NewChar)
}

type char struct {
	*tmpl.Character
	particleGenerationsRemaining    int
	dreamDrifterExtensionsRemaining int
	cloudAttack                     combat.AttackInfo
	cloudSnap                       combat.Snapshot
	cloudSrc                        int
	a4Buff                          []float64
	c1EM                            float64
	c2Buff                          []float64
	c4EnergyGenerationsRemaining    int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillInit()
	c.a1()
	c.a4()
	c.c1()
	c.c2()
	c.c6()
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	if !c.StatusIsActive(dreamDrifterStateKey) {
		return c.Character.ActionReady(a, p)
	}

	// Allow to press E while in dreamDrifter that cancels the effect
	if a == action.ActionSkill {
		// cancel
		return true, action.NoFailure
	}

	// restrict actions that can be performed while in dreamDrifter stat to swap, dash, burst
	if a != action.ActionDash && a != action.ActionSwap && a != action.ActionBurst {
		return false, action.NoFailure
	}

	return c.Character.ActionReady(a, p)
}
