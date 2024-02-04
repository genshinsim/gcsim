package faruzan

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	aimedFrames   [][]int
	aimedHitmarks = []int{15, 86}
	aimedEFrames  []int
)

const aimedEHitmark = 49

func init() {
	aimedFrames = make([][]int, 2)

	// Aimed Shot
	aimedFrames[0] = frames.InitAbilSlice(26)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Fully-Charged Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(96)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot (Hurricane Arrow)
	aimedEFrames = frames.InitAbilSlice(60)
	aimedEFrames[action.ActionDash] = aimedEHitmark
	aimedEFrames[action.ActionJump] = aimedEHitmark
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

	skillActive := c.StatusIsActive(skillKey) && c.hurricaneCount > 0
	// A1:
	// When Faruzan is in the Manifest Gale state created by Wind Realm of Nasamjnin,
	// the amount of time taken to charge a shot is decreased by 60%.
	shortAim := skillActive && c.Base.Ascension >= 1
	if skillActive {
		c.hurricaneCount -= 1
		if c.hurricaneCount <= 0 {
			c.DeleteStatus(skillKey)
		}
	}

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Anemo,
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
	var skillCb func(a combat.AttackCB)
	if skillActive && hold == attacks.AimParamLv1 {
		ai.Abil += " (Hurricane Arrow)"
		done := false
		skillCb = func(a combat.AttackCB) {
			if done {
				return
			}
			c.pressurizedCollapse(a.Target.Pos())
		}
		if shortAim {
			a = action.Info{
				Frames:          frames.NewAbilFunc(aimedEFrames),
				AnimationLength: aimedEFrames[action.InvalidAction],
				CanQueueAfter:   aimedEHitmark,
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
		skillCb,
	)

	return a, nil
}
