package travelerdendro

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames [][]int

const burstHitmark = 91

func init() {
	burstFrames = make([][]int, 2)

	// Male
	burstFrames[0] = frames.InitAbilSlice(58)
	burstFrames[0][action.ActionSwap] = 57 // Q -> Swap

	// Female
	burstFrames[1] = frames.InitAbilSlice(58)
	burstFrames[1][action.ActionSwap] = 57 // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	c.SetCD(action.ActionBurst, 1200)
	c.ConsumeEnergy(2)

	procAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lea Lotus Lamp (Q)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}
	c.burstSnap = c.Snapshot(&procAI)
	c.burstAtk = &combat.AttackEvent{
		Info:     procAI,
		Snapshot: c.burstSnap,
	}

	burstDur := 12 * 60
	if c.Base.Cons >= 2 {
		burstDur += 3 * 60
	}
	c.burstExpire = c.Core.F + burstDur + burstHitmark

	c.Core.Status.Add("dmc-burst", burstDur+burstHitmark) // starts on first hitmark

	// A1 adds a stack per second
	for delay := c.Core.F + burstHitmark; delay < c.burstExpire; delay += 60 {
		c.a1Stack(delay)
	}

	// A1/C6 buff ticks every 0.3s and applies for 1s. probably counting from gadget spawn - Kolbiri
	for delay := c.Core.F + burstHitmark; delay < c.burstExpire; delay += 0.3 * 60 {
		c.a1Buff(delay)
	}

	if c.Base.Cons >= 6 {
		for delay := c.Core.F + burstHitmark; delay < c.burstExpire; delay += 0.3 * 60 {
			c.c6Buff(delay)
		}
	}

	c.burstTransfig = attributes.NoElement

	c.Core.Tasks.Add(c.burstTick, burstHitmark)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames[c.gender]),
		AnimationLength: burstFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   burstFrames[c.gender][action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

// This is a janked up version of DMC burst. This doesn't allow it to hold a Cryo aura,
// nor does it block transforms based on the gadget based on cryo blocking the application
// It also doesn't allow "missing" where ST attacks won't cause the lamp to transform
func (c *char) burstTransfigurationInit() {

	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)

		if c.Core.F <= c.burstExpire && //burst isn't expired
			c.burstTransfig == attributes.NoElement && //burst hasn't transfiged already
			ae.Info.Durability > 0 { // attack has gauge
			switch ae.Info.Element {
			case attributes.Electro:
				fallthrough
			case attributes.Hydro:
				fallthrough
			case attributes.Pyro:
				c.burstTransfig = ae.Info.Element
				c.Core.Log.NewEvent("dmc-burst-transfig-"+ae.Info.Element.String(), glog.LogCharacterEvent, c.Index)
				if c.Base.Cons >= 4 {
					c.c4()
				}
			case attributes.Cryo:
				c.Core.Log.NewEvent("This implementation of DMC burst is janky. Does not hold cryo auras right now", glog.LogCharacterEvent, c.Index)
			}
		}

		return false
	}, "dmc-lotuslight-transfiguration")
}

// timing is a bit off, the transfig timers should restart on transfig
func (c *char) burstTick() {

	thinkInterval := int(1.5 * 60)
	if c.Core.F > c.burstExpire {
		return
	}
	switch c.burstTransfig {
	case attributes.Pyro:
		c.burstAtk.Info.Abil = "Lea Lotus Lamp Explosion (Q)"
		c.burstAtk.Info.Durability = 50
		c.burstAtk.Info.ICDTag = combat.ICDTagNone
		c.burstAtk.Info.Mult = burstExplode[c.TalentLvlBurst()]

		// change this to gadget.self instead of Player() when gadgets are implemented
		// The timing on this explosion is also off
		c.Core.QueueAttackWithSnap(c.burstAtk.Info, c.burstAtk.Snapshot, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy), 0)
		c.burstExpire = c.Core.F
		c.Core.Status.Delete("dmc-burst")
		return
	case attributes.Electro:
		thinkInterval = int(0.9 * 60)
		c.Core.QueueAttackWithSnap(c.burstAtk.Info, c.burstAtk.Snapshot, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), 0)
	case attributes.Hydro:
		c.Core.QueueAttackWithSnap(c.burstAtk.Info, c.burstAtk.Snapshot, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy), 0)
	case attributes.NoElement:
		c.Core.QueueAttackWithSnap(c.burstAtk.Info, c.burstAtk.Snapshot, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy), 0)
	}
	c.Core.Tasks.Add(c.burstTick, thinkInterval)
}
