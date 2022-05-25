package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Tartaglia, NewChar)
}

const riptideDuration = 18 * 60 // riptide duration lasts 18 sec

// tartaglia specific character implementation
type char struct {
	*character.Tmpl
	eCast         int // the frame tartaglia casts E to enter melee stance
	rtParticleICD int
	mlBurstUsed   bool // used for c6
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

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.SkillCon = 3
	c.BurstCon = 5
	c.NormalHitNum = 6

	c.eCast = 0
	c.rtParticleICD = 0

	if c.Base.Cons >= 6 {
		c.mlBurstUsed = false
	}

	c.Core.Flags.ChildeActive = true

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.onExitField()
	c.onDefeatTargets()
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionCharge:
		return 20
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

// Hook to end Tartaglia's melee stance prematurely if he leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("tartagliamelee") > 0 {
			//TODO: need to verify if this is correct
			//but if childe is currently in melee stance and skill is on CD that means that
			//the button has lit up yet from original skill press
			//in which case we need to reset the cooldown first
			c.ResetActionCooldown(core.ActionSkill)
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
			c.AddEnergy("tartaglia-c2", 4)
		}
		return false
	}, "tartaglia-on-enemy-death")
}

func (c *char) c4(t core.Target) {
	if t.GetTag(riptideKey) < c.Core.F {
		return
	}

	if c.Core.Status.Duration("tartagliamelee") > 0 {
		c.rtSlashTick(t)
	} else {
		c.rtFlashTick(t)
	}

	c.AddTask(func() { c.c4(t) }, "tartaglia-c4", 60*4)
}
