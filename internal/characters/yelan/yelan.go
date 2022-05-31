package yelan

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Yelan, NewChar)
}

type char struct {
	*character.Tmpl
	c2icd        int
	burstDiceICD int
	burstTickSrc int
	c6count      int
	c4count      int //keep track of number of enemies tagged
}

const eCD = 600

const (
	breakthroughStatus = "yelan_breakthrough"
	c6Status           = "yelan_c6"
	burstStatus        = "yelanburst"
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 70
	}
	c.Energy = float64(e)
	c.EnergyMax = 70
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 4
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	barb, ok := p.Params["barb"]
	if ok && barb > 0 {
		c.AddTag(breakthroughStatus, 1)
	}

	c.burstStateHook()
	c.SetNumCharges(core.ActionSkill, 1)
	if c.Base.Cons >= 1 {
		c.SetNumCharges(core.ActionSkill, 2)
	}

	return &c, nil
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

func (c *char) Init() {
	c.Tmpl.Init()
	c.c2icd = 0
	c.c6count = 0
	c.a1()

}

func (c *char) a1() {
	m := make(map[core.EleType]bool)
	for _, char := range c.Core.Chars {
		m[char.Ele()] = true
	}
	val := make([]float64, core.EndStatType)
	l := float64(len(m))
	if l > 3 {
		val[core.HPP] = 0.3
	} else {
		val[core.HPP] = l * 0.06
	}
	c.AddMod(core.CharStatMod{
		Key:    "yelan-a1",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})
}
