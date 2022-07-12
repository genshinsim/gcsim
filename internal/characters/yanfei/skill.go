package yanfei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillFrames []int

const skillHitmark = 46

func init() {
	skillFrames = frames.InitAbilSlice(46)
}

// Yanfei skill - Straightforward as it has little interactions with the rest of her kit
// Summons flames that deal AoE Pyro DMG. Opponents hit by the flames will grant Yanfei the maximum number of Scarlet Seals.
func (c *char) Skill(p map[string]int) action.ActionInfo {
	done := false
	addSeal := func(a combat.AttackCB) {
		if done {
			return
		}
		// Create max seals on hit
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"] = c.maxTags
		}
		c.sealExpiry = c.Core.F + 600
		c.Core.Log.NewEvent("yanfei gained max seals", glog.LogCharacterEvent, c.Index).
			Write("current_seals", c.Tags["seal"]).
			Write("expiry", c.sealExpiry)
		done = true
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Signed Edict",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	// TODO: Not sure of snapshot timing
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, skillHitmark, addSeal)

	c.Core.QueueParticle("yanfei", 3, attributes.Pyro, skillHitmark+c.Core.Flags.ParticleDelay)

	c.SetCD(action.ActionSkill, 540)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,
		State:           action.SkillState,
	}
}
