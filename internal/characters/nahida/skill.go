package nahida

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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

func (c *char) Skill(p map[string]int) action.ActionInfo {
	c.markCount = 0
	if p["hold"] == 0 {
		return c.skillPress(p)
	} else {
		return c.skillHold(p)
	}
}

func (c *char) skillPress(p map[string]int) action.ActionInfo {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "All Schemes to Know (Press)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 0.2}, 4.6),
		0, //TODO: snapshot delay?
		skillPressHitmark,
		c.skillMarkTargets,
	)

	c.SetCDWithDelay(action.ActionSkill, skillPressCD, 11)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	// earliest hold can be let go is roughly 16.5. max is set to 317 so that
	// it aligns with max cd at 330
	if hold > 317 {
		hold = 317
	}
	if hold < 17 {
		hold = 17
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "All Schemes to Know (Hold)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 25),
		0, // TODO: snapshot timing
		hold+3,
		c.skillMarkTargets,
	)

	c.SetCDWithDelay(action.ActionSkill, skillHoldCD, hold-17+30)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return hold - 17 + skillHoldFrames[next] },
		AnimationLength: hold - 17 + 30 + skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   hold - 17 + 30 + skillHoldFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}

}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
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
	//assuming it's mark per skill cast; in this case refresh regardless if
	//target already marked up until 8
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

func (c *char) triKarmaOnReaction(rx event.Event) func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		c.triggerTriKarmaDamageIfAvail(t)
		return false
	}
}

func (c *char) triKarmaOnBloomDamage(args ...interface{}) bool {
	t, ok := args[0].(*enemy.Enemy)
	if !ok {
		return false
	}
	//only on bloom, burgeon, hyperbloom damage
	ae, ok := args[1].(*combat.AttackEvent)
	if !ok {
		return false
	}
	switch ae.Info.AttackTag {
	case combat.AttackTagBloom:
	case combat.AttackTagHyperbloom:
	case combat.AttackTagBurgeon:
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
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNahidaSkill,
			ICDGroup:   combat.ICDGroupNahidaSkill,
			StrikeType: combat.StrikeTypeDefault,
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
			4,
			cb,
		)
	}

}
