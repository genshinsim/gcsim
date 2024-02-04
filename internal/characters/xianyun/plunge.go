package xianyun

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var leapFrames []int
var plungeHoldFrames []int

// a1 is 1 frame before this
const plungePressHitmark = 36
const plungeHoldHitmark = 41

// TODO: missing plunge -> skill
func init() {
	// skill (press) -> high plunge -> x
	leapFrames = frames.InitAbilSlice(55) // max
	leapFrames[action.ActionDash] = 43
	leapFrames[action.ActionJump] = 50
	leapFrames[action.ActionSwap] = 50
}

func (c *char) HighPlungeAttack(p map[string]int) (action.Info, error) {
	// last action must be skill (for leap)
	if c.Core.Player.LastAction.Type != action.ActionSkill {
		return action.Info{}, errors.New("xiangyun plunge used without skill; last action before plunge must be skill")
	}

	act := action.Info{
		State: action.PlungeAttackState,
	}

	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, c.skillRadius)
	skillHitmark := 1
	act.Frames = frames.NewAbilFunc(leapFrames)
	act.AnimationLength = leapFrames[action.InvalidAction]
	act.CanQueueAfter = leapFrames[action.ActionSkill] // can only plunge after skill

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Chasing Crane %v", c.eCounter),
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       leap[c.eCounter-1][c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(skillArea.Shape.Pos(), nil, 8, 120),
		skillHitmark,
		skillHitmark,
		c.particleCB,
		c.a1(),
	)
	// reset window after leap
	c.DeleteStatus(eWindowKey)
	return act, nil
}
