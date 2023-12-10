package freminet

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var burstFrames []int

const (
	burstKey     = "freminet-stalking"
	burstHitmark = 44
)

func init() {
	burstFrames = frames.InitAbilSlice(65)
	burstFrames[action.ActionAttack] = 52
	burstFrames[action.ActionSkill] = 52
	burstFrames[action.ActionDash] = 53
	burstFrames[action.ActionJump] = 52
	burstFrames[action.ActionSwap] = 51
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.AddStatus(burstKey, 10*60, true)

	c.ResetActionCooldown(action.ActionSkill)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Shadowhunter's Ambush",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
		burstHitmark,
		burstHitmark,
	)

	c.SetCD(action.ActionBurst, 60*15)
	c.ConsumeEnergy(4)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// do nothing if previous char wasn't freminet
		prev := args[0].(int)
		if prev != c.Index {
			return false
		}
		if !c.StatusIsActive(burstKey) {
			return false
		}
		c.DeleteStatus(burstKey)

		return false
	}, "freminet-exit")
}
