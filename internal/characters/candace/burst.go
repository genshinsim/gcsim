package candace

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var burstFrames []int

const (
	burstHitmark = 102 // TODO: find correct hitmark
	burstKey     = "candace-burst"
	waveHitmark  = 16 // TODO: find correct hitmark
)

func init() {
	burstFrames = frames.InitAbilSlice(102)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.waveCount = 0
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sacred Rite: Wagtail's Tide (Q)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       burstDmg[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
		burstHitmark,
		burstHitmark,
	)

	duration := 540
	if c.Base.Cons >= 1 {
		duration = 720
	}

	// TODO: check if this is the right implementation
	for _, char := range c.Core.Player.Chars() {
		char.AddStatus(burstKey, duration, true) // TODO: find correct buff timing
		c.Core.Player.AddWeaponInfuse(
			c.Index,
			"candace-infuse",
			attributes.Hydro,
			duration,
			true,
			combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
		) // TODO: does this refresh constantly or one time?
	}
	c.ConsumeEnergy(4)                 // TODO: find correct energy timing
	c.SetCD(action.ActionBurst, 15*60) // TODO: find correct CD timing

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstSwap() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.waveCount > 3 {
			return false
		}
		next := args[1].(int)
		char := c.Core.Player.Chars()[next]
		if !char.StatusIsActive(burstKey) {
			return false
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Sacred Rite: Wagtail's Tide (Wave)",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Hydro,
			Durability: 25,
			Mult:       burstWaveDmg[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(c.Core.Combat.Player(), 1, false, combat.TargettableEnemy),
			waveHitmark,
			waveHitmark,
		) // TODO: find correct timing
		c.waveCount++
		return false
	}, "candace-burst-swap")
}
