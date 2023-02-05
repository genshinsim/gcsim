package alhaitham

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillTapFrames []int
var skillHoldFrames []int

const (
	skillTapHitmark  = 19
	skillHoldHitmark = 28
	projectionICDKey = "alhaitham-projection-icd"
	particleICDKey   = "alhaitham-particle-icd"
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
	skillTapFrames[action.ActionBurst] = 28
	skillTapFrames[action.ActionDash] = 33
	skillTapFrames[action.ActionJump] = 33
	skillTapFrames[action.ActionSwap] = 36

	// skill (hold)-> x
	skillHoldFrames = frames.InitAbilSlice(87)
	skillHoldFrames[action.ActionAttack] = 86
	skillHoldFrames[action.ActionLowPlunge] = 35
	skillHoldFrames[action.ActionSkill] = 80
	skillHoldFrames[action.ActionWalk] = 86
	skillHoldFrames[action.ActionSwap] = 85

}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	if hold == 1 {
		return c.SkillHold()
	}

	c.Core.Tasks.Add(c.skillMirrorGain, 15)

	ai := combat.AttackInfo{
		Abil:               "Universality: An Elaboration on Form",
		ActorIndex:         c.Index,
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               rushAtk[c.TalentLvlSkill()],
		FlatDmg:            rushEm[c.TalentLvlSkill()] * c.Stat(attributes.EM),
		HitlagHaltFrames:   0.04 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), combat.Point{Y: 1}, 2.25), skillTapHitmark, skillTapHitmark)

	c.SetCDWithDelay(action.ActionSkill, 18*60, 15)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillTapFrames),
		AnimationLength: skillTapFrames[action.InvalidAction],
		CanQueueAfter:   skillTapFrames[action.ActionAttack], // earliest cancel
		State:           action.SkillState,
	}
}
func (c *char) SkillHold() action.ActionInfo {
	c.Core.Tasks.Add(c.skillMirrorGain, 23)

	ai := combat.AttackInfo{
		Abil:               "Universality: An Elaboration on Form (Hold)",
		ActorIndex:         c.Index,
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               rushAtk[c.TalentLvlSkill()],
		FlatDmg:            rushEm[c.TalentLvlSkill()] * c.Stat(attributes.EM),
		HitlagHaltFrames:   0.04 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), combat.Point{Y: 2}, 2.25), skillHoldHitmark, skillHoldHitmark)

	c.SetCDWithDelay(action.ActionSkill, 18*60, 23)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionLowPlunge], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillMirrorGain() {
	if c.mirrorCount == 0 { //extra mirror if 0 when cast
		c.mirrorGain(2)
		return
	}
	c.mirrorGain(1)

}
func (c *char) mirrorGain(generated int) {
	if generated == 0 {
		return
	}

	if c.mirrorCount == 0 {
		c.lastInfusionSrc = c.Core.F
		c.Core.Tasks.Add(c.mirrorLoss(c.Core.F, 1), 234)
		c.Core.Log.NewEvent("infusion added", glog.LogCharacterEvent, c.Index)

	}

	c.mirrorCount += generated
	if c.Base.Cons >= 2 { //triggers on overflow
		c.c2(generated)
	}

	if c.mirrorCount > 3 { //max 3 mirrors at a time.
		if c.Base.Cons >= 6 {
			c.c6(c.mirrorCount - 3)
		}
		c.mirrorCount = 3
		if c.Core.F != c.lastInfusionSrc { //this avoids multiple queues of mirror loss if mirror overflow multiple times in same frame
			c.lastInfusionSrc = c.Core.F
			c.Core.Tasks.Add(c.mirrorLoss(c.Core.F, 1), 234)
		}
		c.Core.Log.NewEvent("mirror overflowed", glog.LogCharacterEvent, c.Index).
			Write("mirrors gained", generated).
			Write("current mirrors", c.mirrorCount)

		return
	}
	c.Core.Log.NewEvent(fmt.Sprintf("Gained %v mirror(s)", generated), glog.LogCharacterEvent, c.Index).
		Write("current mirrors", c.mirrorCount)

}

func (c *char) mirrorLoss(src int, consumed int) func() {
	return func() {
		if consumed <= 0 {
			return
		}
		if c.lastInfusionSrc != src {
			c.Core.Log.NewEvent("mirror decrease ignored, src diff", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.lastInfusionSrc)
			return
		}
		if c.mirrorCount == 0 { //just in case
			c.Core.Log.NewEvent("Mirror count is 0, omitting reduction", glog.LogCharacterEvent, c.Index)
			return
		}

		c.mirrorCount -= consumed
		if c.mirrorCount < 0 { //This shouldn't happen but just in case
			c.mirrorCount = 0
		}

		c.Core.Log.NewEvent(fmt.Sprintf("Consumed %v mirror(s)", consumed), glog.LogCharacterEvent, c.Index).
			Write("current mirrors", c.mirrorCount)

		// queue up again if we still have mirrors
		if c.mirrorCount > 0 {
			c.Core.Tasks.Add(c.mirrorLoss(src, 1), 214) //not affected by hitlag, 448-234
		}

	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Dendro, c.ParticleDelay)
}

func (c *char) projectionAttack(a combat.AttackCB) {

	ae := a.AttackEvent
	//ignore if projection on icd
	if c.StatusIsActive(projectionICDKey) {
		return
	}
	//ignore if alhaitham is not on field
	if c.Core.Player.Active() != c.Index {
		return
	}
	//ignore if it doesn't have at least a mirror
	if c.mirrorCount == 0 {
		return
	}
	//ignore if it isn't NA/CA/Plunge
	if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra && ae.Info.AttackTag != combat.AttackTagPlunge {
		return
	}
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}

	var c1cb combat.AttackCBFunc
	if c.Base.Cons >= 1 {
		c1cb = c.c1
	}

	snapshotTiming := snapshotTimings[c.mirrorCount-1]
	strikeType := combat.StrikeTypeSlash
	if c.mirrorCount == 3 {
		strikeType = combat.StrikeTypeSpear
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Chisel-Light Mirror: Projection Attack %v", c.mirrorCount),
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupAlhaithamProjectionAttack,
		StrikeType: strikeType,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       mirrorAtk[c.TalentLvlSkill()],
		FlatDmg:    mirrorEm[c.TalentLvlSkill()] * c.Stat(attributes.EM),
	}

	player := c.Core.Combat.Player()
	var ap combat.AttackPattern
	var mirrorsHitmark []int
	switch c.mirrorCount {
	case 3:
		ap = combat.NewCircleHitOnTarget(player, combat.Point{Y: 4}, 4)
		mirrorsHitmark = mirror3Hitmarks
	case 2:
		ap = combat.NewCircleHitOnTargetFanAngle(player, combat.Point{Y: -0.1}, 5.5, 180)
		mirrorsHitmark = mirror2HitmarksLeft
		if c.Core.Rand.Float64() < 0.5 { //50% of using right/left hitmark frames
			mirrorsHitmark = mirror2HitmarksRight
		}
	default:
		ap = combat.NewBoxHitOnTarget(player, nil, 7, 3)
		mirrorsHitmark = mirror1HitmarkLeft
		if c.Core.Rand.Float64() < 0.5 { //50% of using right/left hitmark frames
			mirrorsHitmark = mirror1HitmarkRight
		}
	}

	for i := 0; i < c.mirrorCount; i++ {
		c.Core.QueueAttack(ai, ap, snapshotTiming, mirrorsHitmark[i], c1cb, c.particleCB)
	}
	c.AddStatus(projectionICDKey, 96, true) //1.6 sec icd

}
