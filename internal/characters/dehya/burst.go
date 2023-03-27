package dehya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var burstFrames []int
var kickFrames []int

const burstKey = "dehya-burst"
const burstDoT1Hitmark = 102
const fastPunchHitmark = 24 //10 hits max on 240 f
const slowPunchHitmark = 40 //6 hits minimum
const kickHitmark = 46      //6 hits minimum
var punchHitmarks = []int{42, 30, 30, 27, 27, 24, 24, 24, 24, 24, 24}

func init() {
	burstFrames = frames.InitAbilSlice(102) // Q -> E/D/J
	burstFrames[action.ActionSwap] = 102    // Q -> Swap

	kickFrames = frames.InitAbilSlice(101) // Q -> E/D/J
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// punches, ok := p["hits"]
	// if !ok || punches > 10 {
	// 	punches = 10
	// }
	if c.sanctumActive {
		c.sanctumActive = false
		c.sanctumExpiry += c.sanctumPickupExtension
		c.sanctumRetrieved = true
		c.sanctumICD = c.StatusDuration("dehya-skill-icd")
	}

	c.QueueCharTask(func() {
		c.AddStatus(burstKey, 280, false)
		c.burstCast = c.Core.F
		c.punchSrc = true
		c.burstCounter = 0
		c.burstPunch(c.punchSrc, true)
	}, burstDoT1Hitmark)

	c.ConsumeEnergy(5)
	c.SetCDWithDelay(action.ActionBurst, 18*60, 1)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.ActionAttack],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstPunch(src bool, auto bool) action.ActionInfo {

	hitmark := 44
	if !auto {
		hitmark = punchHitmarks[c.burstCounter]
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Flame-Mane's Fist",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurstPyro,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstPunchAtk[c.TalentLvlBurst()],
		FlatDmg:    (c.c1var + burstPunchHP[c.TalentLvlBurst()]) * c.MaxHP(),
	}

	c.QueueCharTask(func() {
		if c.punchSrc != src {
			return
		}
		if !c.StatusIsActive(burstKey) {
			return
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, 4),
			0,
			0,
		)
		if c.burstCast+240 > c.Core.F {
			c.burstCounter++
			c.punchSrc = true
			c.burstPunch(c.punchSrc, true)
		} else {
			c.punchSrc = true
			c.burstKick(c.punchSrc)
		}
	}, hitmark)
	if auto {
		return action.ActionInfo{
			Frames:          func(action.Action) int { return 44 },
			AnimationLength: 44,
			CanQueueAfter:   44,
			State:           action.BurstState,
		}
	}
	return action.ActionInfo{
		Frames:          func(action.Action) int { return punchHitmarks[c.burstCounter] },
		AnimationLength: punchHitmarks[c.burstCounter],
		CanQueueAfter:   punchHitmarks[c.burstCounter], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstKick(src bool) action.ActionInfo {
	if !c.StatusIsActive(burstKey) || src != c.punchSrc {
		return action.ActionInfo{
			Frames:          func(action.Action) int { return 0 },
			AnimationLength: 0,
			CanQueueAfter:   0,
			State:           action.BurstState,
		}
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Incineration Drive",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurstPyro,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstKickAtk[c.TalentLvlBurst()],
		FlatDmg:    (c.c1var + burstKickHP[c.TalentLvlBurst()]) * c.MaxHP(),
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, 4),
		kickHitmark,
		kickHitmark,
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(kickFrames),
		AnimationLength: kickFrames[action.ActionAttack],
		CanQueueAfter:   kickFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}
}
