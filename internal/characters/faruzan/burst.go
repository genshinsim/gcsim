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
	burstHitmark  = 54
	burstBuffKey  = "faruzan-q-dmg-bonus"
	burstShredKey = "faruzan-q-shred"
)

func init() {
	burstFrames = frames.InitAbilSlice(71)
	burstFrames[action.ActionAttack] = 61 // Q -> N1
	burstFrames[action.ActionAim] = 61    // Q -> Aim
	burstFrames[action.ActionSkill] = 61  // Q -> E
	burstFrames[action.ActionDash] = 61   // Q -> D
	burstFrames[action.ActionJump] = 63   // Q -> J
	burstFrames[action.ActionSwap] = 60   // Q -> Swap
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

	c.Core.Tasks.Add(func() {
		snap = c.Snapshot(&ai)
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 6.3), 0, applyBurstShred)
		for _, char := range c.Core.Player.Chars() {
			c.applyBurstBuff(char)
		}
	}, burstHitmark) // initial hit

	// C2: The duration of the Dazzling Polyhedron created by
	// The Wind's Secret Ways increased by 6s.
	duration := 745
	if c.Base.Cons >= 2 {
		duration += 360
	}

	frequency, ok := p["frequency"]
	if !ok {
		frequency = 1
	}
	if frequency < 1 {
		frequency = 1
	}
	if frequency > 3 {
		frequency = 3
	}

	// following hits
	whirl_ai := ai
	whirl_ai.Abil = "Whirlwind Pulse (Q)"
	whirl_ai.Mult = 0 // is this a 0 damage hit?
	whirl_ai.Element = attributes.NoElement
	hitCount := 0
	for i := 71; i <= duration; i += 120 {
		c.Core.Tasks.Add(func() {
			for _, char := range c.Core.Player.Chars() {
				c.applyBurstBuff(char)
			}
			hitCount++
			if hitCount%3 >= frequency {
				return
			}
			c.Core.QueueAttackWithSnap(whirl_ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 5), 0, applyBurstShred)
		}, burstHitmark+i)
	}

	c.SetCD(action.ActionBurst, 1200)
	c.ConsumeEnergy(3)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) applyBurstBuff(char *character.CharWrapper) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.AnemoP] = burstBuff[c.TalentLvlBurst()]
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(burstBuffKey, 240),
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	if c.Base.Cons >= 6 {
		c.c6Buff(char)
	}
}

func applyBurstShred(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	t.AddResistMod(enemy.ResistMod{
		Base:  modifier.NewBaseWithHitlag(burstShredKey, 240),
		Ele:   attributes.Anemo,
		Value: -0.3,
	})
}
