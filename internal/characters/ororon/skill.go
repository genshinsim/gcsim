package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const skillHitmark = 31

func init() {
	skillFrames = frames.InitAbilSlice(30)
	skillFrames[action.ActionDash] = 31
	skillFrames[action.ActionSwap] = 29
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.particlesGenerated = false
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Spirit Orb DMG",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagElementalArt,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           spiritOrb[c.TalentLvlSkill()],
	}

	enemies := []targets.TargetKey{c.Core.Combat.PrimaryTarget().Key()}
	maxHits := 3 + c.c1ExtraBounce()
	for i := 0; len(enemies) < maxHits && i < c.Core.Combat.EnemyCount(); i++ {
		newKey := c.Core.Combat.Enemies()[i].Key()
		if newKey == c.Core.Combat.PrimaryTarget().Key() {
			continue
		}
		enemies = append(enemies, newKey)
	}
	for i, e := range enemies {
		c.Core.QueueAttack(
			ai,
			combat.NewSingleTargetHit(e), // TODO: find out if this single target
			skillHitmark,
			skillHitmark+travel*(i+1),
			c.particleCB,
			c.makeA4cb(),
			c.makeC1cb(),
		)
	}

	c.SetCDWithDelay(action.ActionSkill, 15*60, 7)
	c.a1OnSkill()
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.particlesGenerated {
		return
	}
	c.particlesGenerated = true
	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Electro, c.ParticleDelay)
}
