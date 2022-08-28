package collei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

const (
	skillKey    = "collei-skill"
	skillReturn = 157
)

var (
	skillHitmarks = []int{34, 138}
	skillFrames   []int
)

func init() {
	skillFrames = frames.InitAbilSlice(68)
	skillFrames[action.ActionAttack] = 65
	skillFrames[action.ActionAim] = 65
	skillFrames[action.ActionSkill] = 67
	skillFrames[action.ActionDash] = 54
	skillFrames[action.ActionJump] = 53
	skillFrames[action.ActionSwap] = 66
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Floral Brush",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snapDelay := skillHitmarks[0]
	for _, hitmark := range skillHitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy),
			snapDelay,
			hitmark,
		)
	}

	c.AddStatus(skillKey, skillReturn, false)

	c.sproutShouldExtend = false
	c.sproutShouldProc = c.Base.Cons >= 2
	c.Core.Tasks.Add(func() {
		if c.sproutShouldProc {
			c.sproutSrc = c.Core.F
			duration := 180
			if c.sproutShouldExtend {
				duration += 180
			}
			c.AddStatus(sproutKey, duration, true)
			c.QueueCharTask(func() { c.a1Ticks(c.sproutSrc) }, sproutHitmark)
		}
	}, skillReturn)

	// 50% chance of 3 orbs
	count := 2.0
	if c.Core.Rand.Float64() < .50 {
		count = 3.0
	}
	c.Core.QueueParticle("collei", count, attributes.Cryo, skillHitmarks[0]+c.ParticleDelay)

	c.SetCDWithDelay(action.ActionSkill, 720, 20)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}
