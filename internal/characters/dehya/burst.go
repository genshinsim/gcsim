package dehya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int
var kickFrames []int
var remainingFieldDur int

const burstKey = "dehya-burst"
const burstDoT1Hitmark = 102
const kickHitmark = 46 //6 hits minimum
var punchHitmarks = []int{42, 30, 30, 27, 27, 24, 24, 24, 24, 24, 24}

func init() {
	//TODO:Deprecate bursty frames in favor of a constant?
	burstFrames = frames.InitAbilSlice(102) // Q -> E/D/J
	burstFrames[action.ActionSwap] = 102    // Q -> Swap

	kickFrames = frames.InitAbilSlice(72)       // Q -> Dash/Walk
	kickFrames[action.ActionAttack] = 75        // Q -> N1
	kickFrames[action.ActionSkill] = 71         // Q -> E
	kickFrames[action.ActionJump] = 73          // Q -> J
	kickFrames[action.ActionSwap] = kickHitmark //Q -> Swap

}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// punches, ok := p["hits"]
	// if !ok || punches > 10 {
	// 	punches = 10
	// }
	remainingFieldDur = 0
	if c.StatusIsActive(dehyaFieldKey) {
		// pick up field at start
		remainingFieldDur = c.StatusExpiry(dehyaFieldKey) + sanctumPickupExtension - c.Core.F //dur gets extended on field recast by a low margin, apparently
		c.Core.Log.NewEvent("sanctum removed", glog.LogCharacterEvent, c.Index).
			Write("Duration Remaining ", remainingFieldDur+sanctumPickupExtension).
			Write("DoT tick CD", c.StatusDuration(skillICDKey))
	}
	c.DeleteStatus(dehyaFieldKey)

	c.QueueCharTask(func() {
		c.AddStatus(burstKey, 280, false)
		c.burstCast = c.Core.F
		c.punchSrc = true
		c.burstCounter = 0
		c.burstPunch(c.punchSrc, true)
	}, burstDoT1Hitmark)

	c.ConsumeEnergy(15) //TODO: If this is ping related, this could be closer to 1 at 0 ping
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
		FlatDmg:    (c.c1var[0] + burstPunchHP[c.TalentLvlBurst()]) * c.MaxHP(),
	}

	c.Core.Tasks.Add(func() {
		if c.punchSrc != src {
			return
		}
		if !c.StatusIsActive(burstKey) {
			return
		}
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -2.8}, 5, 7.8),
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
		FlatDmg:    (c.c1var[0] + burstKickHP[c.TalentLvlBurst()]) * c.MaxHP(),
	}
	if remainingFieldDur > 0 {
		c.QueueCharTask(func() { //place field
			c.addField(remainingFieldDur)
		}, kickHitmark)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 6.5),
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
