package gaming

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const (
	burstKey      = "gaming-wushou"
	lionKey       = "gaming-manchai"
	burstHitmark  = 44 // TODO
	gamingBurstCd = 15 // seconds
)

func init() {
	// TODO : TAKEN FROM FREMINET
	burstFrames = frames.InitAbilSlice(65)
	burstFrames[action.ActionAttack] = 52
	burstFrames[action.ActionSkill] = 52
	burstFrames[action.ActionDash] = 53
	burstFrames[action.ActionJump] = 52
	burstFrames[action.ActionSwap] = 51
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.AddStatus(burstKey, 12*60, true)
	// spawn manchai manually
	if !c.StatusIsActive(lionKey) {
		c.AddStatus(lionKey, lionWalkBack, false)
		c.QueueCharTask(func() {
			c.ResetActionCooldown(action.ActionSkill)
			c.DeleteStatus(lionKey)
			c.c1()
		}, lionWalkBack)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Horned Lion's Gilded Dance (Q)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
		burstHitmark,
		burstHitmark,
	)

	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: "Horned Lion's Gilded Dance (Q)",
		Type:    player.HealTypePercent,
		Src:     0.3,
		Bonus:   c.Stat(attributes.Heal),
	})

	c.SetCD(action.ActionBurst, 60*gamingBurstCd)
	c.ConsumeEnergy(4) // TODO

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// do nothing if previous char wasn't gaming
		prev := args[0].(int)
		if prev != c.Index {
			return false
		}
		if !c.StatusIsActive(burstKey) {
			return false
		}
		c.DeleteStatus(burstKey)

		return false
	}, "gaming-exit")
}
