package xiao

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
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
	c6Src        int
	c6Count      int
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
	if c.Base.Cons >= 6 {
		c.c6Count = 0
		c.c6()
	}

	c.Tags["eCharge"] = c.eChargeMax
	c.onExitField()

	return &c, nil
}

// Implements Xiao C2:
// When in the party and not on the field, Xiao's Energy Recharge is increased by 25%
func (c *char) c2() {
	stat_mod := make([]float64, core.EndStatType)
	stat_mod[core.ER] = 0.25
	c.AddMod(core.CharStatMod{
		Key:    "xiao-c2",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if c.Core.ActiveChar != c.Index {
				return stat_mod, true
			}
			return nil, false
		},
	})
}

// Implements Xiao C6:
// While under the effect of Bane of All Evil, hitting at least 2 opponents with Xiao's Plunge Attack will immediately grant him 1 charge of Lemniscatic Wind Cycling, and for the next 1s, he may use Lemniscatic Wind Cycling while ignoring its CD.
// Adds an OnDamage event checker - if we record two or more instances of plunge damage, then activate C6
func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		if !((ds.Abil == "High Plunge") || (ds.Abil == "Low Plunge")) {
			return false
		}
		if c.Core.Status.Duration("xiaoburst") == 0 {
			return false
		}
		// Stops after reaching 2 hits on a single plunge.
		// Plunge frames are greater than duration of C6 so this will always refresh properly.
		if c.Core.Status.Duration("xiaoc6") > 0 {
			return false
		}
		if c.c6Src != ds.SourceFrame {
			c.c6Src = ds.SourceFrame
			c.c6Count = 0
			return false
		}

		c.c6Count++

		// Prevents activation more than once in a single plunge attack
		if c.c6Count == 2 {
			c.recoverCharge(c.eTickSrc)
			c.eTickSrc = c.Core.F

			c.Core.Status.AddStatus("xiaoc6", 60)
			c.Core.Log.Debugw("Xiao C6 activated", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "new E charges", c.Tags["eCharge"], "expiry", c.Core.F+60)

			c.c6Count = 0
			return false
		}
		return false
	}, "xiao-c6")
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
