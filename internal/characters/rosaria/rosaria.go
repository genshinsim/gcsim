package rosaria

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type char struct {
	*character.Tmpl
}

func init() {
	core.RegisterCharFunc(keys.Rosaria, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Cryo
	c.Energy = 60
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
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if !crit {
			return false
		}
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		c.AddMod(core.CharStatMod{
			Key:    "rosaria-c1",
			Expiry: c.Core.F + 240,
			Amount: func(a core.AttackTag) ([]float64, bool) {
				val := make([]float64, core.EndStatType)
				val[core.AtkSpd] = 0.1
				val[core.DmgP] = 0.1
				if a != core.AttackTagNormal {
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
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if !(crit && (atk.Info.AttackTag == core.AttackTagElementalArt)) {
			return false
		}
		// Use an icd to make it only once per skill cast. Use 30 frames as two hits occur 20 frames apart
		if c.Core.F < icd {
			return false
		}
		icd = c.Core.F + 30

		c.AddEnergy(5)
		c.Core.Log.Debugw("Rosaria C4 recovering 5 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "new energy", c.Energy)
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
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Key.String())
		return 0
	}
}
