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

func init() {
	skillPressFrames = frames.InitAbilSlice(37 * 2)
}

const (
	skillPressCD        = 300
	skillHoldCD         = 360
	skillMarkKey        = "nahida-tri-karma"
	skillICDKey         = "nahida-tri-karma-icd"
	triKarmaParticleICD = "nahida-tri-karma-particle-icd"
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
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

	c.generateTriKarmaSnapshot()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 5),
		0,
		37*2, //TODO: snapshot frame and hitmark
		c.skillMarkTargets,
	)

	//reduce charge by 1
	c.SetCDWithDelay(action.ActionSkill, skillPressCD, 9)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	if hold > 300 {
		hold = 300
	}
	if hold < 30 {
		hold = 30 //TODO: hold appears to have a min before the camera appears
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

	c.generateTriKarmaSnapshot()

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 5),
		0,
		hold, //TODO: snapshot frame and hitmark
		c.skillMarkTargets,
	)

	//reduce charge by 1
	c.SetCDWithDelay(action.ActionSkill, skillHoldCD, 9)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return hold }, //TODO: ??? this is prob not right
		AnimationLength: hold,
		CanQueueAfter:   hold, // earliest cancel
		State:           action.SkillState,
	}

}

func (c *char) particlesOnDmg(_ combat.AttackCB) {
	if c.StatusIsActive(triKarmaParticleICD) {
		return
	}
	c.AddStatus(triKarmaParticleICD, 7*60, false) //TODO: assumed this icd is 7 seconds for now
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Dendro, c.ParticleDelay)
}

func (c *char) skillMarkTargets(a combat.AttackCB) {
	//TODO: unsure if it's 8 target globally or 8 target per cast
	//assuming globally for now
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	//assuming refresh if already exists; don't need to check for limits
	//in this case assuming shouldn't have been marked in the first place if
	//limit is exceeded
	if t.StatusIsActive(skillMarkKey) {
		//TODO: assumed this mark is affected by hitlag
		t.AddStatus(skillMarkKey, 1500, true)
		return
	}
	//TODO: this code is kinda inefficient
	count := 0
	for _, v := range c.Core.Combat.Enemies() {
		e, ok := v.(*enemy.Enemy)
		if !ok {
			continue
		}
		if e.StatusIsActive(skillMarkKey) {
			count++
		}
	}
	if count < 8 {
		t.AddStatus(skillMarkKey, 1200, true)
	}
}

// TODO: this implementation will only affect the next icd; not sure
// if it cuts short current as well
func (c *char) triKarmaInterval() int {
	if c.electroCount > 0 && c.Core.Status.Duration(burstKey) > 0 {
		cd := int((2.5 - burstTriKarmaCDReduction[c.electroCount-1][c.TalentLvlBurst()]) * 60)
		c.Core.Log.NewEvent("tri-karam cd reduced", glog.LogCharacterEvent, c.Index).Write("cooldown", cd)
		return cd

	}
	return int(2.5 * 60)
}

func (c *char) generateTriKarmaSnapshot() {
	//TODO: assuming tri karma snapshot her stats??

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Tri-Karma Purification",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNahidaSkill,
		ICDGroup:   combat.ICDGroupNahidaSkill,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       triKarmaAtk[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	em := snap.Stats[attributes.EM]
	ai.FlatDmg = em * triKarmaEM[c.TalentLvlSkill()]

	c.triKarmaSnapshot = combat.AttackEvent{
		Info:        ai,
		Snapshot:    snap,
		SourceFrame: c.Core.F,
	}

	if em > 200 {
		dmgBuff, crBuff := c.a4(em)
		c.triKarmaSnapshot.Snapshot.Stats[attributes.DmgP] += dmgBuff
		c.triKarmaSnapshot.Snapshot.Stats[attributes.CR] += crBuff
	}
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
	c.AddStatus(skillICDKey, c.triKarmaInterval(), true) //TODO: this is affected by hitlag?
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
			cb = c.particlesOnDmg
			done = true
		}
		c.Core.QueueAttackWithSnap(
			c.triKarmaSnapshot.Info,
			c.triKarmaSnapshot.Snapshot,
			combat.NewDefSingleTarget(e.Key()),
			1,
			cb,
		)
	}

}
