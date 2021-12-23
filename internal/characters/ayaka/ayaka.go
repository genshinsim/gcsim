package ayaka

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type char struct {
	*character.Tmpl
}

func init() {
	core.RegisterCharFunc(keys.Ayaka, NewChar)
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
	c.Weapon.Class = core.WeaponClassSword
	c.CharZone = core.ZoneInazuma
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		f, ok := p["f"]
		if !ok {
			return 10 //tap = 36 frames, so under 1 second
		}
		//for every 1 second passed, consume extra 15
		extra := f / 60
		return float64(10 + 15*extra)
	case core.ActionCharge:
		return 20
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Key.String())
		return 0
	}
}
