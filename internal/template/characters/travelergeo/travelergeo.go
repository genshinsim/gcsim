package travelergeo

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl
}

func init() {
	core.RegisterCharFunc(core.TravelerGeo, NewChar)
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
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

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
			Amount: func() ([]float64, bool) {
				if c.Core.Constructs.CountByType(core.GeoConstructTravellerBurst) == 0 {
					return nil, false
				}
				return val, true
			},
		})
	}
}
