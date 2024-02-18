package xianyun

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
	StarwickerKey = player.XianyunAirborneBuff

	burstHeal    = 2.5 * 60
	burstHitmark = 75
	burstKey     = "xianyun-burst"
	// 16 seconds duration
	burstDuration  = 16 * 60
	burstRadius    = 7
	burstDoTRadius = 4.8
	burstDoTDelay  = 5
)

// TODO: dummy frame data from shenhe
func init() {
	burstFrames = frames.InitAbilSlice(103) // Q -> J
	burstFrames[action.ActionAttack] = 101  // Q -> N1
	burstFrames[action.ActionCharge] = 102  // Q -> CA
	burstFrames[action.ActionSkill] = 101   // Q -> E
	burstFrames[action.ActionDash] = 101    // Q -> D
	burstFrames[action.ActionWalk] = 101    // Q -> Walk
	burstFrames[action.ActionSwap] = 99     // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(18)
	c.BurstCast()
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) BurstCast() {
	// init heal
	c.QueueCharTask(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Stars Gather at Dusk (Initial)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       burst[c.TalentLvlBurst()],
		}

		burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7)
		c.Core.QueueAttack(ai, burstArea, 0, 0)

		atk := c.Base.Atk*(1+c.Stat(attributes.ATKP)) + c.Stat(attributes.ATK)
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Starwicker-Heal-Initial",
			Src:     healInstantP[c.TalentLvlBurst()]*atk + healInstantFlat[c.TalentLvlBurst()],
			Bonus:   c.Stat(attributes.Heal),
		})

		c.AddStatus(burstKey, burstDuration, false)
		for _, char := range c.Core.Player.Chars() {
			char.AddStatus(StarwickerKey, burstDuration, true)
		}
		c.starwickerStacks = 8

		for i := burstHeal; i <= burstHeal+burstDuration; i += 2.5 * 60 {
			c.Core.Tasks.Add(c.BurstHealDoT, int(i))
		}
	}, burstHitmark)
}

func (c *char) burstPlungeDoTTrigger() {
	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...interface{}) bool {
		// ApplyAttack occurs only once per attack, so we do not need to add an ICD status
		atk := args[0].(*combat.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}

		if atk.Info.Durability == 0 {
			// plunge collisions have 0 durability
			return false
		}

		if !c.Core.Player.ActiveChar().StatusIsActive(StarwickerKey) {
			return false
		}

		if c.starwickerStacks <= 0 {
			return false
		}

		aoe := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, burstDoTRadius)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Starwicker Plunge DoT Damage",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       burstDot[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(
			ai,
			aoe,
			burstDoTDelay,
			burstDoTDelay,
		)
		c.QueueCharTask(func() {
			c.starwickerStacks--
			if c.starwickerStacks == 0 {
				// Delay stack reduction and status removal until after the burstDotDelay
				// so that A4 can still proc on the attack that triggers the burstDot.
				for _, char := range c.Core.Player.Chars() {
					char.DeleteStatus(StarwickerKey)
				}
			}
		}, burstDoTDelay)
		return false
	}, "xianyun-starwicker-plunge-DoT-hook")
}

func (c *char) BurstHealDoT() {
	atk := c.Base.Atk*(1+c.Stat(attributes.ATKP)) + c.Stat(attributes.ATK)
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Starwicker-Heal-DoT",
		Src:     healDotP[c.TalentLvlBurst()]*atk + healDotFlat[c.TalentLvlBurst()],
		Bonus:   c.Stat(attributes.Heal),
	})
}
