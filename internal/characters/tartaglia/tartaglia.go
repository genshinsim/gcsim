package tartaglia

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("tartaglia", NewChar)
}

const rtA1 = 18 * 60 // riptide duration lasts 18 sec

// tartaglia specific character implementation
type char struct {
	*character.Tmpl
	eCast         int // the frame tartaglia casts E to enter melee stance
	rtParticleICD int
	rtFlashICD    []int
	rtSlashICD    []int
	rtExpiry      []int
	mlBurstUsed   bool // used for c6
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
	if c.Base.Cons >= 6 {
		c.mlBurstUsed = false
	}

	c.rtParticleICD = 0
	c.rtFlashICD = make([]int, len(c.Core.Targets))
	c.rtSlashICD = make([]int, len(c.Core.Targets))
	c.rtExpiry = make([]int, len(c.Core.Targets))

	c.Core.Flags.ChildeActive = true
	c.onExitField()
	c.onDefeatTargets()
	c.applyRT()
	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionCharge:
		return 20
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}

// Hook to end Tartaglia's melee stance prematurely if he leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			c.onExitMeleeStance()
		}
		return false
	}, "tartaglia-exit")
}

// Handles Childe riptide burst and C2 on death effects
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

			c.Core.Log.Debugw("Riptide Burst ticked", "frame", c.Core.F, "event", core.LogCharacterEvent)
		}, "Riptide Burst", 5)
		//TODO: re-index riptide expiry frame array if needed

		if c.Base.Cons >= 2 {
			c.AddEnergy(4)
			c.Core.Log.Debugw("Tartaglia C2 restoring 4 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "new energy", c.Energy)
		}
		return false
	}, "tartaglia-on-enemy-death")
}

//apply riptide status to enemy hit
func (c *char) applyRT() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		crit := args[3].(bool)

		if c.Core.Status.Duration("tartagliamelee") > 0 {
			if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
				return false
			}
			if !crit {
				return false
			}

			//dont log if it just refreshes riptide status
			if c.rtExpiry[t.Index()] <= c.Core.F {
				c.Core.Log.Debugw("Tartaglia applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t.Index(), "rtExpiry", c.Core.F+rtA1)
			}
			c.rtExpiry[t.Index()] = c.Core.F + rtA1
		} else {
			if ds.AttackTag != core.AttackTagElementalBurst && ds.AttackTag != core.AttackTagExtra {
				return false
			}

			//ranged burst or aim mode
			//dont log if it just refreshes riptide status
			if c.rtExpiry[t.Index()] <= c.Core.F {
				c.Core.Log.Debugw("Tartaglia applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t.Index(), "rtExpiry", c.Core.F+rtA1)
			}
			c.rtExpiry[t.Index()] = c.Core.F + rtA1
		}

		return false
	}, "tartaglia-apply-riptide")
}
