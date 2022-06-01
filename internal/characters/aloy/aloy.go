package aloy

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl
	coilICDExpiry int
	lastFieldExit int
}

func init() {
	core.RegisterCharFunc(core.Aloy, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 4

	c.coilICDExpiry = 0
	c.lastFieldExit = 0

	c.Tags["coil_stacks"] = 0

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.coilMod()
	c.onExitField()
}

// Add coil mod at the beginning of the sim
// Can't be made dynamic easily as coils last until 30s after when Aloy swaps off field
func (c *char) coilMod() {
	val := make([]float64, core.EndStatType)
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "aloy-coil-stacks",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag == core.AttackTagNormal && c.Tags["coil_stacks"] > 0 {
				val[core.DmgP] = skillCoilNABonus[c.Tags["coil_stacks"]-1][c.TalentLvlSkill()]
				return val, true
			}
			return nil, false
		},
	})
}

// Exit Field Hook to start timer to clear coil stacks
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev != c.Index {
			return false
		}
		c.lastFieldExit = c.Core.F

		c.AddTask(func() {
			if c.lastFieldExit != (c.Core.F - 30*60) {
				return
			}
			c.Tags["coil_stacks"] = 0
		}, "aloy-on-field-exit", 30*60)

		return false
	}, "aloy-exit")
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}
