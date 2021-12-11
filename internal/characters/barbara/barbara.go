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
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 3

	c.onExitField()

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return &c, nil
}

// Hook that clears barbara burst status and seals when she leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		c.Tags["seal"] = 0
		c.Core.Status.DeleteStatus("barbaraburst")
		return false
	}, "barbara-exit")
}

// Hook for C2:
// Increases Yan Fei's Charged Attack CRIT Rate by 20% against enemies below 50% HP.
func (c *char) c2() {
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "barbara-c2",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([core.EndStatType]float64, bool) {
			var m [core.EndStatType]float64
			if atk.Info.AttackTag != core.AttackTagExtra {
				return m, false
			}
			if t.HP()/t.MaxHP() >= .5 {
				return m, false
			}
			m[core.CR] = 0.20
			return m, true
		},
	})
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}
