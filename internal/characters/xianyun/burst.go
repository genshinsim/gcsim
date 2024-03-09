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
	burstHeal    = 4 // First heal is about 4f after hitmark
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
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) BurstCast() {
	// initial heal
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

		atk := c.getTotalAtk()
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Stars Gather at Dusk Heal (Initial)",
			Src:     healInstantP[c.TalentLvlBurst()]*atk + healInstantFlat[c.TalentLvlBurst()],
			Bonus:   c.Stat(attributes.Heal),
		})

		c.AddStatus(burstKey, burstDuration, false)

		c.c6()

		for _, char := range c.Core.Player.Chars() {
			// Due to the mechanism for how other characters check if they can do higher jumps
			// The other characters need to have the buff status on themselves.
			char.AddStatus(player.XianyunAirborneBuff, burstDuration, false)
		}

		c.adeptalAssistStacks = 8

		// TODO: From the frames sheet the heal timings are kind of all over the place
		for i := burstHeal; i <= burstHeal+burstDuration; i += 2.5 * 60 {
			// Unaffected by hitlag
			c.Core.Tasks.Add(c.BurstHealDoT, i)
		}

		c.a4StartUpdate()
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

		if !c.Core.Player.ActiveChar().StatusIsActive(player.XianyunAirborneBuff) {
			return false
		}

		if c.adeptalAssistStacks <= 0 {
			return false
		}

		aoe := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, burstDoTRadius)
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Starwicker",
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
			c.adeptalAssistStacks--
			if c.adeptalAssistStacks == 0 {
				// Delay stack reduction and status removal until after the attack lands
				// so that A4 can still proc on the attack that triggers the burstDot.
				for _, char := range c.Core.Player.Chars() {
					char.DeleteStatus(player.XianyunAirborneBuff)
				}
			}
		}, 1)
		return false
	}, "xianyun-starwicker-plunge-hook")
}

func (c *char) BurstHealDoT() {
	atk := c.getTotalAtk()
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Starwicker Heal",
		Src:     healDotP[c.TalentLvlBurst()]*atk + healDotFlat[c.TalentLvlBurst()],
		Bonus:   c.Stat(attributes.Heal),
	})
}
