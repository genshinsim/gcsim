package dehya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const (
	burstKey              = "dehya-burst"
	burstDuration         = 4.1 * 60
	kickKey               = "dehya-burst-kick"
	burstPunch1Hitmark    = 105
	burstPunchSlowHitmark = 50
	burstKickHitmark      = 46
)

var (
	kickFrames    []int
	punchHitmarks = []int{30, 30, 28, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27}
)

func init() {
	kickFrames = frames.InitAbilSlice(76) // Q -> N1
	kickFrames[action.ActionSkill] = 71
	kickFrames[action.ActionDash] = 73
	kickFrames[action.ActionJump] = 73
	kickFrames[action.ActionSwap] = burstKickHitmark
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.c6Count = 0
	c.sanctumSavedDur = 0
	if c.StatusIsActive(dehyaFieldKey) {
		// pick up field at start
		c.pickUpField()
	}

	c.Core.Tasks.Add(func() {
		c.burstHitSrc = 0
		c.burstCounter = 0
		c.AddStatus(burstKey, burstDuration, true)
		c.burstPunchFunc(c.burstHitSrc)()
	}, burstPunch1Hitmark)

	c.ConsumeEnergy(15) //TODO: If this is ping related, this could be closer to 1 at 0 ping
	c.SetCDWithDelay(action.ActionBurst, 18*60, 1)

	return action.Info{
		Frames:          func(action.Action) int { return burstPunch1Hitmark },
		AnimationLength: burstPunch1Hitmark,
		CanQueueAfter:   burstPunch1Hitmark,
		State:           action.BurstState,
	}, nil
}

func (c *char) burstPunchFunc(src int) func() {
	return func() {
		if c.burstHitSrc != src {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
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
			FlatDmg:    (c.c1FlatDmgRatioQ + burstPunchHP[c.TalentLvlBurst()]) * c.MaxHP(),
		}
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -2.8}, 5, 7.8),
			0,
			0,
			c.c4CB(),
			c.c6CB(),
		)
		if !c.StatusIsActive(burstKey) {
			c.burstHitSrc++
			c.AddStatus(kickKey, burstKickHitmark, true)
			c.Core.Tasks.Add(c.burstKickFunc(c.burstHitSrc), burstKickHitmark)
			return
		}
		c.burstCounter++
		c.burstHitSrc++
		c.Core.Tasks.Add(c.burstPunchFunc(c.burstHitSrc), burstPunchSlowHitmark)
	}
}

func (c *char) burstKickFunc(src int) func() {
	return func() {
		if src != c.burstHitSrc { // prevents duplicates
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
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
			FlatDmg:    (c.c1FlatDmgRatioQ + burstKickHP[c.TalentLvlBurst()]) * c.MaxHP(),
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 1}, 6.5),
			0,
			0,
			c.c4CB(),
		)
		if dur := c.sanctumSavedDur; dur > 0 { // place field with 1f delay to avoid self-trigger
			c.sanctumSavedDur = 0
			c.Core.Tasks.Add(func() {
				c.AddStatus(skillICDKey, c.sanctumICD, false)
				c.addField(dur)
			}, 1)
		}
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

func (c *char) burstPunch(src int, auto bool) action.Info {
	hitmark := burstPunchSlowHitmark
	if !auto {
		hitmark = punchHitmarks[c.burstCounter]
	}

	c.Core.Tasks.Add(c.burstPunchFunc(src), hitmark)

	return action.Info{
		Frames:          func(action.Action) int { return hitmark },
		AnimationLength: hitmark,
		CanQueueAfter:   hitmark,
		State:           action.Idle, // TODO: cannot use burst state because burst state implies iframes
	}
}

func (c *char) burstKick(src int) action.Info {
	c.Core.Tasks.Add(c.burstKickFunc(src), burstKickHitmark)
	return action.Info{
		Frames:          frames.NewAbilFunc(kickFrames),
		AnimationLength: kickFrames[action.ActionAttack],
		CanQueueAfter:   kickFrames[action.ActionSwap], // earliest cancel
		State:           action.Idle,                   // TODO: cannot use burst state because burst state implies iframes
	}
}
