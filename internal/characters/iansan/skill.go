package iansan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillHitmark = 21
	skillFrames  []int

	fastSkill = "fast-skill"
)

func init() {
	skillFrames = frames.InitAbilSlice(43)
	skillFrames[action.ActionDash] = 31
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 42
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Thunderbolt Rush",
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		AttackTag:      attacks.AttackTagElementalArt,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeSlash,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			1,
		),
		skillHitmark,
		skillHitmark,
		c.particleCB,
	)

	c.AddStatus(fastSkill, 5*60, true)
	c.enterNightsoul(c.nightsoulState.MaxPoints)
	c.particleGenerated = false
	c.SetCD(action.ActionSkill, 16*60)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
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
	c.nightsoulState.EnterBlessing(points)
	c.nightsoulPointReduceTask(c.nightsoulSrc)
	c.setNightsoulExitTimer(16 * 60)
}

func (c *char) exitNightsoul() {
	c.nightsoulSrc = -1
	c.nightsoulState.ExitBlessing()
	c.DeleteStatus(burstStatus)
	c.DeleteStatus(a1Status)
}

func (c *char) nightsoulPointReduceTask(src int) {
	// reduce 0.6 point every 6f, which is 6 per second
	const tickInterval = .1

	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}

		points := 0.6
		c.nightsoulState.ConsumePoints(points)
		c.c1(points)
		c.updateATKBuff()
		if c.nightsoulState.Points() < 0.001 {
			c.exitNightsoul()
			return
		}

		c.nightsoulPointReduceTask(src)
	}, 60*tickInterval)
}

func (c *char) setNightsoulExitTimer(duration int) {
	src := c.nightsoulSrc
	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}
		c.nightsoulState.ClearPoints()
		c.exitNightsoul()
	}, duration)
}
