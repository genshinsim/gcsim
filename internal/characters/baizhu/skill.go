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

const particleICDKey = "baizhu-particle-icd"

func init() {
	skillFrames = frames.InitAbilSlice(51)
	skillFrames[action.ActionAttack] = 51
	skillFrames[action.ActionCharge] = 51
	skillFrames[action.ActionSkill] = 51
	skillFrames[action.ActionBurst] = 51
	skillFrames[action.ActionDash] = 51
	skillFrames[action.ActionJump] = 51
	skillFrames[action.ActionSwap] = 51

}

const (
	skillFirstHitmark = 40 //TODO:Freims
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
	c.c6done = false

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

	//trigger a chain of attacks starting at the first target
	atk := *c.skillAtk
	atk.SourceFrame = c.Core.F
	atk.Pattern = combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key())
	cb := c.chain(c.Core.F, 1)
	if cb != nil {
		atk.Callbacks = append(atk.Callbacks, c.particleCB, cb)
		if c.Base.Cons >= 6 {
			atk.Callbacks = append(atk.Callbacks, c.c6)
		}
	}
	c.Core.QueueAttackEvent(&atk, skillFirstHitmark)

	c.SetCDWithDelay(action.ActionSkill, 10*60, 1) //TODO:Delay on CD?

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) chain(src int, count int) combat.AttackCBFunc {
	if count == 3 {
		c.skillHealing()
		return nil
	}
	return func(a combat.AttackCB) {
		//on hit figure out the next target
		next := c.Core.Combat.RandomEnemyWithinArea(combat.NewCircleHitOnTarget(a.Target, nil, 8), func(t combat.Enemy) bool {
			return a.Target.Key() != t.Key()
		})
		if next == nil {
			c.skillHealing()
			return
		}
		//queue an attack vs next target
		atk := *c.skillAtk
		atk.SourceFrame = src
		atk.Pattern = combat.NewSingleTargetHit(next.Key())
		cb := c.chain(src, count+1)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.QueueAttackEvent(&atk, 10) //TODO: Modify delay

	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Dendro, c.ParticleDelay)
}

func (c *char) skillHealing() {
	c.Core.Tasks.Add(func() {
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Universal Diagnosis Healing",
			Src:     skillHealPP[c.TalentLvlBurst()] * c.MaxHP(),
			Bonus:   skillHealFlat[c.TalentLvlBurst()],
		})

	}, 22) //TODO: change delay

}
