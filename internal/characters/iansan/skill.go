package iansan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillHitmark = 9
	skillFrames  []int

	fastSkill = "fast-skill"
)

func init() {
	skillFrames = frames.InitAbilSlice(137) // E -> E
	skillFrames[action.ActionAttack] = 32
	skillFrames[action.ActionBurst] = 31
	skillFrames[action.ActionDash] = 27
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionWalk] = 38
	skillFrames[action.ActionSwap] = 33
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Thunderbolt Rush",
		AdditionalTags:     []attacks.AttackTag{attacks.AttackTagNightsoul},
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Electro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.8,
		),
		skillHitmark,
		skillHitmark,
		c.particleCB,
	)

	c.AddStatus(fastSkill, 5*60, true)
	c.Core.Tasks.Add(func() {
		c.enterNightsoul(c.nightsoulState.MaxPoints)
	}, 3)
	c.particleGenerated = false
	c.SetCD(action.ActionSkill, 16*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.particleGenerated {
		return
	}
	c.particleGenerated = true

	count := 4.0
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Electro, c.ParticleDelay)
}

func (c *char) enterNightsoul(points float64) {
	c.nightsoulSrc = c.Core.F
	c.nightsoulState.EnterTimedBlessing(points, 16*60, c.exitNightsoul)
	c.nightsoulPointReduceTask(c.nightsoulSrc)
}

func (c *char) exitNightsoul() {
	c.nightsoulSrc = -1
	c.burstSrc = -1
	c.nightsoulState.ExitBlessing()
	c.nightsoulState.ClearPoints()
	c.DeleteStatus(burstStatus)
	c.DeleteStatus(a1Key)
}

func (c *char) reduceNightsoulPoints(points float64) {
	c.nightsoulState.ConsumePoints(points)
	c.c1(points)
	if c.nightsoulState.Points() <= 0.2 {
		c.exitNightsoul()
	}
}

func (c *char) nightsoulPointReduceTask(src int) {
	// reduce 0.6 point every 6f, which is 6 per second
	const tickInterval = .1

	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}

		c.reduceNightsoulPoints(0.6)
		if c.nightsoulState.Points() <= 0.2 {
			return
		}

		c.nightsoulPointReduceTask(src)
	}, 60*tickInterval)
}
