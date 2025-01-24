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
	skillFrames = frames.InitAbilSlice(31) // E -> Dash
	skillFrames[action.ActionAttack] = 30
	skillFrames[action.ActionCharge] = 30
	skillFrames[action.ActionBurst] = 30
	skillFrames[action.ActionJump] = 30
	skillFrames[action.ActionWalk] = 30
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
		HitlagFactor:   0.05,
	}

	enemies := []combat.Target{c.Core.Combat.PrimaryTarget()}
	maxHits := 3 + c.c1ExtraBounce()
	for i := 0; len(enemies) < maxHits && i < c.Core.Combat.EnemyCount(); i++ {
		newEnemy := c.Core.Combat.Enemy(i)
		if newEnemy.Key() == enemies[0].Key() {
			continue
		}
		enemies = append(enemies, newEnemy)
	}
	for i, e := range enemies {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(e, nil, 0.6),
			skillHitmark,
			skillHitmark+travel*(i+1),
			c.particleCB,
			c.makeA4cb(),
			c.makeC1cb(),
		)
	}

	c.SetCDWithDelay(action.ActionSkill, 15*60, 14)
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
