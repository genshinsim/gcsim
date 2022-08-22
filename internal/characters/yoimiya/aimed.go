package yoimiya

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var aimedFrames [][]int

var aimedHitmarks = []int{15 - 12, 15, 86, 103, 121, 139}

func init() {
	aimedFrames = make([][]int, 6)

	// Aimed Shot (ARCC)
	aimedFrames[0] = frames.InitAbilSlice(26 - 12)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(26)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot
	aimedFrames[2] = frames.InitAbilSlice(97)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]

	// Fully-Charged Aimed Shot (1 Kindling Arrow)
	aimedFrames[3] = frames.InitAbilSlice(114)
	aimedFrames[3][action.ActionDash] = aimedHitmarks[3]
	aimedFrames[3][action.ActionJump] = aimedHitmarks[3]

	// Fully-Charged Aimed Shot (2 Kindling Arrows)
	aimedFrames[4] = frames.InitAbilSlice(132)
	aimedFrames[4][action.ActionDash] = aimedHitmarks[4]
	aimedFrames[4][action.ActionJump] = aimedHitmarks[4]

	// Fully-Charged Aimed Shot (3 Kindling Arrows)
	aimedFrames[5] = frames.InitAbilSlice(150)
	aimedFrames[5][action.ActionDash] = aimedHitmarks[5]
	aimedFrames[5][action.ActionJump] = aimedHitmarks[5]
}

// Standard aimed attack
func (c *char) Aimed(p map[string]int) action.ActionInfo {
	hold, ok := p["hold"]
	if !ok || hold < 0 {
		hold = 2
	}
	if hold > 5 {
		hold = 2
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]
	if !ok || weakspot < 0 {
		weakspot = 0
	}
	if weakspot > 4 {
		weakspot = 4
	}

	// used to adjust how long it takes for the kindling arrows to hit starting from CA arrow release
	// does nothing if hold < 3
	kindling_travel, ok := p["kindling_travel"]
	if !ok {
		kindling_travel = 30
	}

	// Normal Arrow
	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone,
		ICDGroup:             combat.ICDGroupDefault,
		Element:              attributes.Pyro,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot >= 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	if hold < 2 {
		ai.Abil = "Aimed Shot"
		if hold == 0 {
			ai.Abil += " (ARCC)"
		}
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
	}

	c.Core.QueueAttack(
		ai,
		combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
		aimedHitmarks[hold],
		aimedHitmarks[hold]+travel,
	)

	if hold >= 3 {
		// Kindling Arrows
		ai.ICDTag = combat.ICDTagExtraAttack
		ai.Mult = aimExtra[c.TalentLvlAttack()]

		// no hitlag
		ai.HitlagHaltFrames = 0
		ai.HitlagFactor = 0.01
		ai.HitlagOnHeadshotOnly = false
		ai.IsDeployable = false

		for i := 3; i <= hold; i++ {
			aiCpy := ai // avoid closure issues
			aiCpy.Abil = fmt.Sprintf("Kindling Arrow %v", i-2)
			aiCpy.HitWeakPoint = weakspot >= i-1
			c.Core.Tasks.Add(func() {
				// E:
				// "During this time, Normal Attack: Firework Flare-Up will not generate Kindling Arrows at Charge Level 2.""
				// but if her E expires at the time of the arrow release, then the Kindling Arrows should be released:
				// https://youtu.be/XEAJ9ssLmv4
				if !c.StatusIsActive(skillKey) {
					c.Core.QueueAttack(
						aiCpy,
						combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy),
						0,
						kindling_travel, // add a bit of extra delay for kindling arrows
					)
				}
			}, aimedHitmarks[hold])
		}
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames[hold]),
		AnimationLength: aimedFrames[hold][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[hold],
		State:           action.AimState,
	}
}
