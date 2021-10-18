package childe

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("tartaglia", NewChar)
}

// childe specific character implementation
type char struct {
	*character.Tmpl
	eCast         int // the frame childe cast E to enter melee stance
	caFrame       int // Idkw on kqm lib, childe has 2 different charged attack frames
	rtParticleICD int
	rtflashICD    []int // procs by aiming (the enemy must have riptide status)
	rtslashICD    []int // rt slash procs on normal, charge in melee form (the enemy must have riptide status)
	rtExpiry      []int // riptide expired frames
	rtA1          int   // riptide duration lasts 18 sec
	funcC4        []bool
	mlBurstUsed   bool // used for c4. After clearing riptide, remove c4 tickers
	c6            bool // if true reset E cd; otherwise not
}

// Initializes character
func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.SkillCon = 3
	c.BurstCon = 5
	c.NormalHitNum = 6
	c.eCast = 0
	c.caFrame = 73
	c.rtA1 = 18 * 60
	if c.Base.Cons >= 4 {
		c.mlBurstUsed = false
	}
	if c.Base.Cons >= 6 {
		c.c6 = false
	}

	c.Core.Flags.ChildeActive = true
	c.initRTICD()
	c.onExitField()
	c.onDefeatTargets()
	c.rtParticleGen()
	c.rtHook()
	return &c, nil
}

func (c *char) initRTICD() {
	c.rtParticleICD = 0
	c.rtflashICD = make([]int, len(c.Core.Targets))
	c.rtslashICD = make([]int, len(c.Core.Targets))
	c.rtExpiry = make([]int, len(c.Core.Targets))

	if c.Base.Cons >= 4 {
		c.funcC4 = make([]bool, len(c.Core.Targets))
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionCharge:
		return 20
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}

// Hook to end Childe's melee stance prematurely if he leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.CharIndex() {
			return false
		}

		if c.Core.Status.Duration("childemelee") > 0 {
			c.Core.Status.DeleteStatus("childemelee")
			//TODO: preemptive cd doesnt seem fixed at 6s
			newCD := float64(c.Core.F - c.eCast + 6*60)
			//Foul Legacy: Tide Withholder. Decreases the CD of Foul Legacy: Raging Tide by 20%
			if c.Base.Cons >= 1 {
				newCD *= 0.8
			}
			if c.Base.Cons >= 6 && c.c6 {
				newCD = 0
				c.c6 = false
			}
			c.Core.Log.Debugw("Childe leaving melee stance", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
				c.Core.Status.Duration("childemelee"))
			c.SetCD(core.ActionSkill, int(newCD))
		}
		return false
	}, "childe-exit")
}

func (c *char) onDefeatTargets() {
	c.Core.Events.Subscribe(core.OnTargetDied, func(args ...interface{}) bool {
		c.AddTask(func() {
			d := c.Snapshot(
				"Riptide Burst",
				core.AttackTagNormal,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Hydro,
				50,
				rtBurst[c.TalentLvlAttack()],
			)
			d.Targets = core.TargetAll

			c.Core.Combat.ApplyDamage(&d)

			// apply riptide status
			for _, t := range c.Core.Targets {
				if c.rtExpiry[t.Index()] < c.Core.F {
					c.Core.Log.Debugw("Childe applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t.Index(), "Expiry", c.Core.F+c.rtA1)
				}
				c.rtExpiry[t.Index()] = c.Core.F + c.rtA1
			}
			c.Core.Log.Debugw("Riptide Burst ticked", "frame", c.Core.F, "event", core.LogCharacterEvent)
		}, "Riptide Burst", 5)
		//TODO: re-index riptide expired frame list

		if c.Base.Cons >= 2 {
			c.AddEnergy(4)
			c.Core.Log.Debugw("Childe C2 restoring 4 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "new energy", c.Energy)
		}
		return false
	}, "childe-riptide-burst")
}
