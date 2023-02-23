package kaeya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillFrames []int

const (
	skillHitmark   = 28
	particleICDKey = "kaeya-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(53) // E -> N1
	skillFrames[action.ActionBurst] = 52   // E -> Q
	skillFrames[action.ActionDash] = 25    // E -> D
	skillFrames[action.ActionJump] = 26    // E -> J
	skillFrames[action.ActionSwap] = 49    // E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Frostgnaw",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	cb := func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}

		// A1:
		// Every hit with Frostgnaw regenerates HP for Kaeya equal to 15% of his ATK.
		if c.Base.Ascension >= 1 {
			heal := .15 * (a.AttackEvent.Snapshot.BaseAtk*(1+a.AttackEvent.Snapshot.Stats[attributes.ATKP]) + a.AttackEvent.Snapshot.Stats[attributes.ATK])
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "Cold-Blooded Strike",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			})
		}
	}
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -0.2}, 4, 8),
		0,
		skillHitmark,
		cb,
		c.particleCB,
		c.makeA4ParticleCB(),
	)

	c.SetCDWithDelay(action.ActionSkill, 360, 25)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, true)

	count := 2.0
	if c.Core.Rand.Float64() < 0.67 {
		count = 3
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Cryo, c.ParticleDelay)
}

// Opponents Frozen by Frostgnaw will drop additional Elemental Particles.
// Frostgnaw may only produce a maximum of 2 additional Elemental Particles per use.
func (c *char) makeA4ParticleCB() combat.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	a4Count := 0
	return func(a combat.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		if a4Count == 2 {
			return
		}
		if !e.AuraContains(attributes.Frozen) {
			return
		}
		c.Core.Log.NewEvent("kaeya a4 proc", glog.LogCharacterEvent, c.Index)
		a4Count++
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Cryo, c.ParticleDelay)
	}
}
