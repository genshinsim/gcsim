package thoma

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Thoma, NewChar)
}

type char struct {
	*character.Tmpl
	MaxShield float64
	a1Stack   int
	a1icd     int
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
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 4
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = core.ZoneInazuma
	c.Base.Element = core.Pyro
	c.MaxShield = (shieldppmax[c.TalentLvlSkill()]*c.MaxHP() + shieldflatmax[c.TalentLvlSkill()])
	c.a1Stack = 0
	c.a1icd = 0
	c.a1()

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

func (c *char) a1() {
	c.Core.Shields.AddBonus(func() float64 {
		if c.Tags["shielded"] == 0 {
			return 0
		}
		if c.Core.Status.Duration("thoma-a1") <= 0 {
			return 0
		}
		return float64(c.a1Stack) * 0.05
	})

	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		c.a1Stack = 0
		return false
	}, "thoma-a1-swap")
}
