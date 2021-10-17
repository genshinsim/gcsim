package tartaglia

import (
	"strings"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

type char struct {
	*character.Tmpl

	// Tracks riptide application and last application frame on each target
	riptideStatusLastProc   []int
	riptideSlashLastProc    []int
	riptideFlashLastProc    []int
	riptideParticleLastProc int

	c6BurstUsed bool

	// Required to track for skill cooldown time calculation
	skillStartFrame int
}

func init() {
	core.RegisterCharFunc("tartaglia", NewChar)
}

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
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 6

	// Initialize these at some negative value to avoid issues at the start of the sim
	c.riptideStatusLastProc = make([]int, len(s.Targets))
	c.riptideSlashLastProc = make([]int, len(s.Targets))
	c.riptideFlashLastProc = make([]int, len(s.Targets))
	c.riptideParticleLastProc = -9999

	for i := range c.riptideStatusLastProc {
		c.riptideStatusLastProc[i] = -9999
	}
	for i := range c.riptideSlashLastProc {
		c.riptideSlashLastProc[i] = -9999
	}
	for i := range c.riptideFlashLastProc {
		c.riptideFlashLastProc[i] = -9999
	}

	c.c6BurstUsed = false

	c.applyRiptide()
	c.onDeathChecks()
	c.onExitField()

	return &c, nil
}

// On damage hook to check for riptide application
// TODO: Grouping all riptide applications here instead of putting some as on hit callbacks results in some minor performance issues - should refactor despite slightly higher complexity
// Sources are: Ranged CA, Ranged Burst, and A4
// When Tartaglia is in Foul Legacy: Raging Tide’s Melee Stance, on dealing a CRIT hit, Normal and Charged Attacks apply the Riptide status effect to opponents.
func (c *char) applyRiptide() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		target := args[0].(core.Target)
		ds := args[1].(*core.Snapshot)
		crit := args[3].(bool)

		if ds.ActorIndex != c.Index {
			return false
		}

		flagRiptide := false
		source := "ERROR"

		// Ranged Burst
		if (ds.AttackTag == core.AttackTagElementalBurst) && (ds.Abil == "Ranged Burst: Flash of Havoc") {
			flagRiptide = true
			source = "Ranged Burst"
		} else if (ds.AttackTag == core.AttackTagExtra) && (ds.Abil == "Aimed Shot") {
			// Ranged CA
			flagRiptide = true
			source = "Aimed Shot"
		} else if ((strings.Contains(ds.Abil, "Melee Normal")) || (ds.Abil == "Melee Charge Attack")) && crit {
			// A4 - crits on melee stance NA/CA
			flagRiptide = true
			source = "Melee Attack/CA"
		}

		if flagRiptide {
			c.Core.Log.Debugw("Riptide status inflicted from "+source, "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "target", target.Index(), "expiry", c.Core.F+18*60)
			c.riptideStatusLastProc[target.Index()] = c.Core.F
		}

		return false
	}, "tartaglia-a4")
}

// Adds event hooker for Riptide Burst and C2
// Riptide Burst: Defeating an opponent affected by riptide creates a Hydro burst that inflicts the Riptide status on nearby opponents hit.
// C2: When opponents affected by Riptide are defeated, Tartaglia regenerates 4 Elemental Energy
func (c *char) onDeathChecks() {
	c.Core.Events.Subscribe(core.OnTargetDied, func(args ...interface{}) bool {
		t := args[0].(core.Target)

		// Stop if not affected by Riptide
		if c.riptideStatusLastProc[t.Index()]+18*60 < c.Core.F {
			return false
		}

		// Proc Riptide burst - hydro damage on death
		d := c.Snapshot(
			"Riptide Burst",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypePierce,
			core.Hydro,
			50,
			riptideBurst[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll

		c.Core.Combat.ApplyDamage(&d)

		// C2 - recovers energy
		if c.Base.Cons >= 2 {
			c.AddEnergy(4)
			c.Core.Log.Debugw("Tartaglia C2 recovering 4 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "char", c.Index)
		}
		return false
	}, "tartaglia-onDeathChecks")
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		if prev == c.Index && c.Core.Status.Duration("tartagliamelee") > 0 {
			c.onExitMeleeStance()
		}
		return false
	}, "tartaglia-exit")
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 20
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Name)
		return 0
	}
}
