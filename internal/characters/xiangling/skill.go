package xiangling

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	infuseWindow     = 30
	infuseDurability = 20
	particleICDKey   = "xiangling-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(39)
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionJump] = 14
	skillFrames[action.ActionSwap] = 38
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Guoba",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       guobaTick[c.TalentLvlSkill()],
	}

	// delay in frames from guoba expiry until the a4 chili pepper is picked up
	a4Delay, ok := p["a4_delay"]
	if !ok {
		a4Delay = 0
	}
	if a4Delay < 0 {
		a4Delay = 0
	}
	if a4Delay > 10*60 {
		a4Delay = 10 * 60
	}

	// guoba spawns at cd frame
	// lasts 7.3 seconds, shoots every 100 frames
	c.Core.Tasks.Add(func() {
		guoba := c.newGuoba(ai)
		c.AddStatus("xianglingguoba", guoba.Duration, false)
		c.Core.Combat.AddGadget(guoba)
		// queue up a4 relative to guoba expiry
		c.a4(guoba.Duration + a4Delay)
	}, 13)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 13)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1*60, false)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Pyro, c.ParticleDelay)
}
