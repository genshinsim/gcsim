package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	a4Release   = []int{16, 18, 21, 25}
	dashFramesE []int
)

const a4Hitmark = 30

func init() {
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

	if c.StatusIsActive(skillKey) {
		return c.WindfavoredDash(p)
	}

	ai := c.Character.Dash(p)
	ai.Frames = func(action action.Action) int { return delay + ai.Frames(action) }
	ai.AnimationLength = delay + ai.AnimationLength
	ai.CanQueueAfter = delay + ai.CanQueueAfter

	return ai
}

func (c *char) WindfavoredDash(p map[string]int) action.ActionInfo {
	ai := action.ActionInfo{
		Frames:          func(next action.Action) int { return dashFramesE[next] },
		AnimationLength: dashFramesE[action.InvalidAction],
		CanQueueAfter:   dashFramesE[action.ActionSkill],
		State:           action.DashState,
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
				a4Release[i], a4Release[i]+a4Hitmark)
		}
	} else {
		// TODO: Check Point consumption
		c.skydwellerPoints -= 15
	}

	return ai
}
