package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var (
	a4Release        = []int{16, 18, 21, 25}
	dashFramesNormal []int
	dashFramesE      []int
)

const a4Hitmark = 30

func init() {
	dashFramesNormal = frames.InitAbilSlice(21)

	dashFramesE = frames.InitAbilSlice(24)
	dashFramesE[action.ActionAttack] = 21
	dashFramesE[action.ActionCharge] = 21
	dashFramesE[action.ActionSkill] = 5
	dashFramesE[action.ActionDash] = 22
	dashFramesE[action.ActionJump] = 22
	dashFramesE[action.ActionWalk] = 22
}

func (c *char) Dash(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	relevantFrames := dashFramesNormal
	if c.StatusIsActive(skillKey) {
		relevantFrames = dashFramesE
	}

	ai := action.ActionInfo{
		Frames:          func(next action.Action) int { return delay + relevantFrames[next] },
		AnimationLength: delay + relevantFrames[action.InvalidAction],
		CanQueueAfter:   delay + relevantFrames[action.ActionSkill],
		State:           action.DashState,
	}

	if !c.StatusIsActive(skillKey) {
		c.Core.Tasks.Add(func() {
			req := c.Core.Player.AbilStamCost(c.Index, action.ActionDash, p)
			c.Core.Player.Stam -= req
			//this really shouldn't happen??
			if c.Core.Player.Stam < 0 {
				c.Core.Player.Stam = 0
			}
			c.Core.Player.LastStamUse = delay + c.Core.F
			c.Core.Events.Emit(event.OnStamUse, action.DashState)
		}, delay+dashFramesNormal[action.ActionAttack]-1)

		return ai
	}

	if c.StatusIsActive(a4Key) {
		c.DeleteStatus(a4Key)

		a4Mult := 0.35

		// TODO: Write as attack mod?
		if c.StatusIsActive("wanderer-c1-atkspd") {
			a4Mult = 0.6
		}

		a4Info := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Gales of Reverie",
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagWandererA4,
			ICDGroup:   combat.ICDGroupWandererA4,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       a4Mult,
		}

		for i := 0; i < 4; i++ {
			c.Core.QueueAttack(a4Info, combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 1),
				delay+a4Release[i], delay+a4Release[i]+a4Hitmark)
		}
	} else {
		// TODO: Check Point consumption
		c.skydwellerPoints -= 15
	}

	return ai
}
