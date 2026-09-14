package kachina

import (
	"errors"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames       []int
	skillRecastFrames []int
)

func init() {
	skillFrames = frames.InitAbilSlice(733)
	skillFrames[action.ActionAttack] = 56
	skillFrames[action.ActionBurst] = 30
	skillFrames[action.ActionSkill] = 51
	skillFrames[action.ActionDash] = 55
	skillFrames[action.ActionJump] = 57
	skillFrames[action.ActionWalk] = 45
	skillFrames[action.ActionSwap] = 44

	skillRecastFrames = frames.InitAbilSlice(30)
	skillRecastFrames[action.ActionAttack] = 30
	skillRecastFrames[action.ActionBurst] = 30
	skillRecastFrames[action.ActionDash] = 30
	skillRecastFrames[action.ActionJump] = 30
	skillRecastFrames[action.ActionWalk] = 30
	skillRecastFrames[action.ActionSwap] = 30
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if p["recast"] > 0 {
		if p["hold"] > 0 {
			return action.Info{}, errors.New("cannot hold E while recasting Turbo Twirly")
		}
		if !c.StatusIsActive(twirlyKey) {
			return action.Info{}, errors.New("cannot recast E while Turbo Twirly is inactive")
		}
		return c.recastTwirly(), nil
	}

	if c.StatusIsActive(twirlyKey) {
		return c.recastTwirlyWithHold(p["hold"] > 0), nil
	}

	c.startTwirly(p["hold"] > 0)
	c.SetCDWithDelay(action.ActionSkill, int(skill*60), 10)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst],
		State:           action.SkillState,
	}, nil
}

func (c *char) recastTwirly() action.Info {
	return c.recastTwirlyWithHold(false)
}

func (c *char) recastTwirlyWithHold(hold bool) action.Info {
	if hold {
		c.mounted = true
	} else {
		c.mounted = !c.mounted
	}
	if !c.mounted {
		c.queueIndependent(c.twirlySrc, 30)
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionAttack],
		State:           action.SkillState,
	}
}

func (c *char) queueIndependent(src, delay int) {
	c.QueueCharTask(func() {
		if c.twirlySrc != src || !c.StatusIsActive(twirlyKey) || c.mounted {
			return
		}
		ai := c.twirlyIndependentAttackInfo("Turbo Twirly", independent[c.TalentLvlSkill()])
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, c.twirlyRadius()),
			0,
			0,
			c.twirlyParticleCB,
			func(info.AttackCB) { c.consumeTwirlyPoints(10) },
		)
		c.QueueCharTask(func() { c.queueIndependent(src, 118) }, 118)
	}, delay)
}

func (c *char) twirlyParticleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy || c.StatusIsActive("kachina-particle-icd") {
		return
	}
	c.AddStatus("kachina-particle-icd", 12, true)
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Geo, c.ParticleDelay)
}
