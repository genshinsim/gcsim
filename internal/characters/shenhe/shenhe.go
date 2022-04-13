package shenhe

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Shenhe, NewChar)
}

type char struct {
	*character.Tmpl
	quillcount []int
	c4count    int
	c4expiry   int
}

const (
	quillKey = "shenhequill"
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 5
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = core.ZoneLiyue
	c.Base.Element = core.Cryo

	c.c4count = 0
	c.c4expiry = 0

	if c.Base.Cons >= 1 {
		c.SetNumCharges(core.ActionSkill, 2)
	}

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.quillcount = make([]int, len(c.Core.Chars))

	c.a1()
	c.quillDamageMod()

	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}
}

// Helper function to update tags that can be used in configs
// Should be run whenever c.quillcount is updated
func (c *char) updateBuffTags() {
	for _, char := range c.Core.Chars {
		c.Tags["quills_"+char.Name()] = c.quillcount[char.CharIndex()]
		c.Tags[fmt.Sprintf("quills_%v", char.CharIndex())] = c.quillcount[char.CharIndex()]
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}

}

// inspired from barbara c2
// TODO: technically always assumes you are inside shenhe burst
func (c *char) a1() {
	val := make([]float64, core.EndStatType)
	val[core.CryoP] = 0.15
	for _, char := range c.Core.Chars {
		// if i == c.Index {
		// 	continue
		// }
		char.AddMod(core.CharStatMod{
			Key:    "shenhe-a1",
			Expiry: -1,
			Amount: func() ([]float64, bool) {
				if c.Core.Status.Duration("shenheburst") > 0 {
					return val, true
				} else {
					return nil, false
				}
			},
		})
	}
}

func (c *char) c6() {
	val := make([]float64, core.EndStatType)
	val[core.CD] = 0.15
	for _, char := range c.Core.Chars {
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "shenhe-c2",
			Expiry: -1,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				//check if tags active
				if c.Core.Status.Duration("shenheburst") <= 0 {
					return nil, false
				}
				if atk.Info.Element != core.Cryo {
					return nil, false
				}
				return val, true
			},
		})
	}
}

func (c *char) c4() {
	c.AddPreDamageMod(core.PreDamageMod{
		Expiry: -1,
		Key:    "shenhe-c4",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)

			if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
				return nil, false
			}
			if c.Core.F >= c.c4expiry {
				return nil, false
			}
			val[core.DmgP] += 0.05 * float64(c.c4count)
			c.c4count = 0
			c.c4expiry = 0
			return val, true
		},
	})
}
