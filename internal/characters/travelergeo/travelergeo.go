package travelergeo

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type char struct {
	*character.Tmpl
}

func init() {
	core.RegisterCharFunc(keys.TravelerGeo, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Geo
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)

	if c.Base.Cons > 0 {
		c.c1()
	}

}

//Party members within the radius of Wake of Earth have their CRIT Rate increased by 10%
//and have increased resistance against interruption.
func (c *char) c1() {
	val := make([]float64, core.EndStatType)
	val[core.CR] = .1
	for _, char := range c.Core.Chars {
		char.AddMod(core.CharStatMod{
			Key:    "geo-traveler-c1",
			Expiry: -1,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				if c.Core.Constructs.CountByType(core.GeoConstructTravellerBurst) == 0 {
					return nil, false
				}
				return val, true
			},
		})
	}
}
