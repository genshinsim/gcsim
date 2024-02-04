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
	eCounter    int
	skillHeight float64
	skillRadius float64
	a1Count     int
	// leapFrames  []int
}

const eWindowKey = "xianyun-e-window"

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = 4
	c.SkillCon = 5
	c.BurstCon = 3

	c.skillHeight = 0
	c.skillRadius = 0
	c.eCounter = 0

	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// c.a4()
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(eWindowKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}
