package sara

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Sara, NewChar)
}

type char struct {
	*character.Tmpl
	a4LastProc int
	c1LastProc int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	if c.Base.Cons == 6 {
		c.c6()
	}
}

func (c *char) c6() {
	val := make([]float64, core.EndStatType)
	val[core.CD] = 0.6
	for _, char := range c.Core.Chars {
		this := char
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "sara-c6",
			Expiry: -1,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				//check if tags active
				if this.Tag("sarabuff") < c.Core.F {
					return nil, false
				}
				if atk.Info.Element != core.Electro {
					return nil, false
				}
				return val, true
			},
		})
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Key.String())
		return 0
	}
}
