package barbara

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("barbara", NewChar)
}

type char struct {
	*character.Tmpl
	stacks int
	// burstBuffExpiry   int
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
	c.Weapon.Class = core.WeaponClassCatalyst
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 4

	c.a1()
	c.onSkillStackCount() //doesnt' do anything yet
	return &c, nil
}

func (c *char) a1() {
	c.Core.AddStamMod(func(a core.ActionType) (float64, bool) { // @srl does this activate for the active char?
		if c.Core.Status.Duration("barbara-field") >= 0 {
			return -0.12, false
		}
		return 0, false
	})
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	/*
		Returns character stamina consumption for specified action.
	*/
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 50
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}
