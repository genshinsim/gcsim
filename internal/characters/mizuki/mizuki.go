package mizuki

import (
	"fmt"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
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
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.registerSkillCallbacks()
	c.a1()
	c.a4()
	return nil
}

func (c *char) NextQueueItemIsValid(targetChar keys.Char, a action.Action, p map[string]int) error {
	if !c.StatusIsActive(dreamDrifterStateKey) {
		return nil
	}
	if a != action.ActionDash && a != action.ActionSwap && a != action.ActionBurst && a != action.ActionSkill {
		return fmt.Errorf("%v: Tried to execute %v when not dreamdrifter state", c.Base.Key, a)
	}
	return nil
}
