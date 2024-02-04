package amber

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

var aimedHitmarks = []int{15, 86}
var c1Delay = 9 // C1 arrow comes out 9f after the normal one, still comes out even if you cancel at aimedHitmark

func init() {
	aimedFrames = make([][]int, 2)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(25)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(96)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]
}

func (c *char) Aimed(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = attacks.AimParamLv1
	}
	switch hold {
	case attacks.AimParamPhys:
	case attacks.AimParamLv1:
	default:
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	b := p["bunny"]

	// only works for fully-charged aimed shot
	if c.Base.Cons >= 2 && b != 0 && hold == attacks.AimParamLv1 {
		// explode the first bunny
		c.Core.Tasks.Add(func() {
			c.manualExplode()
		}, aimedHitmarks[hold]+travel)

		// also don't do any dmg since we're shooting at bunny
		return action.Info{
			Frames:          frames.NewAbilFunc(aimedFrames[hold]),
			AnimationLength: aimedFrames[hold][action.InvalidAction],
			CanQueueAfter:   aimedHitmarks[hold],
			State:           action.AimState,
		}, nil
	}

	ai := combat.AttackInfo{
		Abil:         "Fully-Charged Aimed Shot",
		ActorIndex:   c.Index,
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       attacks.ICDTagExtraAttack,
		ICDGroup:     attacks.ICDGroupAmber,
		StrikeType:   attacks.StrikeTypePierce,
		Element:      attributes.Pyro,
		Durability:   50,
		Mult:         fullaim[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	if hold < attacks.AimParamLv1 {
		ai.Abil = "Aimed Shot"
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
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
		c.makeA4CB(),
	)

	if c.Base.Cons >= 1 {
		ai.Mult = .2 * ai.Mult
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			aimedHitmarks[hold]+c1Delay,
			aimedHitmarks[hold]+c1Delay+travel,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(aimedFrames[hold]),
		AnimationLength: aimedFrames[hold][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[hold],
		State:           action.AimState,
	}, nil
}
