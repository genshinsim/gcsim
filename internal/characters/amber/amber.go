package amber

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Amber, NewChar)
}

type char struct {
	*character.Tmpl
	bunnies      []bunny
	eCharge      int
	eChargeMax   int
	eNextRecover int
	eTickSrc     int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Pyro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.eChargeMax = 1
	if c.Base.Cons >= 4 {
		c.eChargeMax = 2
	}
	c.eCharge = c.eChargeMax

	if c.Base.Cons >= 2 {
		c.overloadExplode()
	}
	c.a2()
	c.bunnies = make([]bunny, 0, 2)

	return &c, nil
}

func (c *char) a2() {
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "amber-a2",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			v := make([]float64, core.EndStatType)
			v[core.CR] = .1
			return v, atk.Info.AttackTag == core.AttackTagElementalBurst
		},
	})
}
