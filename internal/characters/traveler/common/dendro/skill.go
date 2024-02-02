package dendro

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames [][]int

const (
	skillHitmark   = 28
	particleICDKey = "travelerdendro-particle-icd"
	cdStart        = 25
)

func init() {
	skillFrames = make([][]int, 2)

	// Male
	skillFrames[0] = frames.InitAbilSlice(37) // E -> N1
	skillFrames[0][action.ActionDash] = 29    // E -> D
	skillFrames[0][action.ActionJump] = 29    // E -> J
	skillFrames[0][action.ActionSwap] = 36    // E -> Swap

	// Female
	skillFrames[1] = frames.InitAbilSlice(37) // E -> N1/Q
	skillFrames[1][action.ActionDash] = 28    // E -> D
	skillFrames[1][action.ActionJump] = 28    // E -> J
	skillFrames[1][action.ActionSwap] = 35    // E -> Swap
}

func (c *Traveler) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Razorgrass Blade",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	var skillCB func(a combat.AttackCB)
	if c.Base.Cons >= 1 {
		c.skillC1 = true
		skillCB = c.c1cb()
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), geometry.Point{Y: -0.3}, 6.5, 130),
		skillHitmark,
		skillHitmark,
		skillCB,
		c.particleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 8*60, cdStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[c.gender]),
		AnimationLength: skillFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *Traveler) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, true)

	count := 2.0
	if c.Core.Rand.Float64() < 0.5 {
		count = 3
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Dendro, c.ParticleDelay)
}
