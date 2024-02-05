package xianyun

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	c.SetNumCharges(action.ActionSkill, 2)
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.a4Max = 18000
	c.a4Ratio = 4
}

func (c *char) c2buff() {
	if c.Base.Cons < 2 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.20
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("xianyun-c2", 15*60),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

// TODO: C4 Xianyun
// TODO: C6 Xianyun
func (c *char) c6() {}
