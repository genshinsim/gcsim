package fischl

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const burstHitmark = 18

func init() {
	burstFrames = frames.InitAbilSlice(148)
	burstFrames[action.ActionDash] = 111
	burstFrames[action.ActionJump] = 115
	burstFrames[action.ActionSwap] = 24
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//set on field oz to be this one
	//TODO: Oz should spawn and snapshot when the burst animation is cancelled
	//for now, the common burst->swap combo (24 frames) is used.
	c.Core.Tasks.Add(func() {
		c.queueOz("Burst")
	}, 24)

	//initial damage; part of the burst tag
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Midnight Phantasmagoria",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupFischl,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	//check for C4 damage
	if c.Base.Cons >= 4 {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Midnight Phantasmagoria",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupFischl,
			StrikeType: combat.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 50,
			Mult:       2.22,
		}
		// C4 damage always occurs before burst damage.
		c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 8, 8)
		//heal at end of animation
		heal := c.MaxHP() * 0.2
		c.Core.Tasks.Add(func() {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Index,
				Message: "Her Pilgrimage of Bleak",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			})
		}, burstHitmark) // TODO: should be at end of burst and not hitmark?

	}

	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 15*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel

		State: action.BurstState,
	}
}
