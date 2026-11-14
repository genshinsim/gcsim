package vodyanitsa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillHitmark   = 34
	skillFirstTick = 190
	skillFirstHeal = 139
	skillDur       = 1042
	particleICDKey = "vodyanitsa-particle-icd"
	skillKey       = "vodyanitsa-skill"
)

func init() {
	skillFrames = frames.InitAbilSlice(59)
	skillFrames[action.ActionAttack] = 57
	skillFrames[action.ActionCharge] = 57
	skillFrames[action.ActionSkill] = 58
	skillFrames[action.ActionBurst] = 58
	skillFrames[action.ActionDash] = 57
	skillFrames[action.ActionWalk] = 57
	skillFrames[action.ActionSwap] = 57
}

// Reveals her true Water Imp form, dealing AoE Hydro DMG and gaining the Song of Ages Past effect:
// While this effect is active, Vodyanitsa summons the ancient Horn of Spring's Call, periodically
// moving closer to an opponent in a swift movement to attack them, dealing Hydro DMG.
// The above two types of attacks will also decrease the Hydro RES and Cryo RES of opponents hit.
// Additionally, while the Song of Ages Past effect is active, your active character gains increased
// resistance to interruption and periodically regenerates HP based on Vodyanitsa's Max HP.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Rechitativ: Sonorous Dawn",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		UseHP:      true,
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5.5)

	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB, c.skillShredCB)

	c.c2OnSkillAttack()

	c.skillSrc = c.Core.F
	c.Core.Tasks.Add(func() { c.skillHealTask(c.skillSrc) }, skillFirstHeal)
	c.Core.Tasks.Add(func() { c.skillAttackTask(c.skillSrc) }, skillHitmark+skillFirstTick)

	c.AddStatus(skillKey, skillDur+c.c2SkillDurBonus(), false)

	c.SetCD(action.ActionSkill, 16*60)

	c.a1OnSkill(c.skillSrc)
	c.a4OnSkill()
	c.c6OnSkill()

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}, nil
}

func (c *char) skillAttackTask(src int) {
	if !c.StatusIsActive(skillKey) {
		return
	}

	if c.skillSrc != src {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Horn of Spring's Call",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		UseHP:      true,
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3)

	c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB, c.skillShredCB)

	c.c2OnSkillAttack()
	c.Core.Tasks.Add(func() { c.skillAttackTask(src) }, 176)
}

func (c *char) skillHealTask(src int) {
	if !c.StatusIsActive(skillKey) {
		return
	}

	if c.skillSrc != src {
		return
	}

	maxhp := c.MaxHP()

	c.c4BeforeHeal()

	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  c.Core.Player.Active(),
		Message: "Song of Ages Past",
		Src:     skillHealRatio[c.TalentLvlSkill()]*maxhp + skillHealFlat[c.TalentLvlSkill()],
		Bonus:   c.Stat(attributes.Heal),
	})
	c.c1OnHeal()

	c.Core.Tasks.Add(func() { c.skillHealTask(src) }, 88)
}

func (c *char) skillShredCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	e, ok := ac.Target.(*enemy.Enemy)
	if !ok {
		return
	}

	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag(skillKey+"-hydro", 6*60),
		Ele:   attributes.Hydro,
		Value: -skillShred[c.TalentLvlSkill()],
	})

	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag(skillKey+"-cryo", 6*60),
		Ele:   attributes.Cryo,
		Value: -skillShred[c.TalentLvlSkill()],
	})
}

func (c *char) particleCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, 0.2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Hydro, c.ParticleDelay)
}
