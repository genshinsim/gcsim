package chevreuse

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillPressFrames []int
var skillHoldFrames []int

const (
	skillPressCDStart = 16
	skillPressHitmark = 17
	skillPressTravel  = 1

	skillHoldCDStart = 16
	skillHoldHitmark = 12
	skillHoldTravel  = 3

	skillKey       = "chev-heal"
	particleICDKey = "chev-particle-icd"
)

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(39) // E -> N1/Q
	skillPressFrames[action.ActionDash] = 34
	skillPressFrames[action.ActionJump] = 35
	skillPressFrames[action.ActionWalk] = 19
	skillPressFrames[action.ActionSwap] = 37

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(46) // E -> Swap
	skillHoldFrames[action.ActionAttack] = 38
	skillHoldFrames[action.ActionBurst] = 37
	skillHoldFrames[action.ActionDash] = 30
	skillHoldFrames[action.ActionJump] = 30
	skillHoldFrames[action.ActionWalk] = 30
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if p["hold"] != 0 {
		return c.skillHold(), nil
	}
	return c.skillPress(), nil
}

func (c *char) skillPress() action.Info {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Short-Range Rapid Interdiction Fire",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skillPress[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.4}, 2.5, 10),
		skillPressHitmark,
		skillPressHitmark+skillPressTravel,
		c.makeParticleCB(),
	)

	skillDur := 12 * 60
	skillHealInterval := 2
	c.Core.Tasks.Add(func() {
		c.Core.Status.Add(skillKey, skillDur)
		c.Core.Tasks.Add(c.startSkillHealing(), skillHealInterval*60)
		c.Core.Tasks.Add(c.c6TeamHeal(), 12*60) // Assuming this executes every 90 frames = 1.5s
	}, 23)

	c.SetCDWithDelay(action.ActionSkill, 15*60, skillPressCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionWalk], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold() action.Info {

	var ai combat.AttackInfo
	if c.overChargedBall {
		ai = combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Short-Range Rapid Interdiction Fire [Overcharged]",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       skillHold[c.TalentLvlSkill()],
		}
		// remove status once overcharged is ball shot
		c.a4()
		c.overChargedBall = false

	} else {
		ai = combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Short-Range Rapid Interdiction Fire [Hold]",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       skillHold[c.TalentLvlSkill()],
		}

	}

	skillDur := 12 * 60
	skillHealInterval := 2
	c.Core.Tasks.Add(func() {
		c.Core.Status.Add(skillKey, skillDur)
		c.Core.Tasks.Add(c.startSkillHealing(), skillHealInterval*60)
		c.Core.Tasks.Add(c.c6TeamHeal(), skillDur)
	},
		23)

	c.Core.QueueAttack(
		ai,
		combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()),
		skillHoldHitmark,
		skillHoldHitmark+skillHoldTravel,
		c.makeParticleCB(),
		c.C2(),
	)

	// c4
	if c.StatModIsActive(c4StatusKey) {
		c.c4ShotsLeft -= 1
		if c.c4ShotsLeft == 0 {
			c.DeleteStatus(c4StatusKey)
		}
	} else {
		c.SetCDWithDelay(action.ActionSkill, 15*60, skillHoldCDStart+1)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}

		if c.StatusIsActive(particleICDKey) {
			return
		}

		c.AddStatus(particleICDKey, 10*60, false) // chev has 10s particle icd
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Pyro, c.ParticleDelay)
	}
}

func (c *char) startSkillHealing() func() {

	return func() {
		if c.Core.Status.Duration(skillKey) == 0 {
			return
		}

		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Short-Range Rapid Interdiction Fire Healing",
			Src:     skillHpRegen[c.TalentLvlBurst()]*c.MaxHP() + skillHpFlat[c.TalentLvlBurst()],
			Bonus:   c.Stat(attributes.Heal),
		})
		c.c6(c.Core.Player.ActiveChar())
		c.Core.Tasks.Add(c.startSkillHealing(), 120)
	}
}
