package amber

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Amber, NewChar)
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
	c.Energy = 40
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
	c.AddMod(core.CharStatMod{
		Key:    "amber-a2",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			v := make([]float64, core.EndStatType)
			v[core.CR] = .1
			return v, a == core.AttackTagElementalBurst
		},
	})
}
