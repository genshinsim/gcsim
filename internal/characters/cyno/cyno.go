package cyno

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Cyno, NewChar)
}

type char struct {
	*tmpl.Character
	burstExtension int
	burstSrc       int
	c2Stacks       int
	c4Counter      int
	c6Stacks       int
	a1Extended     bool
	normalBCounter int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.a1Extension()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return nil
}

func (c *char) AdvanceNormalIndex() {
	if c.StatusIsActive(BurstKey) {
		c.normalBCounter++
		if c.normalBCounter == burstHitNum {
			c.normalBCounter = 0
		}
		return
	}
	c.NormalCounter++
	if c.NormalCounter == c.NormalHitNum {
		c.NormalCounter = 0
	}
}

func (c *char) ResetNormalCounter() {
	c.normalBCounter = 0
	c.NormalCounter = 0
}

func (c *char) NextNormalCounter() int {
	if c.StatusIsActive(BurstKey) {
		return c.normalBCounter + 1
	}
	return c.NormalCounter + 1
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.ActionFailure) {
	if a != action.ActionSkill {
		return c.Character.ActionReady(a, p)
	}
	if c.StatusIsActive(BurstKey) {
		if c.AvailableCDCharge[action.ActionLowPlunge] <= 0 {
			return false, action.SkillCD
		}
		return true, action.NoFailure
	}
	if c.AvailableCDCharge[action.ActionSkill] <= 0 {
		return false, action.SkillCD
	}
	return true, action.NoFailure
}

func (c *char) ReduceActionCooldown(a action.Action, v int) {
	c.Character.ReduceActionCooldown(a, v)
	if a == action.ActionSkill {
		c.Character.ReduceActionCooldown(action.ActionLowPlunge, v)
	}
}

func (c *char) ResetActionCooldown(a action.Action) {
	c.Character.ResetActionCooldown(a)
	if a == action.ActionSkill {
		c.Character.ResetActionCooldown(action.ActionLowPlunge)
	}
}
