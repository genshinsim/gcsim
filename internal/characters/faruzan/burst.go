package faruzan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstBuffKey  = "faruzan-q-dmg-bonus"
	burstShredKey = "faruzan-q-shred"
)

func init() {
	burstFrames = frames.InitAbilSlice(55)
}

// Faruzan deploys a Dazzling Polyhedron that deals AoE Anemo DMG and releases
// a Whirlwind Pulse. While the Dazzling Polyhedron persists, it will
// continuously move along a triangular path. Once it reaches each corner of
// that triangular path, it will unleash 1 more Whirlwind Pulse.
//
// Whirlwind Pulse
// - When the Whirlwind Pulse hits opponents, it will apply Perfidious Wind's
// Ruin to them, decreasing their Anemo RES.
// - The Whirlwind Pulse will also apply Prayerful Wind's Gift to all nearby
// characters when it is unleashed, granting them Anemo DMG Bonus.
func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "The Wind's Secret Ways (Q)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	shredCb := func(a combat.AttackCB) {
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		t.AddResistMod(enemy.ResistMod{
			Base:  modifier.NewBaseWithHitlag(burstShredKey, 240),
			Ele:   attributes.Anemo,
			Value: -0.4,
		})
	}

	// C2: The duration of the Dazzling Polyhedron created by
	// The Wind's Secret Ways increased by 6s.
	duration := 720
	if c.Base.Cons >= 2 {
		duration += 360
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.AnemoP] = burstBuff[c.TalentLvlBurst()]
	for i := 55; i <= duration+55; i += 120 {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 5), 0, shredCb)
			for _, char := range c.Core.Player.Chars() {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(burstBuffKey, 240),
					AffectedStat: attributes.CR,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}
		}, i)
	}

	c.SetCDWithDelay(action.ActionBurst, 1200, 18)
	c.ConsumeEnergy(21)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
