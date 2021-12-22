package ningguang

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Ningguang, NewChar)
}

type char struct {
	*character.Tmpl
	c2reset     int
	lastScreen  int
	particleICD int
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
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 1
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue
	// Initialize at some very low value so reset happens correctly at start of sim
	c.c2reset = -9999

	c.a4()

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Tags["jade"] > 0 {
			return 0
		}
		return 50
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}

}

func (c *char) a4() {
	//activate a4 if screen is down and character uses dash
	c.Core.Events.Subscribe(core.OnDash, func(args ...interface{}) bool {
		if c.Core.Constructs.CountByType(core.GeoConstructNingSkill) > 0 {
			val := make([]float64, core.EndStatType)
			val[core.GeoP] = 0.12
			char := c.Core.Chars[c.Core.ActiveChar]
			char.AddMod(core.CharStatMod{
				Key: "ning-screen",
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return val, true
				},
				Expiry: c.Core.F + 600,
			})
		}
		return false
	}, "ningguang-a4")
}
