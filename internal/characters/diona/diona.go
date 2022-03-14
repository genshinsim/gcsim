package diona

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Diona, NewChar)
}

type char struct {
	*character.Tmpl
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = coretype.Cryo

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
		if c.Core.Player.GetShield(core.ShieldDionaSkill) != nil {
			return -0.1, false
		}
		return 0, false
	}, "diona")
}

func (c *char) c2() {
	c.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "diona-c2",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			val[core.DmgP] = .15
			return val, atk.Info.AttackTag == core.AttackTagElementalArt
		},
	})
}
func (c *char) c6() {
	c.Core.Health.AddIncHealBonus(func(healedCharIndex int) float64 {
		if c.Core.StatusDuration("dionaburst") == 0 {
			return 0
		}
		char := c.Core.Chars[c.Core.ActiveChar]
		if healedCharIndex != char.Index() {
			return 0
		}
		if char.HP()/char.MaxHP() <= 0.5 {
			c.coretype.Log.NewEvent("diona c6 activated", coretype.LogCharacterEvent, c.Index)
			return 0.3
		}
		return 0
	})
}
