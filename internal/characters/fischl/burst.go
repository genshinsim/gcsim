package fischl

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
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
		PoiseDMG:   150,
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

	// check for C4
	var c4HealFunc func()
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

		// heal
		c4HealFunc = func() {
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index,
				Target:  c.Index,
				Message: "Her Pilgrimage of Bleak (C4)",
				Src:     0.2 * c.MaxHP(),
				Bonus:   c.Stat(attributes.Heal),
			})
		}
	}

	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 15*60)

	// set oz to active at the start of the action
	c.ozActive = true
	c.burstOzSpawnSrc = c.Core.F
	burstFullOzFunc := c.burstOzSpawn(c.Core.F, 0, burstFullOzFirstTick, c4HealFunc)
	burstShortOzFunc := c.burstOzSpawn(c.Core.F, burstShortOzSpawn, burstShortOzFirstTick, c4HealFunc)

	c.Core.Tasks.Add(burstFullOzFunc, burstFullOzSpawn)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
		OnRemoved:       func(next action.AnimationState) { burstShortOzFunc() },
	}, nil
}

func (c *char) burstOzSpawn(src, ozSpawn, firstTick int, c4HealFunc func()) func() {
	return func() {
		if src != c.burstOzSpawnSrc {
			return
		}
		c.burstOzSpawnSrc = -1
		c.queueOz("Burst", ozSpawn, firstTick)
		// C4 heal should happen right after oz spawn/end of animation because buffs proc'd from the heal are not snapped into Oz
		if c4HealFunc != nil {
			c.Core.Tasks.Add(c4HealFunc, ozSpawn+1)
		}
	}
}
