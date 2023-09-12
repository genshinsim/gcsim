package baizhu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(49) // E -> N1
	skillFrames[action.ActionCharge] = 48
	skillFrames[action.ActionSkill] = 40
	skillFrames[action.ActionBurst] = 30
	skillFrames[action.ActionDash] = 29
	skillFrames[action.ActionJump] = 29
	skillFrames[action.ActionWalk] = 47
	skillFrames[action.ActionSwap] = 28
}

const (
	skillFirstHitmark = 13
	skillTickInterval = 48
	skillReturnTravel = 51
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Universal Diagnosis",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillDamage[c.TalentLvlSkill()],
	}

	snap := c.Snapshot(&ai)
	c.skillAtk = &combat.AttackEvent{
		Info:     ai,
		Snapshot: snap,
	}

	// trigger a chain of attacks starting at the first target
	atk := *c.skillAtk
	atk.SourceFrame = c.Core.F
	atk.Pattern = combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.6)
	cb := c.chain(c.Core.F, 1)
	if cb != nil {
		atk.Callbacks = append(atk.Callbacks, c.makeParticleCB(), cb)
		if c.Base.Cons >= 6 {
			atk.Callbacks = append(atk.Callbacks, c.makeC6CB())
		}
	}
	c.Core.QueueAttackEvent(&atk, skillFirstHitmark)

	c.SetCDWithDelay(action.ActionSkill, 10*60, 23)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) chain(src, count int) combat.AttackCBFunc {
	if count == 3 {
		c.skillHealing()
		return nil
	}
	return func(a combat.AttackCB) {
		// on hit figure out the next target
		next := c.Core.Combat.RandomEnemyWithinArea(combat.NewCircleHitOnTarget(a.Target, nil, 10), nil)
		if next == nil {
			c.skillHealing()
			return
		}
		delay := skillTickInterval
		if next.Key() != a.Target.Key() {
			delay += 6 // add some (estimated) delay in case it's a different target
		}
		// queue an attack vs next target
		atk := *c.skillAtk
		atk.SourceFrame = src
		atk.Pattern = combat.NewCircleHitOnTarget(next, nil, 0.6)
		cb := c.chain(src, count+1)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.QueueAttackEvent(&atk, delay)
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

		count := 3.0
		if c.Core.Rand.Float64() < 0.50 {
			count = 4
		}
		c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Dendro, c.ParticleDelay)
	}
}

func (c *char) skillHealing() {
	c.Core.Tasks.Add(func() {
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Universal Diagnosis Healing",
			Src:     skillHealPP[c.TalentLvlBurst()]*c.MaxHP() + skillHealFlat[c.TalentLvlBurst()],
			Bonus:   c.Stat(attributes.Heal),
		})
	}, skillReturnTravel)
}
