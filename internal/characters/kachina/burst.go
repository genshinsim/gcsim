package kachina

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const burstHitmark = 35

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(75)
	burstFrames[action.ActionAttack] = 56
	burstFrames[action.ActionCharge] = 56
	burstFrames[action.ActionSkill] = 53
	burstFrames[action.ActionDash] = 53
	burstFrames[action.ActionJump] = 53
	burstFrames[action.ActionSwap] = 50
}

func (c *char) Burst(_ map[string]int) (action.Info, error) {
	if c.Base.Cons >= 2 && !c.StatusIsActive(twirlyKey) {
		c.startTwirlyWithPoints(20, false)
	} else if c.Base.Cons >= 2 {
		c.nightsoulState.GeneratePoints(20)
	}

	c.fieldCenter = c.Core.Combat.Player().Pos()
	c.QueueCharTask(func() {
		c.AddStatus(fieldKey, int(field*60), true)
		c.applyC4(c.Core.Player.Active())
	}, burstHitmark)

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Time to Get Serious!",
		AttackTag:      attacks.AttackTagElementalBurst,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       150,
		Element:        attributes.Geo,
		Durability:     25,
		Mult:           burst[1][c.TalentLvlBurst()],
		UseDef:         true,
		IgnoreInfusion: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5.2), burstHitmark, burstHitmark)
	c.ConsumeEnergy(burstHitmark)
	c.SetCDWithDelay(action.ActionBurst, int(burst[0][c.TalentLvlBurst()]*60), 1)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}, nil
}
