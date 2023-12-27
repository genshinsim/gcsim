package charlotte

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillPressFrames []int
	skillHoldFrames  []int
)

func init() {
	skillPressFrames = frames.InitAbilSlice(49) // E -> CA
	skillPressFrames[action.ActionAttack] = 42
	skillPressFrames[action.ActionSkill] = 48
	skillPressFrames[action.ActionBurst] = 42
	skillPressFrames[action.ActionDash] = 42
	skillPressFrames[action.ActionJump] = 42
	skillPressFrames[action.ActionWalk] = 42
	skillPressFrames[action.ActionSwap] = 41

	skillHoldFrames = frames.InitAbilSlice(137) // hE -> Dash, Jump
	skillHoldFrames[action.ActionAttack] = 132
	skillHoldFrames[action.ActionCharge] = 132
	skillHoldFrames[action.ActionSkill] = 130
	skillHoldFrames[action.ActionBurst] = 130
	skillHoldFrames[action.ActionWalk] = 136
	skillHoldFrames[action.ActionSwap] = 134
}

const (
	skillPressBoxX = 4
	skillPressBoxY = 8
	// hold hitbox is an approximation
	skillHoldOffsetX        = 1.3
	skillHoldBoxX           = 4.5
	skillHoldBoxY           = 30
	skillPressCD            = 720
	skillHoldCD             = 1080
	skillPressHitmark       = 31
	skillHoldHitmark        = 111
	skillPressDelay         = 29
	skillHoldDelay          = 110
	skillPressParticleCount = 3
	skillHoldParticleCount  = 5
	skillMarkKey            = "charlotte-mark"
	skillPressMarkKey       = "charlotte-e"
	skillHoldMarkKey        = "charlotte-hold-e"
)

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.markCount = 0
	if p["hold"] == 0 {
		return c.skillPress()
	}
	return c.skillHold(p)
}

func (c *char) skillPress() (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Framing: Freezing Point Composition",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, skillPressBoxX, skillPressBoxY)

	c.Core.QueueAttack(
		ai,
		ap,
		0,
		skillPressHitmark,
		c.skillPressParticleCB,
		c.skillPressMarkTargets,
		c.makeC2CB(),
	)

	c.SetCDWithDelay(action.ActionSkill, skillPressCD, skillPressDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHold(p map[string]int) (action.Info, error) {
	hold := p["hold"]
	// earliest hold hitmark is ~111f
	// latest hold hitmark is ~919f
	// hold=1 gives 111f and hold=809 gives a 919f delay until hitmark.
	if hold < 1 {
		hold = 1
	}
	if hold > 809 {
		hold = 809
	}
	hold += skillHoldHitmark
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Framing: Freezing Point Composition (Hold)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{X: skillHoldOffsetX}, skillHoldBoxX, skillHoldBoxY)

	c.Core.QueueAttack(
		ai,
		ap,
		0,
		hold,
		c.skillHoldParticleCB,
		c.skillHoldMarkTargets,
		c.makeC2CB(),
	)

	c.SetCDWithDelay(action.ActionSkill, skillHoldCD, hold-2)

	return action.Info{
		Frames:          func(next action.Action) int { return hold + skillHoldFrames[next] },
		AnimationLength: hold + skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   hold + skillHoldFrames[action.ActionBurst],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillPressParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	c.Core.QueueParticle(c.Base.Key.String(), skillPressParticleCount, attributes.Cryo, c.ParticleDelay)
}

func (c *char) skillHoldParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	c.Core.QueueParticle(c.Base.Key.String(), skillHoldParticleCount, attributes.Cryo, c.ParticleDelay)
}

func (c *char) skillPressMarkTargets(a combat.AttackCB) {
	if c.markCount == 5 {
		return
	}
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	c.markCount++

	if t.StatusIsActive(skillPressMarkKey) {
		t.DeleteStatus(skillPressMarkKey)
	}
	if t.StatusIsActive(skillHoldMarkKey) {
		t.DeleteStatus(skillHoldMarkKey)
	}

	t.SetTag(skillMarkKey, c.Core.F)
	t.AddStatus(skillPressMarkKey, 360+0.8*60, true)
	t.QueueEnemyTask(c.skillPressMark(c.Core.F, t), 1.5*60)
}

func (c *char) skillHoldMarkTargets(a combat.AttackCB) {
	if c.markCount == 5 {
		return
	}
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	c.markCount++

	if t.StatusIsActive(skillPressMarkKey) {
		t.DeleteStatus(skillPressMarkKey)
	}
	if t.StatusIsActive(skillHoldMarkKey) {
		t.DeleteStatus(skillHoldMarkKey)
	}

	t.SetTag(skillMarkKey, c.Core.F)
	t.AddStatus(skillHoldMarkKey, 720+0.9*60, true)
	t.QueueEnemyTask(c.skillHoldMark(c.Core.F, t), 1.5*60)
}

func (c *char) skillPressMark(src int, t *enemy.Enemy) func() {
	return func() {
		if src != t.GetTag(skillMarkKey) {
			return
		}
		if !t.StatusIsActive(skillPressMarkKey) {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Snappy Silhouette Mark",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagCharlotteMark,
			ICDGroup:   attacks.ICDGroupCharlotteMark,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skillPressMark[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(t.Key()), 0, 0)
		c.Core.Tasks.Add(c.skillPressMark(src, t), 1.5*60)
	}
}

func (c *char) skillHoldMark(src int, t *enemy.Enemy) func() {
	return func() {
		if src != t.GetTag(skillMarkKey) {
			return
		}
		if !t.StatusIsActive(skillHoldMarkKey) {
			return
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Focused Impression Mark",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagCharlotteMark,
			ICDGroup:   attacks.ICDGroupCharlotteMark,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skillHoldMark[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(t.Key()), 0, 0)
		c.Core.Tasks.Add(c.skillHoldMark(src, t), 1.5*60)
	}
}
