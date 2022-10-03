package amber

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstStart = 72 // hitmark of the first tick

func init() {
	burstFrames = frames.InitAbilSlice(111) // Q -> N1/E
	burstFrames[action.ActionDash] = 57     // Q -> D
	burstFrames[action.ActionJump] = 58     // Q -> J
	burstFrames[action.ActionWalk] = 62     // Q -> Walk
	burstFrames[action.ActionSwap] = 60     // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       "Fiery Rain",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupAmber,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	//2sec duration, tick every .4 sec in zone 1
	//2sec duration, tick every .6 sec in zone 2
	//2sec duration, tick every .2 sec in zone 3

	//TODO: properly implement random hits and hit box range. right now everything is just radius 3
	for i := 24; i < 120; i += 24 {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy, combat.TargettableGadget), burstStart+i)
	}
	for i := 36; i < 120; i += 36 {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy, combat.TargettableGadget), burstStart+i)
	}
	for i := 12; i < 120; i += 12 {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 3, false, combat.TargettableEnemy, combat.TargettableGadget), burstStart+i)
	}

	if c.Base.Cons >= 6 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.15
		for _, active := range c.Core.Player.Chars() {
			active.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("amber-c6", 900),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
	}

	c.SetCDWithDelay(action.ActionBurst, 720, 56)
	c.ConsumeEnergy(59)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
