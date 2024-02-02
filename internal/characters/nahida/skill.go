package nahida

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillPressFrames []int
var skillHoldFrames []int

func init() {
	skillPressFrames = frames.InitAbilSlice(32)
	skillPressFrames[action.ActionAttack] = 28
	skillPressFrames[action.ActionCharge] = 28
	skillPressFrames[action.ActionSkill] = 32
	skillPressFrames[action.ActionBurst] = 32
	skillPressFrames[action.ActionDash] = 27
	skillPressFrames[action.ActionJump] = 26
	skillPressFrames[action.ActionSwap] = 25

	skillHoldFrames = frames.InitAbilSlice(63)
	skillHoldFrames[action.ActionAttack] = 57
	skillHoldFrames[action.ActionCharge] = 58
	skillHoldFrames[action.ActionSkill] = 62
	skillHoldFrames[action.ActionBurst] = 62
	skillHoldFrames[action.ActionDash] = 59
	skillHoldFrames[action.ActionJump] = 62
	skillHoldFrames[action.ActionSwap] = 57
}

const (
	skillPressCD        = 300
	skillHoldCD         = 360
	skillPressHitmark   = 13
	skillMarkKey        = "nahida-e"
	skillICDKey         = "nahida-e-icd"
	triKarmaParticleICD = "nahida-e-particle-icd"
)

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.markCount = 0
	if p["hold"] == 0 {
		return c.skillPress(), nil
	}
	return c.skillHold(p)
}

func (c *char) skillPress() action.Info {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "All Schemes to Know (Press)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 0.2}, 4.6),
		0, //TODO: snapshot delay?
		skillPressHitmark,
		c.skillMarkTargets,
	)

	c.SetCDWithDelay(action.ActionSkill, skillPressCD, 11)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int) (action.Info, error) {
	hold := p["hold"]
	// earliest hold can be let go is roughly 16.5, max is 317
	// adds the value in hold onto the minimum length of 16, so hold=1 gives 17f and hold=5 gives a 22f delay until hitmark.
	if hold > 300 {
		hold = 300
	}
	if hold < 1 {
		hold = 1
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "All Schemes to Know (Hold)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 25),
		0, // TODO: snapshot timing
		hold+3+16,
		c.skillMarkTargets,
	)

	c.SetCDWithDelay(action.ActionSkill, skillHoldCD, hold-17+30)

	return action.Info{
		Frames:          func(next action.Action) int { return hold - 17 + skillHoldFrames[next] },
		AnimationLength: hold - 17 + 30 + skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   hold - 17 + 30 + skillHoldFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(triKarmaParticleICD) {
		return
	}
	c.AddStatus(triKarmaParticleICD, 7*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Dendro, c.ParticleDelay)
}

func (c *char) skillMarkTargets(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	// assuming it's mark per skill cast; in this case refresh regardless if
	// target already marked up until 8
	if c.markCount < 8 {
		t.AddStatus(skillMarkKey, 1500, true)
		c.markCount++
	}
}

func (c *char) updateTriKarmaInterval() {
	cd := int(2.5 * 60)
	if c.electroCount > 0 && c.StatusIsActive(withinBurstKey) {
		cd -= int(burstTriKarmaCDReduction[c.electroCount-1][c.TalentLvlBurst()] * 60)
	}
	if cd != c.triKarmaInterval {
		c.Core.Log.NewEvent("tri-karma cd reduced", glog.LogCharacterEvent, c.Index).Write("cooldown", cd)
		c.triKarmaInterval = cd
	}
	c.QueueCharTask(c.updateTriKarmaInterval, 60) // check every 1s
}

func (c *char) triKarmaOnReaction(args ...interface{}) bool {
	t, ok := args[0].(*enemy.Enemy)
	if !ok {
		return false
	}
	c.triggerTriKarmaDamageIfAvail(t)
	return false
}

func (c *char) triKarmaOnBloomDamage(args ...interface{}) bool {
	t, ok := args[0].(*enemy.Enemy)
	if !ok {
		return false
	}
	// only on bloom, burgeon, hyperbloom damage
	ae, ok := args[1].(*combat.AttackEvent)
	if !ok {
		return false
	}
	switch ae.Info.AttackTag {
	case attacks.AttackTagBloom:
	case attacks.AttackTagHyperbloom:
	case attacks.AttackTagBurgeon:
	default:
		return false
	}

	c.triggerTriKarmaDamageIfAvail(t)
	return false
}

func (c *char) triggerTriKarmaDamageIfAvail(t *enemy.Enemy) {
	if c.StatusIsActive(skillICDKey) {
		return
	}
	if !t.StatusIsActive(skillMarkKey) {
		return
	}
	c.AddStatus(skillICDKey, c.triKarmaInterval, true) //TODO: this is affected by hitlag?
	done := false
	for _, v := range c.Core.Combat.Enemies() {
		e, ok := v.(*enemy.Enemy)
		if !ok {
			continue
		}
		if !e.StatusIsActive(skillMarkKey) {
			continue
		}
		var cb combat.AttackCBFunc
		if !done {
			cb = c.particleCB
			done = true
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tri-Karma Purification",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNahidaSkill,
			ICDGroup:   attacks.ICDGroupNahidaSkill,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Dendro,
			Durability: 25,
			Mult:       triKarmaAtk[c.TalentLvlSkill()],
		}
		snap := c.Snapshot(&ai)
		em := snap.Stats[attributes.EM]
		ai.FlatDmg = em * triKarmaEM[c.TalentLvlSkill()]

		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewSingleTargetHit(e.Key()),
			3,
			cb,
		)
	}
}
