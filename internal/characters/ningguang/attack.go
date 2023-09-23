package ningguang

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var (
	attackFrames   [][]int
	attackLockout  []int
	attackHitmarks []int
	attackOptions  = map[attackType][]attackType{
		attackTypeLeft:  {attackTypeRight, attackTypeTwirl},
		attackTypeRight: {attackTypeLeft, attackTypeTwirl},
		attackTypeTwirl: {attackTypeLeft, attackTypeRight},
	}
)

const normalHitNum = 1

func init() {
	attackLockout = make([]int, endAttackType)
	attackLockout[attackTypeLeft] = 15
	attackLockout[attackTypeRight] = 5
	attackLockout[attackTypeTwirl] = 13

	attackHitmarks = make([]int, endAttackType)
	attackHitmarks[attackTypeLeft] = 29
	attackHitmarks[attackTypeRight] = 19
	attackHitmarks[attackTypeTwirl] = 27

	attackFrames = make([][]int, endAttackType)
	// NA Left -> x
	attackFrames[attackTypeLeft] = frames.InitNormalCancelSlice(attackLockout[attackTypeLeft], 61)
	attackFrames[attackTypeLeft][action.ActionCharge] = 42
	attackFrames[attackTypeLeft][action.ActionWalk] = 44
	// NA Right -> x
	attackFrames[attackTypeRight] = frames.InitNormalCancelSlice(attackLockout[attackTypeRight], 56)
	attackFrames[attackTypeRight][action.ActionCharge] = 40
	attackFrames[attackTypeRight][action.ActionWalk] = 41
	// NA Twirl -> x
	attackFrames[attackTypeTwirl] = frames.InitNormalCancelSlice(attackLockout[attackTypeTwirl], 66)
	attackFrames[attackTypeTwirl][action.ActionCharge] = 40
	attackFrames[attackTypeTwirl][action.ActionWalk] = 42
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	done := false
	cb := func(_ combat.AttackCB) {
		// doesn't gain jades off-field
		if c.Core.Player.Active() != c.Index {
			return
		}
		if done {
			return
		}
		count := c.jadeCount
		// if we're at 7 dont increase but also dont reset back to 3
		if count != 7 {
			count++
			if count > 3 {
				count = 3
			} else {
				c.Core.Log.NewEvent("adding star jade", glog.LogCharacterEvent, c.Index).
					Write("count", count)
			}
			c.jadeCount = count
		}
		done = true
	}

	r := 0.5
	if c.Base.Cons >= 1 {
		r = 3.5
	}

	nextAttack := attackOptions[c.prevAttack][c.Core.Rand.Intn(2)]
	if c.Core.Player.CurrentState() == action.DashState { // dash > NA will always be left attack
		nextAttack = attackTypeLeft
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal (%s)", nextAttack),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       attack[c.TalentLvlAttack()],
	}

	for i := 0; i < 2; i++ {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				r,
			),
			attackHitmarks[nextAttack],
			attackHitmarks[nextAttack]+travel,
			cb,
		)
	}

	c.prevAttack = nextAttack

	return action.Info{
		Frames: func(next action.Action) int {
			return frames.AtkSpdAdjust(attackFrames[nextAttack][next], c.Stat(attributes.AtkSpd))
		},
		AnimationLength: attackFrames[nextAttack][action.InvalidAction],
		CanQueueAfter:   attackLockout[nextAttack],
		State:           action.NormalAttackState,
	}, nil
}
