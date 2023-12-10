package freminet

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Freminet, NewChar)
}

type char struct {
	*tmpl.Character
	skillStacks int
	persID      int
	c4Stacks    int
	c6Stacks    int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.NormalCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()

	c.c1()
	c.c4c6()

	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.ActionFailure) {
	if a == action.ActionSkill && c.StatusIsActive(persTimeKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}
