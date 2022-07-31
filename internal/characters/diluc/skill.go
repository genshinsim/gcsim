package diluc

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames [][]int
var skillHitmarks = []int{24, 28, 46}
var skillHitlagStages = []float64{.12, .12, .16}

func init() {
	skillFrames = make([][]int, 3)

	// skill (1st) -> x
	skillFrames[0] = frames.InitAbilSlice(32)
	skillFrames[0][action.ActionSkill] = 31
	skillFrames[0][action.ActionDash] = 24
	skillFrames[0][action.ActionJump] = 24
	skillFrames[0][action.ActionSwap] = 30

	// skill (2nd) -> x
	skillFrames[1] = frames.InitAbilSlice(38)
	skillFrames[1][action.ActionSkill] = 37
	skillFrames[1][action.ActionBurst] = 37
	skillFrames[1][action.ActionDash] = 28
	skillFrames[1][action.ActionJump] = 31
	skillFrames[1][action.ActionSwap] = 36

	// skill (3rd) -> x
	// TODO: missing counts for skill -> skill
	skillFrames[2] = frames.InitAbilSlice(66)
	skillFrames[2][action.ActionAttack] = 58
	skillFrames[2][action.ActionSkill] = 57 // uses burst frames
	skillFrames[2][action.ActionBurst] = 57
	skillFrames[2][action.ActionDash] = 47
	skillFrames[2][action.ActionJump] = 48
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// reset counter
	if c.Core.F >= c.eWindow {
		c.eCounter = 0
	}

	hitmark := skillHitmarks[c.eCounter]

	//actual skill cd starts immediately on first cast
	//times out after 4 seconds of not using
	//every hit applies pyro
	//apply attack speed
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("Searing Onslaught %v", c.eCounter),
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skill[c.eCounter][c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   skillHitlagStages[c.eCounter] * 60,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy), hitmark, hitmark)

	var orb float64 = 1
	if c.Core.Rand.Float64() < 0.33 {
		orb = 2
	}
	c.Core.QueueParticle("diluc", orb, attributes.Pyro, hitmark+c.Core.Flags.ParticleDelay)

	//add a timer to activate c4
	if c.Base.Cons >= 4 {
		c.Core.Tasks.Add(func() {
			c.c4()
		}, hitmark+120) // 2seconds after cast
	}

	// allow skill to be used again if 4s hasn't passed since last use
	c.eWindow = c.Core.F + 60*4

	// store skill counter so we can determine which frames to return
	idx := c.eCounter
	c.eCounter++
	switch c.eCounter {
	case 1:
		// TODO: cd delay?
		// set cd on first use
		c.SetCD(action.ActionSkill, 10*60)
	case 3:
		// reset window since we're at 3rd use
		c.eWindow = -1
		c.eCounter = 0
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[idx]),
		AnimationLength: skillFrames[idx][action.InvalidAction],
		CanQueueAfter:   skillFrames[idx][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
