package wanderer

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var lowPlungeFrames []int

const lowPlungeHitmark = 41
const lowPlungeCollisionHitmark = 36

func init() {
	lowPlungeFrames = frames.InitAbilSlice(72)
	lowPlungeFrames[action.ActionAttack] = 65
	lowPlungeFrames[action.ActionCharge] = 64
	lowPlungeFrames[action.ActionBurst] = 65
	lowPlungeFrames[action.ActionDash] = 41
	lowPlungeFrames[action.ActionSwap] = 57
}

func (c *char) LowPlungeAttack(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	// Not in falling state
	if !c.StatusIsActive(plungeAvailableKey) {
		c.Core.Log.NewEvent("only plunge after skill ends", glog.LogActionEvent, c.Index).
			Write("action", action.ActionLowPlunge)
		return action.ActionInfo{
			Frames:          func(action.Action) int { return 1200 },
			AnimationLength: 1200,
			CanQueueAfter:   1200,
			State:           action.Idle,
		}
	}
	c.DeleteStatus(plungeAvailableKey)

	// Decreasing delay due to casting midair
	if delay > 0 {
		delay = 7
	}

	collision, ok := p["collision"]
	if !ok {
		collision = 0 // Whether or not Wanderer does a collision hit
	}

	if collision > 0 {
		c.plungeCollision(lowPlungeCollisionHitmark + delay)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge Attack",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       lowPlunge[c.TalentLvlAttack()],
	}

	// TODO: check snapshot delay
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		delay+lowPlungeHitmark, delay+lowPlungeHitmark)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return lowPlungeFrames[next] },
		AnimationLength: lowPlungeFrames[action.InvalidAction],
		CanQueueAfter:   lowPlungeHitmark,
		State:           action.PlungeAttackState,
	}
}

func (c *char) plungeCollision(fullDelay int) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Plunge Collision",
		AttackTag:  attacks.AttackTagPlunge,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 0,
		Mult:       plunge[c.TalentLvlAttack()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1.5), fullDelay, fullDelay)
}
