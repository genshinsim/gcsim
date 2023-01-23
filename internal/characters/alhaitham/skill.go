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

var skillTapFrames []int
var skillHoldFrames []int

//TODO: defhalt values are needed

const (
	skillTapHitmark  = 19
	skillHoldHitmark = 28
	projectionICDKey = "alhaitham-projection-icd"
)

var mirror1HitmarkLeft = []int{39}
var mirror1HitmarkRight = []int{40}

var mirror2HitmarksLeft = []int{28, 37}
var mirror2HitmarksRight = []int{26, 35}

var mirror3Hitmarks = []int{32, 41, 51}

var snapshotTimings = []int{20, 22, 26}

func init() {
	// skill (tap) -> x
	skillTapFrames = frames.InitAbilSlice(44)
	skillTapFrames[action.ActionAttack] = 27
	skillTapFrames[action.ActionSkill] = 28
	skillTapFrames[action.ActionDash] = 33
	skillTapFrames[action.ActionJump] = 33
	skillTapFrames[action.ActionSwap] = 36

	// skill (hold)-> x
	skillHoldFrames = frames.InitAbilSlice(86)
	skillHoldFrames[action.ActionAttack] = 86
	skillHoldFrames[action.ActionLowPlunge] = 35
	skillHoldFrames[action.ActionSkill] = 86
	skillHoldFrames[action.ActionDash] = 86
	skillHoldFrames[action.ActionJump] = 86
	skillHoldFrames[action.ActionSwap] = 86

}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	if hold == 1 {
		return c.SkillHold()
	}

	c.Core.Tasks.Add(func() {
		if c.mirrorCount == 0 { //extra mirror if 0 when cast
			c.mirrorGain()
		}
		c.mirrorGain()
	}, 15)

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
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), combat.Point{Y: 1}, 2.25), skillTapHitmark, skillTapHitmark)

	c.SetCDWithDelay(action.ActionSkill, 18*60, 15)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillTapFrames),
		AnimationLength: skillTapFrames[action.InvalidAction],
		CanQueueAfter:   skillTapFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
func (c *char) SkillHold() action.ActionInfo {
	c.Core.Tasks.Add(func() {
		if c.mirrorCount == 0 { //extra mirror if 0 when cast
			c.mirrorGain()
		}
		c.mirrorGain()
	}, 23)

	ai := combat.AttackInfo{
		Abil:               "Universality: An Elaboration on Form (Hold)",
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
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), combat.Point{Y: 2}, 2.25), skillHoldHitmark, skillHoldHitmark)

	c.SetCDWithDelay(action.ActionSkill, 18*60, 23)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) mirrorGain() {
	if c.Base.Cons >= 2 { //triggers on overflow
		c.c2()
	}
	if c.mirrorCount > 2 { //max 3 mirrors at a time.
		queueOnFrame := false              //var tracks if this is the first overflowing in this frame
		if c.Core.F == c.lastInfusionSrc { //check if c.lastinfusion has already been called on this frame
			queueOnFrame = true
		}
		if !queueOnFrame { //this avoids multiple queues of mirror loss if mirror overflow multiple times in same frame
			c.lastInfusionSrc = c.Core.F
			c.Core.Tasks.Add(c.mirrorLoss(c.Core.F), 234)
		}
		c.Core.Log.NewEvent("mirror overflowed", glog.LogCharacterEvent, c.Index)

		if c.Base.Cons >= 6 {
			c.c6()
		}
		return
	}
	if c.mirrorCount == 0 {
		c.lastInfusionSrc = c.Core.F
		c.Core.Tasks.Add(c.mirrorLoss(c.Core.F), 234)
		c.Core.Log.NewEvent("infusion added", glog.LogCharacterEvent, c.Index)

	}
	c.mirrorCount++
	c.recentlyMirrorGain = true
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
			if c.recentlyMirrorGain { //if mirror has been gained recently, mirror is lost after 234f
				c.Core.Tasks.Add(c.mirrorLoss(src), 234) //not affected by hitlag
				return
			}
			c.Core.Tasks.Add(c.mirrorLoss(src), 214) //not affected by hitlag, 448-234

		}
		c.recentlyMirrorGain = false
	}
}

func (c *char) projectionAttack(a combat.AttackCB) {

	ae := a.AttackEvent
	//ignore if projection on icd
	if c.StatusIsActive(projectionICDKey) {
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
	mirrorsHitmark := make([]int, 3)
	snapshotTiming := 21
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
	switch c.mirrorCount {
	case 3:
		ai.Mult = mirror1Atk[c.TalentLvlSkill()]
		ai.FlatDmg = mirror1Em[c.TalentLvlSkill()] * c.Stat(attributes.EM)
		ap = combat.NewCircleHitOnTarget(trg, combat.Point{Y: 4}, 4)
		mirrorsHitmark = mirror3Hitmarks
		snapshotTiming = snapshotTimings[2]
	case 2:
		ai.Mult = mirror1Atk[c.TalentLvlSkill()]
		ai.FlatDmg = mirror1Em[c.TalentLvlSkill()] * c.Stat(attributes.EM)
		ap = combat.NewCircleHitOnTargetFanAngle(trg, combat.Point{Y: -0.1}, 5.5, 180)
		snapshotTiming = snapshotTimings[1]
		mirrorsHitmark = mirror2HitmarksLeft
		if c.Core.Rand.Float64() < 0.5 { //50% of using right/left hitmark frames
			mirrorsHitmark = mirror2HitmarksRight
		}
	default:
		snapshotTiming = snapshotTimings[0]
		mirrorsHitmark = mirror1HitmarkLeft
		if c.Core.Rand.Float64() < 0.5 { //50% of using right/left hitmark frames
			mirrorsHitmark = mirror1HitmarkRight
		}
	}

	for i := 0; i < c.mirrorCount; i++ {
		c.Core.QueueAttack(ai, ap, snapshotTiming, mirrorsHitmark[i], c1cb)
	}

	c.Core.QueueParticle("alhaitham", 1, attributes.Dendro, c.ParticleDelay)
	c.AddStatus(projectionICDKey, 96, true) //1.6 sec icd

}
