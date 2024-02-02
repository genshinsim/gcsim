package yoimiya

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var aimedFrames [][]int

var aimedHitmarks = []int{86, 103, 121, 139}

func init() {
	aimedFrames = make([][]int, 4)

	// Normal CA
	aimedFrames[0] = frames.InitAbilSlice(97)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// CA - 1 Kindling Arrow
	aimedFrames[1] = frames.InitAbilSlice(114)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// CA - 2 Kindling Arrows
	aimedFrames[2] = frames.InitAbilSlice(132)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]

	// CA - 3 Kindling Arrows
	aimedFrames[3] = frames.InitAbilSlice(150)
	aimedFrames[3][action.ActionDash] = aimedHitmarks[3]
	aimedFrames[3][action.ActionJump] = aimedHitmarks[3]
}

// Standard aimed attack
func (c *char) Aimed(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	// determines what CA gets executed
	kindling, ok := p["kindling"]
	if !ok || kindling < 0 {
		kindling = 0
	}
	if kindling > 3 {
		kindling = 3
	}

	// used to adjust how long it takes for the kindling arrows to hit starting from CA arrow release
	kindlingTravel, ok := p["kindling_travel"]
	if !ok {
		kindlingTravel = 30
	}

	// Normal Arrow
	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Aimed Shot",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Pyro,
		Durability:           25,
		Mult:                 aim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	c2CB := c.makeC2CB()
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -0.5},
			0.1,
			1,
		),
		aimedHitmarks[kindling],
		aimedHitmarks[kindling]+travel,
		c2CB,
	)

	// Kindling Arrows
	if kindling > 0 {
		ai.ICDTag = attacks.ICDTagExtraAttack
		ai.Mult = aimExtra[c.TalentLvlAttack()]

		// TODO:
		// Kindling Arrows can hit weakspots to proc stuff like Prototype Crescent, but they don't always crit
		// current assumption is that they never hit a weakspot
		ai.HitWeakPoint = false

		// no hitlag
		ai.HitlagHaltFrames = 0
		ai.HitlagFactor = 0.01
		ai.HitlagOnHeadshotOnly = false
		ai.IsDeployable = false

		for i := 1; i <= kindling; i++ {
			ai.Abil = fmt.Sprintf("Kindling Arrow %v", i)
			// add a bit of extra delay for kindling arrows
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					0.6,
				),
				aimedHitmarks[kindling],
				aimedHitmarks[kindling]+kindlingTravel,
				c2CB,
			)
		}
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(aimedFrames[kindling]),
		AnimationLength: aimedFrames[kindling][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[kindling],
		State:           action.AimState,
	}, nil
}
