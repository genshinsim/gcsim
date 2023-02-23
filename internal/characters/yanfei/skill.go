package yanfei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillFrames []int

const (
	skillHitmark   = 32
	particleICDKey = "yanfei-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(46) // E -> N1
	skillFrames[action.ActionCharge] = 35  // E -> CA
	skillFrames[action.ActionBurst] = 43   // E -> Q
	skillFrames[action.ActionDash] = 29    // E -> D
	skillFrames[action.ActionJump] = 34    // E -> J
	skillFrames[action.ActionSwap] = 44    // E -> Swap
}

// Yanfei skill - Straightforward as it has little interactions with the rest of her kit
// Summons flames that deal AoE Pyro DMG. Opponents hit by the flames will grant Yanfei the maximum number of Scarlet Seals.
func (c *char) Skill(p map[string]int) action.ActionInfo {
	done := false
	addSeal := func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		// Create max seals on hit
		if c.sealCount < c.maxTags {
			c.sealCount = c.maxTags
		}
		c.AddStatus(sealBuffKey, 600, true)
		c.Core.Log.NewEvent("yanfei gained max seals", glog.LogCharacterEvent, c.Index).
			Write("current_seals", c.sealCount)
		done = true
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Signed Edict",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	// TODO: Not sure of snapshot timing
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			3.5,
		),
		0,
		skillHitmark,
		c.particleCB,
		addSeal,
	)

	c.SetCDWithDelay(action.ActionSkill, 540, 28)

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
	c.AddStatus(particleICDKey, 0.2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Pyro, c.ParticleDelay)
}
