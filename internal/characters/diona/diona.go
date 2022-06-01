package diona

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Diona, NewChar)
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
	c.Base.Element = core.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a1()

	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons == 6 {
		c.c6()
	}
}

func (c *char) a1() {
	c.Core.AddStamMod(func(a core.ActionType) (float64, bool) {
		if c.Core.Shields.Get(core.ShieldDionaSkill) != nil {
			return -0.1, false
		}
		return 0, false
	}, "diona")
}

func (c *char) c2() {
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "diona-c2",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			val[core.DmgP] = .15
			return val, atk.Info.AttackTag == core.AttackTagElementalArt
		},
	})
}
func (c *char) c6() {
	c.Core.Health.AddIncHealBonus(func(healedCharIndex int) float64 {
		//check tag for if bonus is active
		char := c.Core.Chars[healedCharIndex]
		if c.Core.F > c.Tags["c6bonus-"+char.Key().String()] {
			return 0
		}
		c.Core.Log.NewEvent("diona c6 incomming heal bonus activated", core.LogCharacterEvent, c.Index)
		return 0.3
	})
}
