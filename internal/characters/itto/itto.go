package itto

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Itto, NewChar)
}

type char struct {
	*character.Tmpl
	dasshuUsed  bool
	dasshuCount int
	sCACount    int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Geo
	c.Energy = 70
	c.EnergyMax = 70
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4
	c.dasshuUsed = false
	c.dasshuCount = 0
	c.Tags["strStack"] = 0
	c.sCACount = 0

	c.onExitField()

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Tags["strStack"] > 0 {
			return 0
		}
		return 20
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}

}

// Itto Geo infusion can't be overridden, so it must be a snapshot modification rather than a weapon infuse
func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Core.Status.Duration("ittoq") > 0 {
		//infusion to attacks only
		switch ai.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagPlunge:
		case core.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = core.Geo
	}
	return ds
}
