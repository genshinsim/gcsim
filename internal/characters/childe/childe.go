package childe

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("tartaglia", NewChar)
}

const rtA1 = 18 * 60 // riptide duration lasts 18 sec

// childe specific character implementation
type char struct {
	*character.Tmpl
	eCast         int // the frame childe cast E to enter melee stance
	rtParticleICD int
	rtFlashICD    []int
	rtSlashICD    []int
	rtExpiry      []int
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
	if c.Base.Cons >= 4 {
		c.mlBurstUsed = false
		c.funcC4 = make([]bool, len(c.Core.Targets))
	}
	if c.Base.Cons >= 6 {
		c.c6 = false
	}

	c.rtParticleICD = 0
	c.rtFlashICD = make([]int, len(c.Core.Targets))
	c.rtSlashICD = make([]int, len(c.Core.Targets))
	c.rtExpiry = make([]int, len(c.Core.Targets))

	c.Core.Flags.ChildeActive = true
	c.onExitField()
	c.onDefeatTargets()
	c.rtParticleGen()
	c.rtHook()
	return &c, nil
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
		if c.Core.Status.Duration("childemelee") > 0 {
			c.onExitMeleeStance()
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
			for i := 0; i < len(c.Core.Targets); i++ {
				if c.rtExpiry[i] < c.Core.F {
					c.Core.Log.Debugw("Childe applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", i, "Expiry", c.Core.F+rtA1)
				}
				c.rtExpiry[i] = c.Core.F + rtA1
			}
			c.Core.Log.Debugw("Riptide Burst ticked", "frame", c.Core.F, "event", core.LogCharacterEvent)
		}, "Riptide Burst", 5)
		//re-index riptide expiry frame array if needed

		if c.Base.Cons >= 2 {
			c.AddEnergy(4)
			c.Core.Log.Debugw("Childe C2 restoring 4 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "new energy", c.Energy)
		}
		return false
	}, "childe-riptide-burst")
}
