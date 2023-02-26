package raiden

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const (
	burstHitmark = 98
	BurstKey     = "raidenburst"
)

func init() {
	burstFrames = frames.InitAbilSlice(112) // Q -> J
	burstFrames[action.ActionAttack] = 111  // Q -> N1
	burstFrames[action.ActionCharge] = 500  // TODO: this action is illegal
	burstFrames[action.ActionSkill] = 111   // Q -> E
	burstFrames[action.ActionDash] = 111    // Q -> D
	burstFrames[action.ActionSwap] = 110    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// activate burst, reset stacks
	c.burstCastF = c.Core.F
	c.restoreCount = 0
	c.restoreICD = 0
	c.c6Count = 0
	c.c6ICD = 0

	// use a special modifier to track burst
	c.AddStatus(BurstKey, 420+burstHitmark, true)

	// apply when burst ends
	if c.Base.Cons >= 4 {
		c.applyC4 = true
		src := c.burstCastF
		c.QueueCharTask(func() {
			if src == c.burstCastF && c.applyC4 {
				c.applyC4 = false
				c.c4()
			}
		}, 420+burstHitmark)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Shinsetsu",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Electro,
		Durability: 50,
		Mult:       burstBase[c.TalentLvlBurst()],
	}

	if c.Base.Cons >= 2 {
		ai.IgnoreDefPercent = 0.6
	}

	c.Core.Tasks.Add(func() {
		c.stacksConsumed = c.stacks
		c.stacks = 0
		ai.Mult += resolveBaseBonus[c.TalentLvlBurst()] * c.stacksConsumed
		c.Core.Log.NewEvent("resolve stacks", glog.LogCharacterEvent, c.Index).
			Write("stacks", c.stacksConsumed)
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -0.1}, 13, 8),
			0,
			0,
		)
	}, burstHitmark)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstRestorefunc(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.Core.F > c.restoreICD && c.restoreCount < 5 {
		c.restoreCount++
		c.restoreICD = c.Core.F + 60 // once every 1 second
		energy := burstRestore[c.TalentLvlBurst()] * (1 + c.a4Energy(a.AttackEvent.Snapshot.Stats[attributes.ER]))
		for _, char := range c.Core.Player.Chars() {
			char.AddEnergy("raiden-burst", energy)
		}
	}
}

func (c *char) onSwapClearBurst() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if !c.StatusIsActive(BurstKey) {
			return false
		}
		// i prob don't need to check for who prev is here
		prev := args[0].(int)
		if prev == c.Index {
			c.DeleteStatus(BurstKey)
			if c.applyC4 {
				c.applyC4 = false
				c.c4()
			}
		}
		return false
	}, "raiden-burst-clear")
}

func (c *char) onBurstStackCount() {
	// TODO: this used to be on PostBurst; need to check if it works correctly still
	c.Core.Events.Subscribe(event.OnBurst, func(_ ...interface{}) bool {
		if c.Core.Player.Active() == c.Index {
			return false
		}
		char := c.Core.Player.ActiveChar()
		// add stacks based on char max energy
		stacks := resolveStackGain[c.TalentLvlBurst()] * char.EnergyMax
		if c.Base.Cons > 0 {
			if char.Base.Element == attributes.Electro {
				stacks = stacks * 1.8
			} else {
				stacks = stacks * 1.2
			}
		}
		c.stacks += stacks
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-stacks")
}
