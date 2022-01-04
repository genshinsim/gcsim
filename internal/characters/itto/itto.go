package itto

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Itto, NewChar)
}

type char struct {
	*character.Tmpl
	a2stacks    int
	skillStacks int
	ushiOnField bool
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
	c.a2stacks = 0
	c.skillStacks = 0
	c.ushiOnField = false

	c.a4()
	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.skillStacks > 0 {
			return 0
		}
		return 20
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}

}

func (c *char) a4() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = c.Stats[core.DEF] * 0.35
	c.AddPreDamageMod(core.PreDamageMod{
		Key: "itto-a4",
		Amount: func(ae *core.AttackEvent, t core.Target) ([]float64, bool) {
			if c.skillStacks > 0 && ae.Info.AttackTag == core.AttackTagExtra {
				return val, true
			}
			return nil, false
		},
		Expiry: -1,
	})
}
