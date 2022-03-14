package rosaria

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

type char struct {
	*character.Tmpl
}

func init() {
	core.RegisterCharFunc(core.Rosaria, NewChar)
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = coretype.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSpear
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 5

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

// Adds event checker for C1: Unholy Revelation
// When Rosaria deals a CRIT Hit, her ATK Speed increase by 10% and her Normal Attack DMG increases by 10% for 4s (can trigger vs shielded enemies)
// TODO: Description is unclear whether attack speed affects NA + CA - assume that it only affects NA for now
func (c *char) c1() {
	// Add hook that monitors for crit hits. Mirrors existing favonius code
	// No log value saved as stat mod already shows up in debug view
	c.Core.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		crit := args[3].(bool)
		if !crit {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		val := make([]float64, core.EndStatType)
		val[core.AtkSpd] = 0.1
		val[core.DmgP] = 0.1
		c.AddPreDamageMod(coretype.PreDamageMod{
			Key:    "rosaria-c1",
			Expiry: c.Core.Frame + 240,
			Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				if atk.Info.AttackTag != coretype.AttackTagNormal {
					return nil, false
				}
				return val, true
			},
		})
		return false
	}, "rosaria-c1")
}

// Adds event checker for C4 Painful Grace
// Ravaging Confession's CRIT Hits regenerate 5 Energy for Rosaria. Can only be triggered once each time Ravaging Confession is cast.
// Only applies when a crit hit is resolved, so can't be handled within skill code directly
// TODO: Since this only is needed for her E, can change this so it spawns a subscription in her E code
// Then it can return true, which kills the callback
// However, would also need a timeout function as well since her E can not crit
// Requires additional work and references - will leave implementation for later
func (c *char) c4() {
	icd := 0
	c.Core.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if !(crit && (atk.Info.AttackTag == core.AttackTagElementalArt)) {
			return false
		}
		// Use an icd to make it only once per skill cast. Use 30 frames as two hits occur 20 frames apart
		if c.Core.Frame < icd {
			return false
		}
		icd = c.Core.Frame + 30

		c.AddEnergy("rosaria-c4", 5)
		c.coretype.Log.NewEvent("Rosaria C4 recovering 5 energy", coretype.LogEnergyEvent, c.Index, "new energy", c.Energy)
		return false
	}, "rosaria-c4")
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.coretype.Log.NewEvent("ActionStam not implemented", coretype.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}
