package amber

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames [][]int

var aimedHitmarks = []int{15 - 12, 15, 86}
var c1Hitmarks = []int{24 - 12, 24, 95} // C1 arrow comes out 9f after the normal one, still comes out even if you cancel at aimedHitmark

func init() {
	aimedFrames = make([][]int, 3)

	// Aimed Shot (ARCC)
	aimedFrames[0] = frames.InitAbilSlice(25 - 12)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(25)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot
	aimedFrames[2] = frames.InitAbilSlice(96)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	hold, ok := p["hold"]
	if !ok || hold < 0 {
		hold = 2
	}
	if hold > 2 {
		hold = 2
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	b := p["bunny"]

	if c.Base.Cons >= 2 && b != 0 {
		//explode the first bunny
		c.Core.Tasks.Add(func() {
			c.manualExplode()
		}, aimedHitmarks[hold]+travel)

		//also don't do any dmg since we're shooting at bunny
		return action.ActionInfo{
			Frames:          frames.NewAbilFunc(aimedFrames[hold]),
			AnimationLength: aimedFrames[hold][action.InvalidAction],
			CanQueueAfter:   aimedHitmarks[hold],
			State:           action.AimState,
		}
	}

	ai := combat.AttackInfo{
		Abil:         "Fully-Charged Aimed Shot",
		ActorIndex:   c.Index,
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagExtraAttack,
		ICDGroup:     combat.ICDGroupAmber,
		Element:      attributes.Pyro,
		Durability:   50,
		Mult:         fullaim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	if hold < 2 {
		ai.Abil = "Aimed Shot"
		if hold == 0 {
			ai.Abil += " (ARCC)"
		}
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
	}
	c.Core.QueueAttack(ai, combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy), aimedHitmarks[hold], aimedHitmarks[hold]+travel, c.a4)

	if c.Base.Cons >= 1 {
		ai.Abil += " (C1)"
		ai.Mult = .2 * ai.Mult
		c.Core.QueueAttack(ai, combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy), c1Hitmarks[hold], c1Hitmarks[hold]+travel)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames[hold]),
		AnimationLength: aimedFrames[hold][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[hold],
		State:           action.AimState,
	}
}
