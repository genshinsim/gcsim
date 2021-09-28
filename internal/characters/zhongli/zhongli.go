package zhongli

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

type char struct {
	*character.Tmpl
	maxSteele   int
	steeleCount int
	energyICD   int
}

func init() {
	core.RegisterCharFunc("zhongli", NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 90
	c.EnergyMax = 90
	c.Weapon.Class = core.WeaponClassSpear
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	c.maxSteele = 1
	if c.Base.Cons >= 1 {
		c.maxSteele = 2
	}

	c.a2()

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) a2() {
	c.Core.Shields.AddBonus(func() float64 {
		if c.Tags["shielded"] == 0 {
			return 0
		}
		return float64(c.Tags["a2"]) * 0.05
	})
}
