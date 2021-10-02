package xiao

import (
	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/character"
)

func init() {
	core.RegisterCharFunc("xiao", NewChar)
}

// Xiao specific character implementation
type char struct {
	*character.Tmpl
	eCharge      int
	eChargeMax   int
	eNextRecover int
	eTickSrc     int
	qStarted     int
	a4Expiry     int
}

// Initializes character
// TODO: C4 is not implemented - don't really care about def
func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 70
	c.EnergyMax = 70
	c.Weapon.Class = core.WeaponClassSpear
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 5

	c.eChargeMax = 2
	if c.Base.Cons >= 1 {
		c.eChargeMax = 3
	}

	if c.Base.Cons >= 2 {
		c.c2()
	}

	c.eCharge = c.eChargeMax
	c.onExitField()

	return &c, nil
}

// Implements Xiao C2:
// When in the party and not on the field, Xiao's Energy Recharge is increased by 25%
func (c *char) c2() {
	stat_mod := make([]float64, core.EndStatType)
	stat_mod[core.ER] = 0.25
	c.AddMod(core.CharStatMod{
		Key: "xiao-c2",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if c.Core.ActiveChar != c.Index {
				return stat_mod, true
			}
			return nil, false
		},
	})
}

// Hook to end Xiao's burst prematurely if he leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		c.Core.Status.DeleteStatus("xiaoburst")
		return false
	}, "xiao-exit")
}

// Stamina usage values
func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}
