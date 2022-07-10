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

const burstStart = 74

func init() {
	burstFrames = frames.InitAbilSlice(74)
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
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(3, false, combat.TargettableEnemy), burstStart+i)
	}
	for i := 36; i < 120; i += 36 {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(3, false, combat.TargettableEnemy), burstStart+i)
	}
	for i := 12; i < 120; i += 12 {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(3, false, combat.TargettableEnemy), burstStart+i)
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

	c.ConsumeEnergy(64)
	c.SetCDWithDelay(action.ActionBurst, 720, 64)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstStart,
		State:           action.BurstState,
	}
}
