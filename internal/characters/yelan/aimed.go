package yelan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames []int
var aimedBarbFrames []int

const aimedHitmark = 86
const aimedBarbHitmark = 32

func init() {
	// TODO: confirm that aim->x is the same for all cancels
	aimedFrames = frames.InitAbilSlice(96)
	aimedFrames[action.ActionDash] = aimedHitmark
	aimedFrames[action.ActionJump] = aimedHitmark

	aimedBarbFrames = frames.InitAbilSlice(42)
	aimedBarbFrames[action.ActionDash] = aimedBarbHitmark
	aimedBarbFrames[action.ActionJump] = aimedBarbHitmark
}

// Aimed charge attack damage queue generator
func (c *char) Aimed(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	if c.breakthrough {
		c.breakthrough = false
		c.Core.Log.NewEvent("breakthrough state deleted", glog.LogCharacterEvent, c.Index)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Breakthrough Barb",
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagYelanBreakthrough,
			ICDGroup:   attacks.ICDGroupYelanBreakthrough,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    barb[c.TalentLvlAttack()] * c.MaxHP(),
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
		Abil:         "Aim Charge Attack",
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       attacks.ICDTagNone,
		ICDGroup:     attacks.ICDGroupDefault,
		StrikeType:   attacks.StrikeTypePierce,
		Element:      attributes.Hydro,
		Durability:   25,
		Mult:         aimed[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
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
		aimedHitmark,
		aimedHitmark+travel,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(aimedFrames),
		AnimationLength: aimedFrames[action.InvalidAction],
		CanQueueAfter:   aimedHitmark,
		State:           action.AimState,
	}, nil
}
