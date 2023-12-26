package charlotte

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

// TODO aoe
var (
	skillPressFrames []int
	skillHoldFrames  []int
)

func init() {
	skillPressFrames = frames.InitAbilSlice(42) // E -> Walk
	skillPressFrames[action.ActionAttack] = 42
	skillPressFrames[action.ActionCharge] = 42
	skillPressFrames[action.ActionSkill] = 49
	skillPressFrames[action.ActionBurst] = 41
	skillPressFrames[action.ActionDash] = 43
	skillPressFrames[action.ActionJump] = 42
	skillPressFrames[action.ActionSwap] = 52

	skillHoldFrames = frames.InitAbilSlice(135) // hE -> Walk
	skillHoldFrames[action.ActionAttack] = 134
	skillHoldFrames[action.ActionCharge] = 131
	skillHoldFrames[action.ActionSkill] = 134
	skillHoldFrames[action.ActionBurst] = 131
	skillHoldFrames[action.ActionDash] = 137
	skillHoldFrames[action.ActionJump] = 138
	skillHoldFrames[action.ActionSwap] = 146
}

const (
	skillPressRadius        = 6
	skillPressAngle         = 120
	skillHoldRadius         = 6
	skillPressCD            = 720
	skillHoldCD             = 1080
	skillPressHitmark       = 32
	skillHoldHitmark        = 112
	skillPressDelay         = 30
	skillHoldDelay          = 111
	skillPressParticleCount = 3
	skillHoldParticleCount  = 5
	skillPressMarkKey       = "charlotte-e"
	skillHoldMarkKey        = "charlotte-hold-e"
)

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if p["hold"] == 1 {
		return c.skillHold()
	}
	return c.skillPress()
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

	ap := combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, skillPressRadius, skillPressAngle)
	if c.Base.Cons >= 2 {
		c.c2(ap)
	}

	c.Core.QueueAttack(
		ai,
		ap,
		0,
		skillPressHitmark,
		c.skillPressParticleCB,
		c.skillPressMarkTargets,
	)

	c.SetCDWithDelay(action.ActionSkill, skillPressCD, skillPressDelay)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionBurst],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHold() (action.Info, error) {
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

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, skillHoldRadius)
	if c.Base.Cons >= 2 {
		c.c2(ap)
	}

	c.Core.QueueAttack(
		ai,
		ap,
		0,
		skillHoldHitmark,
		c.skillHoldParticleCB,
		c.skillHoldMarkTargets,
	)

	c.SetCDWithDelay(action.ActionSkill, skillHoldCD, skillHoldDelay)

	return action.Info{
		Frames:          func(next action.Action) int { return skillHoldDelay + skillHoldFrames[next] },
		AnimationLength: skillHoldDelay + skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldDelay + skillHoldFrames[action.ActionBurst],
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
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	if t.StatusIsActive(skillHoldMarkKey) {
		t.DeleteStatus(skillHoldMarkKey)
		t.AddStatus(skillPressMarkKey, 360, true)
		c.Core.Tasks.Add(c.skillPressMark(t), 1.5*60)
		return
	}
	if t.StatusIsActive(skillPressMarkKey) {
		t.AddStatus(skillPressMarkKey, 360, true)
		return
	}
	if c.markCount < 5 {
		c.markCount++
		t.AddStatus(skillPressMarkKey, 360, true)
		c.Core.Tasks.Add(c.skillPressMark(t), 1.5*60)
		return
	}
}

func (c *char) skillHoldMarkTargets(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	if t.StatusIsActive(skillPressMarkKey) {
		t.DeleteStatus(skillPressMarkKey)
		t.AddStatus(skillHoldMarkKey, 720, true)
		c.Core.Tasks.Add(c.skillHoldMark(t), 1.5*60)
		return
	}
	if t.StatusIsActive(skillHoldMarkKey) {
		t.AddStatus(skillHoldMarkKey, 720, true)
		return
	}
	if c.markCount < 5 {
		c.markCount++
		t.AddStatus(skillHoldMarkKey, 720, true)
		c.Core.Tasks.Add(c.skillHoldMark(t), 1.5*60)
		return
	}
}

func (c *char) skillPressMark(t *enemy.Enemy) func() {
	return func() {
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
		c.Core.Tasks.Add(c.skillPressMark(t), 1.5*60)
	}
}

func (c *char) skillHoldMark(t *enemy.Enemy) func() {
	return func() {
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
		c.Core.Tasks.Add(c.skillHoldMark(t), 1.5*60)
	}
}

func (c *char) charlotteMarkOnTargetDied() {
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if !t.StatusIsActive(skillPressMarkKey) && !t.StatusIsActive(skillHoldMarkKey) {
			return false
		}
		c.markCount--
		return false
	}, "charlotte-on-target-died")
}
