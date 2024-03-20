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
	burstKey         = "gaming-q"
	burstStart       = int(0.6 * 60)
	burstDuration    = 12 * 60
	burstHitmark     = 60
	burstCD          = 15 * 60
	burstEnergyFrame = 7
	manChaiKey       = "gaming-man-chai"
	manChaiParam     = "man_chai_delay"
)

func init() {
	burstFrames = frames.InitAbilSlice(66)
	burstFrames[action.ActionAttack] = 63
	burstFrames[action.ActionSkill] = 63
	burstFrames[action.ActionJump] = 65
	burstFrames[action.ActionWalk] = 63
	burstFrames[action.ActionSwap] = 77
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if p[manChaiParam] > 0 {
		c.manChaiWalkBack = p[manChaiParam]
	} else {
		c.manChaiWalkBack = 100
	}

	c.Core.Tasks.Add(func() {
		c.AddStatus(burstKey, burstDuration, true)
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "Suanni's Gilded Dance (Q)",
			Type:    player.HealTypePercent,
			Src:     0.3,
			Bonus:   c.Stat(attributes.Heal),
		})
	}, burstStart)

	c.Core.Tasks.Add(c.queueManChai, burstHitmark+1)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Suanni's Gilded Dance (Q)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
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

	c.SetCD(action.ActionBurst, burstCD)
	c.ConsumeEnergy(burstEnergyFrame)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) queueManChai() {
	// new man chai can't spawn if one is currently active
	if c.StatusIsActive(manChaiKey) {
		return
	}
	c.AddStatus(manChaiKey, c.manChaiWalkBack, false)
	c.Core.Tasks.Add(func() {
		// can't link up if off-field
		if c.Core.Player.Active() != c.Index {
			return
		}
		c.ResetActionCooldown(action.ActionSkill)
		c.DeleteStatus(manChaiKey)
		c.c1()
	}, c.manChaiWalkBack)
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
