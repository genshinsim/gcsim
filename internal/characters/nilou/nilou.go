package nilou

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Nilou, NewChar)
}

type char struct {
	*tmpl.Character
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExitField()

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.ActionFailure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(pirouetteStatus) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}
