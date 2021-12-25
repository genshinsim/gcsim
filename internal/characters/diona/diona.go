package diona

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Diona, NewChar)
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
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.a2()

	if c.Base.Cons == 6 {
		c.c6()
	}

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return &c, nil
}

func (c *char) a2() {
	c.Core.AddStamMod(func(a core.ActionType) (float64, bool) {
		if c.Core.Shields.Get(core.ShieldDionaSkill) != nil {
			return -0.1, false
		}
		return 0, false
	})
}

func (c *char) c2() {
	c.AddMod(core.CharStatMod{
		Key:    "diona-c2",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			val[core.DmgP] = .15
			return val, a == core.AttackTagElementalArt
		},
	})
}
func (c *char) c6() {
	c.Core.Health.AddIncHealBonus(func(healedCharIndex int) float64 {
		if c.Core.Status.Duration("dionaburst") == 0 {
			return 0
		}
		char := c.Core.Chars[c.Core.ActiveChar]
		if healedCharIndex != char.CharIndex() {
			return 0
		}
		if char.HP()/char.MaxHP() <= 0.5 {
			c.Core.Log.Debugw("diona c6 activated", "frame", c.Core.F, "event", core.LogCharacterEvent)
			return 0.3
		}
		return 0
	})
}
