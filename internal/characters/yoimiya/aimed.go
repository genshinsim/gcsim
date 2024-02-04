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

var aimedHitmarks = []int{15, 86, 103, 121, 139}

func init() {
	aimedFrames = make([][]int, 5)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(26)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(97)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot (1 Kindling Arrow)
	aimedFrames[2] = frames.InitAbilSlice(114)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]

	// Fully-Charged Aimed Shot (2 Kindling Arrows)
	aimedFrames[3] = frames.InitAbilSlice(132)
	aimedFrames[3][action.ActionDash] = aimedHitmarks[3]
	aimedFrames[3][action.ActionJump] = aimedHitmarks[3]

	// Fully-Charged Aimed Shot (3 Kindling Arrows)
	aimedFrames[4] = frames.InitAbilSlice(150)
	aimedFrames[4][action.ActionDash] = aimedHitmarks[4]
	aimedFrames[4][action.ActionJump] = aimedHitmarks[4]
}

// Standard aimed attack
func (c *char) Aimed(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	hold, ok := p["hold"]
	if !ok {
		hold = attacks.AimParamLv1
	}
	switch hold {
	case attacks.AimParamPhys:
	case attacks.AimParamLv1:
	case attacks.AimParamLv2:
	case attacks.AimParamLv3:
	case attacks.AimParamLv4:
	default:
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}

	// used to adjust how long it takes for the kindling arrows to hit starting from CA arrow release
	// does nothing if hold < lv. 2
	kindlingTravel, ok := p["kindling_travel"]
	if !ok {
		kindlingTravel = 30
	}

	// Normal Arrow
	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Pyro,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	c2CB := c.makeC2CB()
	if hold < attacks.AimParamLv1 {
		ai.Abil = "Aimed Shot"
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
		c2CB = nil
	}
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -0.5},
			0.1,
			1,
		),
		aimedHitmarks[hold],
		aimedHitmarks[hold]+travel,
		c2CB,
	)

	// Kindling Arrows
	if hold >= attacks.AimParamLv2 {
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

		for i := 1; i <= hold-1; i++ {
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
				aimedHitmarks[hold],
				aimedHitmarks[hold]+kindlingTravel,
				c2CB,
			)
		}
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(aimedFrames[hold]),
		AnimationLength: aimedFrames[hold][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[hold],
		State:           action.AimState,
	}, nil
}
