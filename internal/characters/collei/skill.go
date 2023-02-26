package collei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	skillKey     = "collei-skill"
	skillRelease = 20
	skillReturn  = 157
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
	// The game has ICD as AttackTagElementalArt, ICDTagElementalArt,
	// ICDGroupColleiBoomerangForward, and ICDGroupColleiBoomerangBack. However,
	// we believe this is unnecessary, so just use ICDTagNone.
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Floral Brush",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Dendro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	var c6Cb func(a combat.AttackCB)
	if c.Base.Cons >= 6 {
		c6Triggered := false
		c6Cb = func(a combat.AttackCB) {
			if c6Triggered {
				return
			}
			c6Triggered = true
			c.c6(a.Target)
		}
	}
	particleCB := c.makeParticleCB()
	//TODO: this should have its own position
	for _, hitmark := range skillHitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2),
			skillRelease,
			hitmark,
			c6Cb,
			particleCB,
		)
	}

	c.Core.Tasks.Add(func() {
		c.AddStatus(skillKey, skillReturn-skillRelease, false)
	}, skillRelease)

	c.sproutShouldExtend = false
	c.sproutShouldProc = c.Base.Cons >= 2 && c.Base.Ascension >= 1
	c.Core.Tasks.Add(func() {
		if !c.sproutShouldProc {
			return
		}
		src := c.Core.F
		c.sproutSrc = src
		duration := 180
		if c.sproutShouldExtend {
			duration += 180
		}
		c.AddStatus(sproutKey, duration, true)
		ai := c.a1AttackInfo()
		snap := c.Snapshot(&ai)
		c.QueueCharTask(func() {
			c.a1Ticks(src, snap)
		}, sproutHitmark)
	}, skillReturn)

	c.SetCDWithDelay(action.ActionSkill, 720, 20)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Dendro, c.ParticleDelay)
	}
}
