package tartaglia

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Tartaglia, NewChar)
}

const riptideDuration = 18 * 60 // riptide duration lasts 18 sec

// tartaglia specific character implementation
type char struct {
	*character.Tmpl
	eCast         int // the frame tartaglia casts E to enter melee stance
	rtParticleICD int
	// rtFlashICD    []int
	// rtSlashICD    []int
	// rtExpiry      []int
	mlBurstUsed bool // used for c6
}

//constants for tags
const (
	riptideKey         = "riptide"
	riptideSlashICDKey = "riptide-slash-icd"
	riptideFlashICDKey = "riptide-flash-icd"
)

// Initializes character
func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro
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
	// c.rtFlashICD = make([]int, len(c.Core.Targets))
	// c.rtSlashICD = make([]int, len(c.Core.Targets))
	// c.rtExpiry = make([]int, len(c.Core.Targets))

	c.Core.Flags.ChildeActive = true
	c.onExitField()
	c.onDefeatTargets()
	// c.applyRT()
	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionCharge:
		return 20
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Key.String())
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

//Riptide Burst: Defeating an opponent affected by Riptide creates a Hydro burst
//that inflicts the Riptide status on nearby opponents hit.
// Handles Childe riptide burst and C2 on death effects
func (c *char) onDefeatTargets() {
	c.Core.Events.Subscribe(core.OnTargetDied, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		//do nothing if no riptide on target
		if t.GetTag(riptideKey) < c.Core.F {
			return false
		}
		c.AddTask(func() {
			ai := core.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Riptide Burst",
				AttackTag:  core.AttackTagNormal,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				StrikeType: core.StrikeTypeDefault,
				Element:    core.Hydro,
				Durability: 50,
				Mult:       rtBurst[c.TalentLvlAttack()],
			}
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 0)
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
// func (c *char) applyRT() {
// 	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
// 		atk := args[1].(*core.AttackEvent)
// 		t := args[0].(core.Target)
// 		crit := args[3].(bool)

// 		if c.Core.Status.Duration("tartagliamelee") > 0 {
// 			if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
// 				return false
// 			}
// 			if !crit {
// 				return false
// 			}

// 			//dont log if it just refreshes riptide status
// 			if c.rtExpiry[t.Index()] <= c.Core.F {
// 				c.Core.Log.Debugw("Tartaglia applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t.Index(), "rtExpiry", c.Core.F+riptideDuration)
// 			}
// 			c.rtExpiry[t.Index()] = c.Core.F + riptideDuration
// 		} else {
// 			if atk.Info.AttackTag != core.AttackTagElementalBurst && atk.Info.AttackTag != core.AttackTagExtra {
// 				return false
// 			}

// 			//ranged burst or aim mode
// 			//dont log if it just refreshes riptide status
// 			if c.rtExpiry[t.Index()] <= c.Core.F {
// 				c.Core.Log.Debugw("Tartaglia applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t.Index(), "rtExpiry", c.Core.F+riptideDuration)
// 			}
// 			c.rtExpiry[t.Index()] = c.Core.F + riptideDuration
// 		}

// 		return false
// 	}, "tartaglia-apply-riptide")
// }
