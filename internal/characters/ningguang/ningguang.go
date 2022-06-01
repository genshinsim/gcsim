package ningguang

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Ningguang, NewChar)
}

type char struct {
	*character.Tmpl
	c2reset       int
	lastScreen    int
	particleICD   int
	skillSnapshot core.Snapshot
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
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 1
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = core.ZoneLiyue

	// Initialize at some very low value so these happen correctly at start of sim
	c.c2reset = -9999
	c.particleICD = -9999

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a4()
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
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}

}

func (c *char) a4() {
	//activate a4 if screen is down and character uses dash
	c.Core.Events.Subscribe(core.PostDash, func(args ...interface{}) bool {
		if c.Core.Constructs.CountByType(core.GeoConstructNingSkill) > 0 {
			val := make([]float64, core.EndStatType)
			val[core.GeoP] = 0.12
			char := c.Core.Chars[c.Core.ActiveChar]
			char.AddMod(core.CharStatMod{
				Key: "ning-screen",
				Amount: func() ([]float64, bool) {
					return val, true
				},
				Expiry: c.Core.F + 600,
			})
		}
		return false
	}, "ningguang-a4")
}
