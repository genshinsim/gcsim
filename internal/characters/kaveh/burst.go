package kaveh

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark     = 36
	burstDuration    = 720
	burstKey         = "kaveh-q"
	burstDmgBonusKey = "kaveh-q-dmg-bonus"
)

func init() {
	burstFrames = frames.InitAbilSlice(49)
	burstFrames[action.ActionAttack] = 48
	burstFrames[action.ActionDash] = 44
	burstFrames[action.ActionJump] = 44
	burstFrames[action.ActionWalk] = 48
	burstFrames[action.ActionSwap] = 42
}

func (c *char) Burst(p map[string]int) action.Info {
	c.a4Stacks = 0

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Painted Dome (Q)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	c.Core.QueueAttack(ai, ap, burstHitmark, burstHitmark)
	c.SetCD(action.ActionBurst, 1200)
	c.ConsumeEnergy(3)

	c.Core.Tasks.Add(func() {
		c.ruptureDendroCores(ap)
		c.AddStatus(burstKey, burstDuration, true)
		c.a4()
		if c.Base.Cons >= 2 {
			c.c2()
		}
		for _, char := range c.Core.Player.Chars() {
			char.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBaseWithHitlag(burstDmgBonusKey, burstDuration),
				Amount: func(ai combat.AttackInfo) (float64, bool) {
					if ai.AttackTag == attacks.AttackTagBloom {
						return burstDmgBonus[c.TalentLvlBurst()], false
					}
					return 0, false
				},
			})
		}
	}, burstHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) addBurstExitHandler() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		c.DeleteStatus(burstKey)
		c.DeleteStatus(a4Key)
		c.DeleteStatus(c2Key)
		for _, char := range c.Core.Player.Chars() {
			char.DeleteStatus(burstDmgBonusKey)
		}
		return false
	}, "kaveh-exit")
}
