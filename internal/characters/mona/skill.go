package mona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames   [][]int
	skillHitmarks = []int{86, 101}
)

const particleICDKey = "mona-particle-icd"

func init() {
	skillFrames = make([][]int, 2)
	// Tap E
	skillFrames[0] = frames.InitAbilSlice(50) // Tap E -> N1
	skillFrames[0][action.ActionCharge] = 46  // Tap E -> CA
	skillFrames[0][action.ActionBurst] = 28   // Tap E -> Q
	skillFrames[0][action.ActionDash] = 36    // Tap E -> D
	skillFrames[0][action.ActionJump] = 28    // Tap E -> J
	skillFrames[0][action.ActionSwap] = 43    // Tap E -> Swap

	// Hold E
	skillFrames[1] = frames.InitAbilSlice(80) // Hold E -> N1
	skillFrames[1][action.ActionCharge] = 76  // Hold E -> CA
	skillFrames[1][action.ActionBurst] = 58   // Hold E -> Q
	skillFrames[1][action.ActionDash] = 66    // Hold E -> D
	skillFrames[1][action.ActionJump] = 59    // Hold E -> J
	skillFrames[1][action.ActionSwap] = 73    // Hold E -> Swap
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	hold := 0
	if p["hold"] != 0 {
		hold = 1
	}

	// DoT
	// ticks 4 times
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Mirror Reflection of Doom (Tick)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skillDot[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)

	// tick every 1s
	for i := skillHitmarks[hold]; i < 300; i += 60 {
		c.Core.QueueAttackWithSnap(ai, snap, ap, i)
	}

	// Explosion
	aiExplode := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Mirror Reflection of Doom (Explode)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(aiExplode, ap, 0, skillHitmarks[hold]+313, c.particleCB)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 24)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[hold]),
		AnimationLength: skillFrames[hold][action.InvalidAction],
		CanQueueAfter:   skillFrames[hold][action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1*60, false)

	count := 3.0
	if c.Core.Rand.Float64() < .33 {
		count = 4
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Hydro, c.ParticleDelay)
}
