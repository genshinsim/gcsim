package dahlia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames   [][]int
	skillHitmarks = []int{27, 28}
	skillCDStarts = []int{23, 21}
)

const (
	particleICDKey = "dahlia-particle-icd"
)

func init() {
	skillFrames = make([][]int, 2)

	// Tap E (tE)
	skillFrames[0] = frames.InitAbilSlice(53) // tE -> W
	skillFrames[0][action.ActionAttack] = 32  // tE -> N1
	skillFrames[0][action.ActionBurst] = 33   // tE -> Q
	skillFrames[0][action.ActionSkill] = 36   // tE -> tE / hE (TO-DO: check if this includes both)
	skillFrames[0][action.ActionDash] = 31    // tE -> D
	skillFrames[0][action.ActionJump] = 29    // tE -> J
	skillFrames[0][action.ActionSwap] = 31    // tE -> Swap

	// Hold E (hE)
	skillFrames[1] = frames.InitAbilSlice(55) // hE -> W
	skillFrames[1][action.ActionAttack] = 35  // hE -> N1
	skillFrames[1][action.ActionBurst] = 35   // hE -> Q
	skillFrames[1][action.ActionSkill] = 37   // hE -> tE / hE (TO-DO: check if this includes both)
	skillFrames[1][action.ActionDash] = 33    // hE -> D
	skillFrames[1][action.ActionJump] = 32    // hE -> J
	skillFrames[1][action.ActionSwap] = 33    // hE -> Swap
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	hold := 0
	if p["short_hold"] != 0 {
		hold = 1
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Immersive Ordinance",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), // TO-DO: Fix when data becomes available
		0, // TO-DO: Should this be skillHitmarks[hold] instead?
		skillHitmarks[hold],
		c.particleCB,
	)

	c.SetCDWithDelay(action.ActionSkill, 9*60, skillCDStarts[hold])

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[hold]),
		AnimationLength: skillFrames[hold][action.InvalidAction],
		CanQueueAfter:   skillFrames[hold][action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, true) // TO-DO: Fix when data becomes available

	c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Hydro, c.ParticleDelay)
}
