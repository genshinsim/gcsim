package xianyun

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Xianyun, NewChar)
}

type char struct {
	*tmpl.Character
	skillCounter int
	skillSrc     int
	a1Count      int
	a1Buffer     []int
	a4Max        float64
	a4Ratio      float64

	// leapFrames  []int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = 4
	c.SkillCon = 5
	c.BurstCon = 3
	c.skillSrc = -1

	// The number of characters in s.Player.Chars() is wrong at this point, so we need to do it in Init()
	// c.a1Buffer = make([]int, len(c.Core.Player.Chars()))

	c.a4Max = 9000
	c.a4Ratio = 2.0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1Buffer = make([]int, len(c.Core.Player.Chars()))
	c.a1()
	c.a4()
	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(skillStateKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}
