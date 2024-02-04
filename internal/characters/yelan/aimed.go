package yelan

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames [][]int
var aimedBarbFrames []int

var aimedHitmarks = []int{15, 86}

const aimedBarbHitmark = 32

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

	// Breakthrough Barb
	aimedBarbFrames = frames.InitAbilSlice(42)
	aimedBarbFrames[action.ActionDash] = aimedBarbHitmark
	aimedBarbFrames[action.ActionJump] = aimedBarbHitmark
}

// Aimed charge attack damage queue generator
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

	if c.breakthrough && hold == attacks.AimParamLv1 {
		c.breakthrough = false
		c.Core.Log.NewEvent("breakthrough state deleted", glog.LogCharacterEvent, c.Index)

		ai := combat.AttackInfo{
			ActorIndex:   c.Index,
			Abil:         "Breakthrough Barb",
			AttackTag:    attacks.AttackTagExtra,
			ICDTag:       attacks.ICDTagYelanBreakthrough,
			ICDGroup:     attacks.ICDGroupYelanBreakthrough,
			StrikeType:   attacks.StrikeTypePierce,
			Element:      attributes.Hydro,
			Durability:   25,
			FlatDmg:      barb[c.TalentLvlAttack()] * c.MaxHP(),
			HitWeakPoint: weakspot == 1,
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				6,
			),
			aimedBarbHitmark,
			aimedBarbHitmark+travel,
		)

		return action.Info{
			Frames:          frames.NewAbilFunc(aimedBarbFrames),
			AnimationLength: aimedBarbFrames[action.InvalidAction],
			CanQueueAfter:   aimedBarbHitmark,
			State:           action.AimState,
		}, nil
	}

	ai := combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Fully-Charged Aimed Shot",
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypePierce,
		Element:      attributes.Hydro,
		Durability:   25,
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
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(aimedFrames[hold]),
		AnimationLength: aimedFrames[hold][action.InvalidAction],
		CanQueueAfter:   aimedHitmarks[hold],
		State:           action.AimState,
	}, nil
}
