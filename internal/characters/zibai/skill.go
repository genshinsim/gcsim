package zibai

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillFrames        []int
	skillSkillFrames   []int
	skillStrideHitmark = []int{30, 30 + 4}
)

const (
	skillHitmark      = 18
	particleICDKey    = "zibai-particle-icd"
	skillKey          = "zibai-skill"
	radianceNAICDKey  = "zibai-radiance-na-icd"
	radianceLCrICDKey = "zibai-radiance-lcr-icd"
	maxRadiance       = 100

	skillAbil1 = "Spirit Steed's Stride"
	skillAbil2 = skillAbil1 + lunarCrystallizeAbil
)

func init() {
	skillFrames = frames.InitAbilSlice(52)
	skillFrames[action.ActionAttack] = 28
	skillFrames[action.ActionSkill] = 28
	skillFrames[action.ActionBurst] = 29
	skillFrames[action.ActionDash] = 27
	skillFrames[action.ActionJump] = 27
	skillFrames[action.ActionWalk] = 41

	skillSkillFrames = frames.InitAbilSlice(61)
	skillSkillFrames[action.ActionAttack] = 45
	skillSkillFrames[action.ActionSkill] = 54
	skillSkillFrames[action.ActionBurst] = 45
	skillSkillFrames[action.ActionDash] = 42
	skillSkillFrames[action.ActionJump] = 43
	skillSkillFrames[action.ActionWalk] = 47
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		// do nothing if previous char wasn't zibai
		prev := args[0].(int)
		if prev != c.Index() {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}

		c.DeleteStatus(skillKey)
	}, "zibai-exit")
}

// Summoning a shadow of her former powers, she switches to the Lunar Phase Shift mode.
// In this mode, Zibai's Normal Attacks and Charged Attacks will deal Geo DMG that cannot be
// overridden by other infusions, and she can accrue special Phase Shift Radiance through different
// methods. Zibai can consume Phase Shift Radiance to unleash the special Elemental Skill Spirit
// Steed's Stride.

// Moonsign: Ascendant Gleam
// While in the Lunar Phase Shift mode, when Zibai performs Normal Attacks, the fourth attack will
// deal an additional instance of Geo DMG, which is considered Lunar-Crystallize Reaction DMG.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillSkill()
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Lunar Phase Shift (0 dmg)",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1.3}, 3.1), skillHitmark, skillHitmark)

	c.AddStatus(skillKey, 15*60+skillHitmark, true)

	c.SetCDWithDelay(action.ActionSkill, 18*60, skillHitmark)

	c.skillsUsed = 0

	c.skillSrc = c.Core.F

	c.QueueCharTask(func() {
		c.radianceTicker(c.skillSrc)()
		c.a1OnSkill()
		c.c1OnSkill()
	}, skillHitmark)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

// Lunar Phase Shift
// After using the Elemental Skill Heaven and Earth Made Manifest, Zibai will change to this mode.
// This mode lasts at most 15s and has the following properties:
// · Normal and Charged Attacks are converted into Geo DMG that cannot be overridden by other
// 	 infusions.
// · The Elemental Skill Heaven and Earth Made Manifest will also be converted to the Elemental
//   Skill Spirit Steed's Stride: When Zibai has at least 70 Phase Shift Radiance, she can consume
//   70 Phase Shift Radiance to unleash Spirit Steed's Stride and cause two instances of Geo DMG.
//   The second instance will be considered Lunar-Crystallize Reaction DMG.

// A maximum of 100 Phase Shift Radiance points can be accrued, and Zibai can accumulate them in the
// following ways:
// · She gains 10 points every second she is in the Lunar Phase Shift mode.
// · She gains 5 points when hitting an opponent with her Normal Attacks. She can gain Phase Shift
//   Radiance points in this way once every 0.5s.

// Moonsign: Ascendant Gleam
// When nearby party members trigger Lunar-Crystallize, they will also accrue 35 Phase Shift
// Radiance points for Zibai. Points can be gained this way once every 4s. Unleashing the special
// elemental skill Spirit Steed's Stride will reset the CD for Phase Shift Radiance gained this way.

// Zibai will leave this mode after unleashing Spirit Steed's Stride 4 times, or when the skill's
// duration ends.
func (c *char) skillSkill() (action.Info, error) {
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 2.5}, 3.5)
	for i := range skillStride {
		c.QueueCharTask(func() {
			ai := info.AttackInfo{
				ActorIndex:       c.Index(),
				Abil:             skillAbil1,
				AttackTag:        attacks.AttackTagElementalArt,
				ICDTag:           attacks.ICDTagNone,
				ICDGroup:         attacks.ICDGroupDefault,
				StrikeType:       attacks.StrikeTypeBlunt,
				PoiseDMG:         75,
				Element:          attributes.Geo,
				UseDef:           true,
				Durability:       25,
				Mult:             skillStride[i][c.TalentLvlSkill()],
				HitlagHaltFrames: 0.02 * 60,
				HitlagFactor:     0.01,
			}

			if i == 1 {
				ai.Abil = skillAbil2
				ai.AttackTag = attacks.AttackTagDirectLunarCrystallize
				ai.Durability = 0
				ai.IgnoreDefPercent = 1
				ai.FlatDmg += c.a1StrideBonusDmg()
				ai.PoiseDMG = 50
			}

			c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB, c.c4SkillCB)
		}, skillStrideHitmark[i])
	}
	c.DeleteStatus(radianceLCrICDKey)
	c.skillsUsed += 1
	c.c6ConsumeRadiance()
	if c.skillsUsed >= c.c1MaxSkillsPerSkill() {
		c.DeleteStatus(skillKey)
	}

	return action.Info{
		Frames:          func(next action.Action) int { return skillSkillFrames[next] },
		AnimationLength: skillSkillFrames[action.InvalidAction],
		CanQueueAfter:   skillSkillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	if c.Core.Rand.Float64() > 0.67 {
		return
	}

	c.AddStatus(particleICDKey, 2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Geo, c.ParticleDelay)
}

func (c *char) radianceTicker(src int) func() {
	return func() {
		if c.skillSrc != src {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		c.addRadiance(1)
		c.Core.Tasks.Add(c.radianceTicker(src), 6)
	}
}

func (c *char) radianceCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(radianceNAICDKey) {
		return
	}

	c.AddStatus(radianceNAICDKey, 0.5*60, true)
	c.addRadiance(5)
}

func (c *char) skillInit() {
	c.Core.Events.Subscribe(event.OnLunarCrystallize, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}
		if c.StatusIsActive(radianceLCrICDKey) {
			return
		}
		c.AddStatus(radianceLCrICDKey, 4*60, true)
		c.addRadiance(35)
	}, "zibai-radiance-lcr")
}

func (c *char) addRadiance(amt float64) {
	amt *= c.c6RadianceEff()
	c.radiance = min(c.radiance+amt, maxRadiance)
	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent(fmt.Sprint("Gained ", amt, " radiance (", c.radiance, ")"), glog.LogCharacterEvent, c.Index())
	}
}

func (c *char) consumeRadiance(amt float64) {
	c.radiance = max(c.radiance-amt, 0)
	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent(fmt.Sprint("Consumed ", amt, " radiance (", c.radiance, ")"), glog.LogCharacterEvent, c.Index())
	}
}
