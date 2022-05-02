package itto

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
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

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 70
	}
	c.Energy = float64(e)
	c.EnergyMax = 70
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4
	c.SkillCon = 3
	c.BurstCon = 5

	c.dasshuUsed = false
	c.dasshuCount = 0
	c.Tags["strStack"] = 0
	c.sCACount = 0

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.onExitField()

	if c.Base.Cons == 6 {
		c.c6()
	}
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
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
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

func (c *char) c6() {
	val := make([]float64, core.EndStatType)
	val[core.CD] = 0.7
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "itto-c6",
		Expiry: -1,
		Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
			if a.Info.AttackTag != core.AttackTagExtra {
				return nil, false
			}
			return val, true
		},
	})
}
