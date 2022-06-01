package xiao

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Xiao, NewChar)
}

// Xiao specific character implementation
type char struct {
	*character.Tmpl
	qStarted int
	a4Expiry int
	c6Src    int
	c6Count  int
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
	c.Base.Element = core.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 70
	}
	c.Energy = float64(e)
	c.EnergyMax = 70
	c.Weapon.Class = core.WeaponClassSpear
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 6

	c.c6Count = 0

	c.SetNumCharges(core.ActionSkill, 2)
	if c.Base.Cons >= 1 {
		c.SetNumCharges(core.ActionSkill, 3)
	}

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()
	c.InitCancelFrames()

	c.onExitField()

	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}
}

func (c *char) a4() {
	c.AddMod(core.CharStatMod{
		Key:    "xiao-a4",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			stacks := c.Tags["a4"]
			if stacks == 0 {
				return nil, false
			}
			m[core.DmgP] += float64(stacks) * 0.15
			return m, true
		},
	})
}

// Implements Xiao C2:
// When in the party and not on the field, Xiao's Energy Recharge is increased by 25%
func (c *char) c2() {
	c.AddMod(core.CharStatMod{
		Key:    "xiao-c2",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			m := make([]float64, core.EndStatType)
			m[core.ER] = 0.25
			if c.Core.ActiveChar != c.Index {
				return m, true
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
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if !((atk.Info.Abil == "High Plunge") || (atk.Info.Abil == "Low Plunge")) {
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
		if c.c6Src != atk.SourceFrame {
			c.c6Src = atk.SourceFrame
			c.c6Count = 0
			return false
		}

		c.c6Count++

		// Prevents activation more than once in a single plunge attack
		if c.c6Count == 2 {
			c.ResetActionCooldown(core.ActionSkill)

			c.Core.Status.AddStatus("xiaoc6", 60)
			c.Core.Log.NewEvent("Xiao C6 activated", core.LogCharacterEvent, c.Index, "new E charges", c.Tags["eCharge"], "expiry", c.Core.F+60)

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
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

// Xiao specific Snapshot implementation for his burst bonuses. Similar to Hu Tao
// Implements burst anemo attack damage conversion and DMG bonus
// Also implements A1:
// While under the effects of Bane of All Evil, all DMG dealt by Xiao is increased by 5%. DMG is increased by an additional 5% for every 3s the ability persists. The maximum DMG Bonus is 25%
func (c *char) Snapshot(a *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(a)

	if c.Core.Status.Duration("xiaoburst") > 0 {
		// Calculate and add A1 damage bonus - applies to all damage
		// Fraction dropped in int conversion in go - acts like floor
		stacks := 1 + int((c.Core.F-c.qStarted)/180)
		if stacks > 5 {
			stacks = 5
		}
		ds.Stats[core.DmgP] += float64(stacks) * 0.05
		c.Core.Log.NewEvent("a1 adding dmg %", core.LogCharacterEvent, c.Index, "stacks", stacks, "final", ds.Stats[core.DmgP], "time since burst start", c.Core.F-c.qStarted)

		// Anemo conversion and dmg bonus application to normal, charged, and plunge attacks
		// Also handle burst CA ICD change to share with Normal
		switch a.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
			a.ICDTag = core.ICDTagNormalAttack
		case core.AttackTagPlunge:
		default:
			return ds
		}
		a.Element = core.Anemo
		bonus := burstBonus[c.TalentLvlBurst()]
		ds.Stats[core.DmgP] += bonus
		c.Core.Log.NewEvent("xiao burst damage bonus", core.LogCharacterEvent, c.Index, "bonus", bonus, "final", ds.Stats[core.DmgP])
	}
	return ds
}
