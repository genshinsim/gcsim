package shenhe

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Shenhe, NewChar)
}

type char struct {
	*character.Tmpl
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 5
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = core.ZoneLiyue

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	if c.Base.Cons >= 6 {
		c.c6()
	}
	c.a2()
	c.a4()
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		case 0:
			f = 12
		case 1:
			f = 38 - 12
		case 2:
			f = 72 - 38
		case 3:
			f = 141 - 72
		case 4:
			f = 167 - 141
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionSkill:
		return 26, 26
	case core.ActionBurst:
		return 99, 99
	case core.ActionCharge:
		return 78, 78
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Key.String(), a)
		return 0, 0
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}

}

// inspired from barbara c2
// technically always assumes you are inside shenhe burst
func (c *char) a2() {
	val := make([]float64, core.EndStatType)
	val[core.CryoP] = 0.15
	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue
		}
		char.AddMod(core.CharStatMod{
			Key:    "shenhe-a2",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if c.Core.Status.Duration("shenheburst") >= 0 {
					return val, true
				} else {
					return nil, false
				}
			},
		})
	}
}

func (c *char) a4() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15
	for i, char := range c.Core.Chars {
		if i == c.Index {
			continue
		}
		char.AddMod(core.CharStatMod{
			Key: "shenhe-a2",
			Expiry: func() int {
				if c.Core.Status.Duration("shenheskillpress") >= 0 {
					return c.Core.F + c.Core.Status.Duration("shenheskillpress")
				} else if c.Core.Status.Duration("shenheskillhold") >= 0 {
					return c.Core.F + c.Core.Status.Duration("shenheskillhold")
				} else {
					return 0
				}
			}(),
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if c.Core.Status.Duration("shenheskillpress") >= 0 {
					if a != core.AttackTagElementalBurst && a != core.AttackTagElementalArt && a != core.AttackTagElementalArtHold {
						return nil, false
					}
					return val, true
				} else if c.Core.Status.Duration("shenheskillhold") >= 0 {
					if a != core.AttackTagNormal && a != core.AttackTagExtra && a != core.AttackTagPlunge {
						return nil, false
					}
					return val, true
				} else {
					return nil, false
				}
			},
		})
	}
}

func (c *char) c6() {
	m := make([]float64, core.EndStatType)
	m[core.PyroP] = 0.15

	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Key:    "xl-c6",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, c.Core.Status.Duration("xlc6") > 0
			},
		})
	}
}
