package chevreuse

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const (
	skillPressCDStart      = 18
	skillPressHitmark      = 25
	skillPressArkheHitmark = 59

	skillHoldCDStart      = 13
	skillHoldHitmark      = 19
	skillHoldArkheHitmark = 55

	skillHealKey      = "chev-skill-heal"
	skillHealInterval = 120
	particleICDKey    = "chev-particle-icd"
	arkheICDKey       = "chev-arkhe-icd"
)

var skillPressFrames []int
var skillHoldFrames []int

func init() {

	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(31) // E -> N1/Q
	skillPressFrames[action.ActionAttack] = 31
	skillPressFrames[action.ActionBurst] = 31
	skillPressFrames[action.ActionDash] = 23
	skillPressFrames[action.ActionJump] = 25
	skillPressFrames[action.ActionWalk] = 24
	skillPressFrames[action.ActionSwap] = 24

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(26) // E -> Q
	skillHoldFrames[action.ActionAttack] = 25
	skillHoldFrames[action.ActionBurst] = 26
	skillHoldFrames[action.ActionDash] = 21
	skillHoldFrames[action.ActionJump] = 23
	skillHoldFrames[action.ActionWalk] = 24
	skillHoldFrames[action.ActionSwap] = 23
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if p["hold"] == 0 {
		return c.skillPress(), nil
	}
	return c.skillHold(p), nil
}

func (c *char) skillPress() action.Info {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Short-Range Rapid Interdiction Fire",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -0.5}, 2, 6),
		skillPressHitmark,
		skillPressHitmark,
		c.particleCB,
		c.arkhe(skillPressArkheHitmark-skillPressHitmark),
	)

	c.skillHeal(skillPressCDStart)
	c.SetCDWithDelay(action.ActionSkill, 15*60, skillPressCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int) action.Info {

	hold := p["hold"]
	// earliest hold hitmark is ~19f
	// latest hold hitmark is ~319f
	// hold=1 gives 19f and hold=301 gives a 319f delay until hitmark.
	if hold < 1 {
		hold = 1
	}
	if hold > 301 {
		hold = 301
	}
	// subtract 1 to account for needing to supply > 0 to indicate hold
	hold -= 1
	hitmark := hold + skillHoldHitmark
	cdStart := hold + skillHoldCDStart

	var ai combat.AttackInfo
	var ap combat.AttackPattern

	if c.overChargedBall {
		ai = combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Short-Range Rapid Interdiction Fire [Overcharged]",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			Element:    attributes.Pyro,
			Durability: 25,
			PoiseDMG:   125,
			Mult:       skillOvercharged[c.TalentLvlSkill()],
		}

		ap = combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5)
		// remove status once overcharged is ball shot
		c.overChargedBall = false
		c.Core.Tasks.Add(c.a4, cdStart)

	} else {
		ai = combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Short-Range Rapid Interdiction Fire [Hold]",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       skillHold[c.TalentLvlSkill()],
		}
		ap = combat.NewBoxHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -0.5}, 3, 7)

	}

	// c4
	if c.StatModIsActive(c4StatusKey) {
		c.c4ShotsLeft -= 1
		if c.c4ShotsLeft == 0 {
			c.DeleteStatus(c4StatusKey)
		}
	} else {
		c.SetCDWithDelay(action.ActionSkill, 15*60, cdStart)
	}

	c.Core.QueueAttack(
		ai,
		ap,
		hitmark,
		hitmark,
		c.particleCB,
		c.c2(),
		c.arkhe(skillHoldArkheHitmark-skillHoldHitmark),
	)

	c.skillHeal(cdStart)

	return action.Info{
		Frames:          func(next action.Action) int { return hold + skillHoldFrames[next] },
		AnimationLength: hold + skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   hold + skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) arkhe(delay int) combat.AttackCBFunc {
	// triggers on hitting anything, not just enemy
	return func(a combat.AttackCB) {
		if c.StatusIsActive(arkheICDKey) {
			return
		}
		c.AddStatus(arkheICDKey, 10*60, true)

		aiArkhe := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Surging Blade (" + c.Base.Key.Pretty() + ")",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 0,
			Mult:       arkhe[c.TalentLvlSkill()],
		}

		c.Core.QueueAttack(
			aiArkhe,
			combat.NewCircleHitOnTarget(a.Target.Pos(), nil, 2),
			delay,
			delay,
		)
	}
}

func (c *char) skillHeal(delay int) {
	skillDur := 12*60 + 1 // heal on last tick of expiry
	if !c.StatusIsActive(skillHealKey) {
		c.Core.Tasks.Add(func() {
			c.AddStatus(skillHealKey, skillDur, false)               // not hitlag extendable
			c.Core.Tasks.Add(c.startSkillHealing, skillHealInterval) // first heal comes after 2s
			// don't queue up c6 team heal if there is one already queued up
			if c.c6HealQueued {
				return
			}
			c.c6HealQueued = true
			c.Core.Tasks.Add(c.c6TeamHeal, 12*60)
		}, delay)
		return
	}
	// extend skill heal on retrigger while still active (c4+)
	c.Core.Tasks.Add(func() {
		c.ExtendStatus(skillHealKey, skillDur)
	}, delay)
}

func (c *char) startSkillHealing() {
	if !c.StatusIsActive(skillHealKey) {
		return
	}

	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Core.Player.Active(),
		Message: "Short-Range Rapid Interdiction Fire Healing",
		Src:     skillHpRegen[c.TalentLvlSkill()]*c.MaxHP() + skillHpFlat[c.TalentLvlSkill()],
		Bonus:   c.Stat(attributes.Heal),
	})
	c.c6(c.Core.Player.ActiveChar())
	c.Core.Tasks.Add(c.startSkillHealing, skillHealInterval)
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 10*60, false) // chev has 10s particle icd, not hitlag extendable
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Pyro, c.ParticleDelay)
}

func (c *char) overchargedBallEventSub() {
	c.Core.Events.Subscribe(event.OnOverload, func(args ...interface{}) bool {
		// don't proc on gadgets
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}
		c.overChargedBall = true
		return false
	}, "chev-overcharged-ball")
}
