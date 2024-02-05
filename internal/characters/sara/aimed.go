package sara

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
var aimedA1Frames []int

var aimedHitmarks = []int{15, 86}

const aimedA1Hitmark = 50

func init() {
	// outside of E status
	aimedFrames = make([][]int, 2)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(25)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(96)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot (Crowfeather)
	aimedA1Frames = frames.InitAbilSlice(60)
	aimedA1Frames[action.ActionDash] = aimedA1Hitmark
	aimedA1Frames[action.ActionJump] = aimedA1Hitmark
}

// Aimed charge attack damage queue generator
// Additionally handles crowfeather state, E skill damage, and A4
// Has two parameters, "travel", used to set the number of frames that the arrow is in the air (default = 10)
// weak_point, used to determine if an arrow is hitting a weak point (default = 1 for true)
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

	// A1:
	// While in the Crowfeather Cover state provided by Tengu Stormcall, Aimed Shot charge times are decreased by 60%.
	skillActive := c.Base.Ascension >= 1 && c.Core.Status.Duration(coverKey) > 0

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Electro,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     .12 * 60,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	if hold < attacks.AimParamLv1 {
		ai.Abil = "Aimed Shot"
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
	}

	var a action.Info

	if skillActive && hold == attacks.AimParamLv1 {
		ai.Abil += " (A1)"
		a = action.Info{
			Frames:          frames.NewAbilFunc(aimedA1Frames),
			AnimationLength: aimedA1Frames[action.InvalidAction],
			CanQueueAfter:   aimedA1Hitmark,
			State:           action.AimState,
		}
	} else {
		a = action.Info{
			Frames:          frames.NewAbilFunc(aimedFrames[hold]),
			AnimationLength: aimedFrames[hold][action.InvalidAction],
			CanQueueAfter:   aimedHitmarks[hold],
			State:           action.AimState,
		}
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
		a.CanQueueAfter,
		a.CanQueueAfter+travel,
	)

	// Cover state handling - drops crowfeather, which explodes after 1.5 seconds
	if skillActive && hold == attacks.AimParamLv1 {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tengu Juurai: Ambush",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}
		ap := combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6)

		// TODO: snapshot?
		// Particles are emitted after the ambush thing hits
		c.Core.QueueAttack(ai, ap, aimedA1Hitmark, aimedA1Hitmark+travel+90, c.makeA4CB(), c.particleCB)
		c.attackBuff(ap, aimedA1Hitmark+travel+90)

		c.Core.Status.Delete(coverKey)
	}

	return a, nil
}
