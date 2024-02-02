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

const burstKey = "dehya-burst"
const kickKey = "dehya-burst-kick"
const burstDoT1Hitmark = 105
const kickHitmark = 46 // 6 hits minimum
const punchSlowHitmark = 43

var burstFrames []int
var kickFrames []int
var punchHitmarks = []int{30, 30, 28, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27}

func init() {
	//TODO: Deprecate bursty frames in favor of a constant?
	burstFrames = frames.InitAbilSlice(102) // Q -> E/D/J
	burstFrames[action.ActionSwap] = 102    // Q -> Swap

	kickFrames = frames.InitAbilSlice(72)       // Q -> Dash/Walk
	kickFrames[action.ActionAttack] = 75        // Q -> N1
	kickFrames[action.ActionSkill] = 71         // Q -> E
	kickFrames[action.ActionJump] = 73          // Q -> J
	kickFrames[action.ActionSwap] = kickHitmark // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.burstJumpCancel = false
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Flame-Mane's Fist",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   50,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstPunchAtk[c.TalentLvlBurst()],
		FlatDmg:    (c.c1var[0] + burstPunchHP[c.TalentLvlBurst()]) * c.MaxHP(),
	}

	c.c6count = 0
	c.sanctumSavedDur = 0
	if c.StatusIsActive(dehyaFieldKey) {
		// pick up field at start
		c.sanctumSavedDur = c.StatusExpiry(dehyaFieldKey) + sanctumPickupExtension - c.Core.F // dur gets extended on field recast by a low margin, apparently
		c.Core.Log.NewEvent("sanctum removed", glog.LogCharacterEvent, c.Index).
			Write("Duration Remaining ", c.sanctumSavedDur+sanctumPickupExtension).
			Write("DoT tick CD", c.StatusDuration(skillICDKey))
		c.DeleteStatus(dehyaFieldKey)
	}

	c.Core.Tasks.Add(func() {
		c.AddStatus(burstKey, 240, false)
		c.burstCast = c.Core.F
		c.burstHitSrc = 0
		c.burstCounter = 0
		c.burstPunch(c.burstHitSrc, true)
	}, burstDoT1Hitmark)

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -2.8}, 5, 7.8),
		burstDoT1Hitmark,
		burstDoT1Hitmark,
		c.c4cb(),
		c.c6cb(),
	)

	c.ConsumeEnergy(15) //TODO: If this is ping related, this could be closer to 1 at 0 ping
	c.SetCDWithDelay(action.ActionBurst, 18*60, 1)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.ActionAttack],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstPunch(src int, auto bool) action.Info {
	hitmark := punchSlowHitmark
	if !auto {
		hitmark = punchHitmarks[c.burstCounter]
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Flame-Mane's Fist",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   50,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstPunchAtk[c.TalentLvlBurst()],
		FlatDmg:    (c.c1var[0] + burstPunchHP[c.TalentLvlBurst()]) * c.MaxHP(),
	}

	c.Core.Tasks.Add(func() {
		if c.burstHitSrc != src {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if c.burstJumpCancel { // prevent punches if you jump cancel burst
			return
		}
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -2.8}, 5, 7.8),
			0,
			0,
			c.c4cb(),
			c.c6cb(),
		)
		if !c.StatusIsActive(burstKey) {
			c.burstHitSrc++
			c.AddStatus(kickKey, kickHitmark, false)
			c.burstKick(c.burstHitSrc)

			return
		}
		c.burstCounter++
		c.burstHitSrc++
		c.burstPunch(c.burstHitSrc, true)
	}, hitmark)

	return action.Info{
		Frames:          func(action.Action) int { return hitmark },
		AnimationLength: hitmark,
		CanQueueAfter:   hitmark,
		State:           action.BurstState,
	}
}

func (c *char) burstKick(src int) action.Info {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Incineration Drive",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   100,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstKickAtk[c.TalentLvlBurst()],
		FlatDmg:    (c.c1var[0] + burstKickHP[c.TalentLvlBurst()]) * c.MaxHP(),
	}

	c.Core.Tasks.Add(func() {
		if src != c.burstHitSrc { // prevents duplicates
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 6.5),
			0,
			0,
			c.c4cb(),
		)
		if dur := c.sanctumSavedDur; dur > 0 { // place field
			c.sanctumSavedDur = 0
			c.addField(dur)
		}
	}, kickHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(kickFrames),
		AnimationLength: kickFrames[action.ActionAttack],
		CanQueueAfter:   kickFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) UseBurstAction() *action.Info {
	var out action.Info
	c.burstHitSrc++
	if c.StatusIsActive(kickKey) {
		out = c.burstKick(c.burstHitSrc)
		return &out
	}
	if c.StatusIsActive(burstKey) {
		out = c.burstPunch(c.burstHitSrc, false)
		return &out
	}
	return nil
}
