package keqing

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Keqing, NewChar)
}

type char struct {
	*tmpl.Character
	a4buff []float64
	c4buff []float64
	c6buff []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4buff = make([]float64, attributes.EndStatType)
	c.a4buff[attributes.CR] = 0.15
	c.a4buff[attributes.ER] = 0.15

	if c.Base.Cons >= 4 {
		c.c4buff = make([]float64, attributes.EndStatType)
		c.c4buff[attributes.ATKP] = 0.25
		c.c4()
	}
	if c.Base.Cons >= 6 {
		c.c6buff = make([]float64, attributes.EndStatType)
		c.c6buff[attributes.ElectroP] = 0.06
	}
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.ActionFailure) {
	// check if stiletto is on-field
	if a == action.ActionSkill && c.Core.Status.Duration(stilettoKey) > 0 {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		return 25
	}
	return c.Character.ActionStam(a, p)
}
