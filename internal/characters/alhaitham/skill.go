package alhaitham

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const skillHitmark = 25

func init() {
	// skill -> x
	skillFrames = frames.InitAbilSlice(44)
	skillFrames[action.ActionAttack] = 44
	skillFrames[action.ActionSkill] = 44
	skillFrames[action.ActionDash] = 30
	skillFrames[action.ActionJump] = 30
	skillFrames[action.ActionSwap] = 30

}

func (c *char) Skill(p map[string]int) action.ActionInfo {

	if c.mirrorCount == 0 { //extra mirror if 0 when cast
		c.mirrorGain()
	}
	c.mirrorGain()
	ai := combat.AttackInfo{
		Abil:               "Universality: An Elaboration on Form",
		ActorIndex:         c.Index,
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               rushAtk[c.TalentLvlSkill()],
		FlatDmg:            rushEm[c.TalentLvlSkill()] * c.Stat(attributes.EM),
		HitlagHaltFrames:   0.09 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}
	//TODO: Add hold support
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), combat.Point{Y: 1}, 2.25), skillHitmark, skillHitmark)

	c.SetCDWithDelay(action.ActionSkill, 18*60, 20) //TODO: delay value if needed on cast

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) mirrorGain() {
	if c.Base.Cons >= 2 { //triggers on overflow
		c.c2()
	}
	if c.mirrorCount > 2 { //max 3 mirrors at a time.
		c.lastInfusionSrc = c.Core.F
		c.Core.Tasks.Add(c.mirrorLoss(c.Core.F), 4*60)
		c.Core.Log.NewEvent("mirror overflowed", glog.LogCharacterEvent, c.Index)

		if c.Base.Cons >= 6 {
			c.c6()
		}
		return
	}
	if c.mirrorCount == 0 {
		c.lastInfusionSrc = c.Core.F
		c.Core.Tasks.Add(c.mirrorLoss(c.Core.F), 4*60)
		c.Core.Log.NewEvent("infusion added", glog.LogCharacterEvent, c.Index)

	}
	c.mirrorCount++
	c.Core.Log.NewEvent("Gained 1 mirror", glog.LogCharacterEvent, c.Index)

}

func (c *char) mirrorLoss(src int) func() {
	return func() {
		if c.lastInfusionSrc != src {
			c.Core.Log.NewEvent("mirror decrease ignored, src diff", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.lastInfusionSrc)
			return
		}
		if c.mirrorCount == 0 { //for case when you swap out and have mirrorloss queue'd. T
			//TODO:change the lastinfusionsrc on swap event perhaps?
			c.Core.Log.NewEvent("Mirror count is 0, ommiting reduction", glog.LogCharacterEvent, c.Index)
			return
		}

		c.mirrorCount--

		c.Core.Log.NewEvent("Lost 1 mirror", glog.LogCharacterEvent, c.Index).
			Write("mirrors", c.mirrorCount)

		// queue up again if we still have mirrors
		if c.mirrorCount > 0 {
			c.Core.Tasks.Add(c.mirrorLoss(src), 4*60) //not affected by hitlag
		}
	}
}

func (c *char) projectionAttack(a combat.AttackCB) {

	ae := a.AttackEvent
	//ignore if projection on icd
	if c.projectionICD > c.Core.F {
		return
	}
	//ignore if it doesn't have at least a mirror
	if c.mirrorCount == 0 {
		return
	}
	//ignore it isn't NA/CA/Plunge
	if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra && ae.Info.AttackTag != combat.AttackTagPlunge {
		return
	}
	var c1cb combat.AttackCBFunc
	if c.Base.Cons >= 1 {
		c1cb = c.c1
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Chisel-Light Mirror: Projection Attack %v", c.mirrorCount),
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagAlhaithamProjectionAttack,
		ICDGroup:   combat.ICDGroupAlhaithamProjectionAttack,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       mirror1Atk[c.TalentLvlSkill()],
		FlatDmg:    mirror1Em[c.TalentLvlSkill()] * c.Stat(attributes.EM),
	} //TODO: hitlag stuff?
	trg, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	ap := combat.NewBoxHitOnTarget(trg, nil, 7, 3)
	//TODO: clean this code later (currently redundant)
	switch c.mirrorCount {
	case 3:
		ai.Mult = mirror1Atk[c.TalentLvlSkill()]
		ai.FlatDmg = mirror1Em[c.TalentLvlSkill()] * c.Stat(attributes.EM)
		ap = combat.NewCircleHitOnTarget(trg, combat.Point{Y: 4}, 4)
	case 2:
		ai.Mult = mirror1Atk[c.TalentLvlSkill()]
		ai.FlatDmg = mirror1Em[c.TalentLvlSkill()] * c.Stat(attributes.EM)
		ap = combat.NewCircleHitOnTargetFanAngle(trg, combat.Point{Y: -0.1}, 5.5, 180)
	default:

	}

	for i := 0; i < c.mirrorCount; i++ {
		c.Core.QueueAttack(ai, ap, 5, 5, c1cb) //TODO: projection hit timings
	}

	c.Core.QueueParticle("alhaitham", 1, attributes.Dendro, c.ParticleDelay)
	c.projectionICD = c.Core.F + 96 //1.6 sec icd
	return

}
