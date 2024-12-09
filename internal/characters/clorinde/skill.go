package clorinde

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames     []int
	skillDashFrames []int
)

const (
	skillStateKey  = "clorinde-night-watch"
	particleICDKey = "clorinde-particle-icd"

	skillDashHitmark = 11
	tolerance        = 0.0000001
	skillStart       = 6
	skillCD          = 16 * 60
)

func init() {
	skillFrames = frames.InitAbilSlice(33) // E -> Q
	skillFrames[action.ActionAttack] = 31
	skillFrames[action.ActionSkill] = 32
	skillFrames[action.ActionDash] = skillStart // ability doesn't start if dash is done before CD
	skillFrames[action.ActionJump] = 25
	skillFrames[action.ActionSwap] = 25
	skillFrames[action.ActionWalk] = 32

	skillDashFrames = frames.InitAbilSlice(43) // E -> Walk
	skillDashFrames[action.ActionAttack] = 24
	skillDashFrames[action.ActionSkill] = 24
	skillDashFrames[action.ActionBurst] = 24
	skillDashFrames[action.ActionDash] = 25
	skillDashFrames[action.ActionJump] = 25
	skillDashFrames[action.ActionSwap] = 42
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// first press activates skill state
	// sequential presses pew pew stuff
	if c.StatusIsActive(skillStateKey) {
		return c.skillDash(p)
	}
	c.AddStatus(skillStateKey, skillStart+int(60*skillStateDuration[0]), true)
	c.QueueCharTask(c.c6skill, 0)
	c.SetCDWithDelay(action.ActionSkill, skillCD, skillStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillDash(p map[string]int) (action.Info, error) {
	c.normalSCounter = 0

	// depending on BOL lvl it does either 1 hit or 3 hit
	ratio := c.currentHPDebtRatio()
	switch {
	case ratio >= 1:
		if c.Base.Cons >= 6 && c.c6Stacks > 0 {
			c.c6()
			c.c6Stacks -= 1
		}
		return c.skillDashFullBOL(p)
	case math.Abs(ratio) < tolerance:
		return c.skillDashNoBOL(p)
	default:
		return c.skillDashRegular(p)
	}
}

func (c *char) gainBOLOnAttack() {
	c.ModifyHPDebtByRatio(skillBOLGain[c.TalentLvlSkill()])
}

func (c *char) skillDashNoBOL(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Impale the Night (0% BoL)",
		AttackTag:      attacks.AttackTagNormal,
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeSlash,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           skillLungeNoBOL[c.TalentLvlSkill()],
		HitlagFactor:   0.01,
		IgnoreInfusion: true,
	}
	// TODO: what's the size of this??
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.6)
	c.Core.QueueAttack(ai, ap, skillDashHitmark, skillDashHitmark, c.particleCB)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillDashFrames),
		AnimationLength: skillDashFrames[action.InvalidAction],
		CanQueueAfter:   skillDashFrames[action.ActionSkill],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillDashFullBOL(_ map[string]int) (action.Info, error) {
	for i := 0; i < 3; i++ {
		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Impale the Night (100%+ BoL)",
			AttackTag:      attacks.AttackTagNormal,
			ICDTag:         attacks.ICDTagNormalAttack,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeSlash,
			Element:        attributes.Electro,
			Durability:     25,
			Mult:           skillLungeFullBOL[c.TalentLvlSkill()],
			HitlagFactor:   0.01,
			IgnoreInfusion: true,
		}
		// TODO: what's the size of this??
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.8)
		c.Core.QueueAttack(ai, ap, skillDashHitmark, skillDashHitmark, c.particleCB)
	}

	// Bond of Life timing is ping dependent
	c.QueueCharTask(func() {
		c.skillHeal(skillLungeFullBOLHeal[0], "Impale the Night (100%+ BoL)")
	}, skillDashHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillDashFrames),
		AnimationLength: skillDashFrames[action.InvalidAction],
		CanQueueAfter:   skillDashFrames[action.ActionSkill],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillDashRegular(_ map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Impale the Night (<100% BoL)",
		AttackTag:      attacks.AttackTagNormal,
		ICDTag:         attacks.ICDTagNormalAttack,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeSlash,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           skillLungeLowBOL[c.TalentLvlSkill()],
		HitlagFactor:   0.01,
		IgnoreInfusion: true,
	}
	// TODO: what's the size of this??
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.8)
	c.Core.QueueAttack(ai, ap, skillDashHitmark, skillDashHitmark, c.particleCB)

	// Bond of Life timing is ping dependent
	c.QueueCharTask(func() {
		c.skillHeal(skillLungeLowBOLHeal[0], "Impale the Night (<100% BoL)")
	}, skillDashHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillDashFrames),
		AnimationLength: skillDashFrames[action.InvalidAction],
		CanQueueAfter:   skillDashFrames[action.ActionSkill],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHeal(bolMult float64, msg string) {
	amt := c.CurrentHPDebt() * bolMult
	c.heal(&info.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: msg,
		Src:     amt,
		Bonus:   c.Stat(attributes.Heal), // TODO: confirms that it scales with healing %
	})
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 2*60, true)

	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Electro, c.ParticleDelay)
}
