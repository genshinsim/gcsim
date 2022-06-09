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

const burstHitmark = 98

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//activate burst, reset stacks
	c.burstCastF = c.Core.F
	c.stacksConsumed = c.stacks
	c.stacks = 0
	c.Core.Status.Add("raidenburst", 420+burstHitmark) //7 seconds
	c.restoreCount = 0
	c.restoreICD = 0
	c.c6Count = 0
	c.c6ICD = 0

	// apply when burst ends
	if c.Base.Cons >= 4 {
		c.applyC4 = true
		src := c.burstCastF
		c.Core.Tasks.Add(func() {
			if src == c.burstCastF && c.applyC4 {
				c.applyC4 = false
				c.c4()
			}
		}, 420+burstHitmark)
	}

	if c.Base.Cons == 6 {
		c.c6Count = 0
	}

	c.Core.Log.NewEvent("resolve stacks", glog.LogCharacterEvent, c.Index, "stacks", c.stacksConsumed)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Musou Shinsetsu",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 50,
		Mult:       burstBase[c.TalentLvlBurst()],
	}
	ai.Mult += resolveBaseBonus[c.TalentLvlBurst()] * c.stacksConsumed
	if c.Base.Cons >= 2 {
		ai.IgnoreDefPercent = 0.6
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		Post:            burstHitmark,
		State:           action.BurstState,
	}
}

func (c *char) burstRestorefunc(a combat.AttackCB) {
	if c.Core.F > c.restoreICD && c.restoreCount < 5 {
		c.restoreCount++
		c.restoreICD = c.Core.F + 60 //once every 1 second
		energy := burstRestore[c.TalentLvlBurst()]
		//apply a4
		excess := int(a.AttackEvent.Snapshot.Stats[attributes.ER] / 0.01)
		c.Core.Log.NewEvent("a4 energy restore stacks", glog.LogCharacterEvent, c.Index, "stacks", excess, "increase", float64(excess)*0.006)
		energy = energy * (1 + float64(excess)*0.006)
		for _, char := range c.Core.Player.Chars() {
			char.AddEnergy("raiden-burst", energy)
		}
	}
}

func (c *char) onSwapClearBurst() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("raidenburst") == 0 {
			return false
		}
		//i prob don't need to check for who prev is here
		prev := args[0].(int)
		if prev == c.Index {
			c.Core.Status.Delete("raidenburst")
			if c.applyC4 {
				c.applyC4 = false
				c.c4()
			}
		}
		return false
	}, "raiden-burst-clear")
}

func (c *char) onBurstStackCount() {
	c.Core.Events.Subscribe(event.PostBurst, func(args ...interface{}) bool {
		if c.Core.Player.Active() == c.Index {
			return false
		}
		char := c.Core.Player.ActiveChar()
		//add stacks based on char max energy
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

	//a4 stack gain
	particleICD := 0
	c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		if particleICD > c.Core.F {
			return false
		}
		particleICD = c.Core.F + 180 // once every 3 seconds
		c.stacks += 2
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-particle-stacks")
}
