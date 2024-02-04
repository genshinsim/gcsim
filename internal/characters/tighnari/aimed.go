package tighnari

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
var aimedHitmarks = []int{14, 86}
var aimedWreathFrames []int

const aimedWreathHitmark = 175

func init() {
	aimedFrames = make([][]int, 2)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(23)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(94)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot (Wreath Arrow)
	aimedWreathFrames = frames.InitAbilSlice(183)
	aimedWreathFrames[action.ActionDash] = aimedWreathHitmark
	aimedWreathFrames[action.ActionJump] = aimedWreathHitmark
}

func (c *char) Aimed(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if !ok {
		hold = attacks.AimParamLv1
	}
	if c.StatusIsActive(vijnanasuffusionStatus) {
		hold = attacks.AimParamLv2
	}
	switch hold {
	case attacks.AimParamPhys:
	case attacks.AimParamLv1:
	case attacks.AimParamLv2:
		return c.WreathAimed(p)
	default:
		return action.Info{}, fmt.Errorf("invalid hold param supplied, got %v", hold)
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Dendro,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
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

func (c *char) WreathAimed(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	wreathTravel, ok := p["wreath"]
	if !ok {
		wreathTravel = 35
	}
	weakspot := p["weakspot"]

	skip := 0
	if c.StatusIsActive(vijnanasuffusionStatus) {
		skip = 142 // 2.4 * 60

		arrows := c.Tag(wreatharrows) - 1
		c.SetTag(wreatharrows, arrows)
		if arrows == 0 {
			c.DeleteStatus(vijnanasuffusionStatus)
		}
	}
	if c.Base.Cons >= 6 {
		skip += 0.9 * 60
	}
	if skip > aimedWreathHitmark {
		skip = aimedWreathHitmark
	}

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot (Wreath Arrow)",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Dendro,
		Durability:           25,
		Mult:                 wreath[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
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
		aimedWreathHitmark-skip,
		aimedWreathHitmark+travel-skip,
	)
	if c.Base.Ascension >= 1 {
		c.Core.Tasks.Add(c.a1, aimedWreathHitmark-skip+1)
	}

	ai = combat.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Clusterbloom Arrow",
		AttackTag:    attacks.AttackTagExtra,
		ICDTag:       attacks.ICDTagExtraAttack,
		ICDGroup:     attacks.ICDGroupTighnari,
		StrikeType:   attacks.StrikeTypePierce,
		Element:      attributes.Dendro,
		Durability:   25,
		Mult:         clusterbloom[c.TalentLvlAttack()],
		HitWeakPoint: false, // TODO: tighnari can hit the weak spot on some enemies (like hilichurls)
	}
	c.Core.Tasks.Add(func() {
		snap := c.Snapshot(&ai)
		for i := 0; i < 4; i++ {
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					1,
				),
				wreathTravel,
			)
		}

		if c.Base.Cons >= 6 {
			ai = combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Karma Adjudged From the Leaden Fruit",
				AttackTag:  attacks.AttackTagExtra,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypePierce,
				Element:    attributes.Dendro,
				Durability: 25,
				Mult:       1.5,
			}
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					1,
				),
				wreathTravel,
			)
		}
	}, aimedWreathHitmark+travel-skip)

	return action.Info{
		Frames:          func(next action.Action) int { return aimedWreathFrames[next] - skip },
		AnimationLength: aimedWreathFrames[action.InvalidAction] - skip,
		CanQueueAfter:   aimedWreathHitmark - skip,
		State:           action.AimState,
	}, nil
}
