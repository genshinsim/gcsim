package heizou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillEndFrames []int

func init() {
	skillEndFrames = frames.InitAbilSlice(19)
	skillEndFrames[action.ActionDash] = 10
	skillEndFrames[action.ActionJump] = 10
	skillEndFrames[action.ActionSwap] = 10
}

const (
	skillHitmark                 = 20
	skillCDStart                 = 18
	holdAtFullStacksPenalty      = 17 // if you hold while at 4 stacks it takes 17 extra frames to release
	skillHitlagHaltFrame         = 0.09
	skillHitlagMaxStackHaltFrame = 0.12
	particleICDKey               = "heizou-particle-icd"
)

func (c *char) skillHoldDuration(stacks int) int {
	// animation duration only
	// diff is the number of stacks we must charge up to reach the desired state
	diff := stacks - c.decStack
	if diff < 0 {
		diff = 0
	}
	if diff > 4 {
		diff = 4
	}
	// it's .75s per stack
	return 45 * diff
}

func (c *char) addDecStack() {
	if c.decStack < 4 {
		c.decStack++
		c.Core.Log.NewEvent("declension stack gained", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.decStack)
	}
}

func (c *char) skillRelease(delay int) action.Info {
	c.Core.Tasks.Add(func() {
		hitDelay := skillHitmark - skillCDStart
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Heartstopper Strike",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypeDefault,
			Element:            attributes.Anemo,
			Durability:         50,
			Mult:               skill[c.TalentLvlSkill()] + float64(c.decStack)*decBonus[c.TalentLvlSkill()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   skillHitlagHaltFrame * 60,
			CanBeDefenseHalted: false,
		}
		offset := -0.3
		if c.decStack == 0 {
			offset = -0.4
		}
		width := 3.0
		height := 5.0
		if c.decStack == 4 {
			ai.Abil = "Heartstopper Strike (Max Stacks)"
			ai.Mult += convicBonus[c.TalentLvlSkill()]
			ai.HitlagHaltFrames = skillHitlagMaxStackHaltFrame * 60
			width += 1
			height += 1
		}

		done := false
		skillCB := func(a combat.AttackCB) {
			c.decStack = 0
			if a.Target.Type() != targets.TargettableEnemy {
				return
			}
			if done {
				return
			}
			done = true

			c.a4()
		}

		snap := c.Snapshot(&ai)
		if c.Base.Cons >= 6 {
			c6CR, c6CD := c.c6()
			snap.Stats[attributes.CR] += c6CR
			snap.Stats[attributes.CD] += c6CD
		}
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: offset}, width, height),
			hitDelay,
			skillCB,
			c.particleCB,
		)
		c.SetCD(action.ActionSkill, 10*60)
	}, skillCDStart+delay)

	return action.Info{
		Frames:          func(next action.Action) int { return delay + skillEndFrames[next] + skillHitmark },
		AnimationLength: delay + skillEndFrames[action.InvalidAction] + skillHitmark,
		CanQueueAfter:   delay + skillEndFrames[action.ActionSwap] + skillHitmark, // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold() action.Info {
	if c.decStack == 4 {
		return c.skillRelease(holdAtFullStacksPenalty)
	}
	for i := c.decStack + 1; i <= 4; i++ {
		c.Core.Tasks.Add(c.addDecStack, c.skillHoldDuration(i))
	}
	return c.skillRelease(c.skillHoldDuration(4))
}

func (c *char) skillPress() action.Info {
	return c.skillRelease(0)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if p["hold"] != 0 {
		return c.skillHold(), nil
	}
	return c.skillPress(), nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, true)

	count := 2.0
	switch c.decStack {
	case 2, 3:
		if c.Core.Rand.Float64() < .5 {
			count = 3
		}
	case 4:
		count = 3
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Anemo, c.ParticleDelay)
}
