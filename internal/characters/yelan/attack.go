package yelan

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{13}, {13}, {18}, {15, 29}}

const normalHitNum = 4

func init() {
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 15)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 21)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 38)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 67)
}

// Normal attack damage queue generator
// relatively standard with no major differences versus other bow characters
// Has "travel" parameter, used to set the number of frames that the arrow is in the air (default = 10)
func (c *char) Attack(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	if c.Base.Cons >= 6 && c.Core.Status.Duration(c6Status) > 0 {
		//c6 is default ICD group for some odd reason
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Breakthrough Barb",
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagExtraAttack,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Hydro,
			Durability: 25,
		}
		ai.FlatDmg = barb[c.TalentLvlAttack()] * c.MaxHP() * 1.56

		for i := range attack[c.NormalCounter] {
			c.c6count++
			if c.c6count >= 5 {
				c.Core.Status.Delete(c6Status) //delete status after 5 arrows
			}
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					4,
				),
				attackHitmarks[c.NormalCounter][i],
				attackHitmarks[c.NormalCounter][i]+travel,
			)
		}
	} else {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Physical,
			Durability: 25,
		}

		for i, mult := range attack[c.NormalCounter] {
			ai.Mult = mult[c.TalentLvlAttack()]
			c.Core.QueueAttack(
				ai,
				combat.NewBoxHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					combat.Point{Y: -0.5},
					0.1,
					1,
				),
				attackHitmarks[c.NormalCounter][i],
				attackHitmarks[c.NormalCounter][i]+travel,
			)
		}
	}

	defer c.AdvanceNormalIndex()

	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}
