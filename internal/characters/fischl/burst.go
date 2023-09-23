package fischl

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const (
	burstHitmark          = 18
	burstFullOzSpawn      = 113 // from start of action
	burstFullOzFirstTick  = 69
	burstShortOzSpawn     = 1 // after swap occurs
	burstShortOzFirstTick = 63
)

func init() {
	burstFrames = frames.InitAbilSlice(148)
	burstFrames[action.ActionDash] = 115 // sheet assumed wrong dash frames
	burstFrames[action.ActionJump] = 115
	burstFrames[action.ActionSwap] = 24
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// initial damage; part of the burst tag
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Midnight Phantasmagoria",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupFischl,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 0.5),
		burstHitmark,
		burstHitmark,
	)

	// check for C4 damage
	if c.Base.Cons >= 4 {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Her Pilgrimage of Bleak (C4)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupFischl,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 50,
			Mult:       2.22,
		}
		// C4 damage always occurs before burst damage.
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 8, 8)
		// heal at end of animation
		heal := c.MaxHP() * 0.2
		c.Core.Tasks.Add(func() {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Index,
				Message: "Her Pilgrimage of Bleak (C4)",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			})
		}, burstHitmark) // TODO: should be at end of burst and not hitmark?
	}

	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 15*60)

	// set oz to active at the start of the action
	c.ozActive = true
	c.burstOzSpawnSrc = c.Core.F
	burstFullOzFunc := c.burstOzSpawn(c.Core.F, 0, burstFullOzFirstTick)
	burstShortOzFunc := c.burstOzSpawn(c.Core.F, burstShortOzSpawn, burstShortOzFirstTick)

	c.Core.Tasks.Add(burstFullOzFunc, burstFullOzSpawn)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
		OnRemoved:       func(next action.AnimationState) { burstShortOzFunc() },
	}, nil
}

func (c *char) burstOzSpawn(src, ozSpawn, firstTick int) func() {
	return func() {
		if src != c.burstOzSpawnSrc {
			return
		}
		c.burstOzSpawnSrc = -1
		c.queueOz("Burst", ozSpawn, firstTick)
	}
}
