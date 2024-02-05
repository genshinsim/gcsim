package xianyun

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var leapFrames []int
var plungeHitmarks = []int{20, 30, 40}
var plungeRadius = []float64{4, 5, 6.5}

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
	if !c.StatusIsActive(skillStateKey) {
		return action.Info{}, fmt.Errorf("xiangyun plunge used while not in cloud transmogrification state")
	}

	skillArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, plungeRadius[c.skillCounter-1])
	skillHitmark := plungeHitmarks[c.skillCounter-1]
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Chasing Crane %v", c.skillCounter),
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       leap[c.skillCounter-1][c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		skillArea,
		skillHitmark,
		skillHitmark,
		c.particleCB,
		c.a1cb(),
	)
	// reset window after leap
	c.DeleteStatus(skillStateKey)
	c.skillCounter = 0
	c.skillSrc = noSrcVal

	return action.Info{
		Frames:          frames.NewAbilFunc(leapFrames),
		State:           action.PlungeAttackState,
		AnimationLength: leapFrames[action.InvalidAction],
		CanQueueAfter:   leapFrames[action.ActionSkill],
	}, nil
}
