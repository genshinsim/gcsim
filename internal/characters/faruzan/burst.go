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

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 6.3),
		burstHitmark,
		burstHitmark,
		applyBurstShredCb,
	)

	// C2: The duration of the Dazzling Polyhedron created by
	// The Wind's Secret Ways increased by 6s.
	duration := 745
	if c.Base.Cons >= 2 {
		duration += 360
	}

	c.burstSrc = c.Core.F
	currSrc := c.burstSrc

	x, y := c.Core.Combat.Player().Pos()
	count := 0
	for i := 137; i <= duration; i += 120 {
		ox, oy := calcGadgetOffsets(count)
		c.Core.Tasks.Add(func() {
			if c.burstSrc != currSrc {
				return
			}
			for id := range c.Core.Combat.EnemiesWithinRadius(x+ox, y+oy, 6) {
				trg, ok := c.Core.Combat.Enemy(id).(*enemy.Enemy)
				if !ok {
					continue
				}
				applyBurstShred(trg)
			}
		}, 43+i)
		count += 1
	}

	field := combat.NewCircleHit(c.Core.Combat.Player(), 40)
	for i := 0; i <= duration; i += 6 {
		c.Core.Tasks.Add(func() {
			if c.burstSrc != currSrc {
				return
			}
			if !combat.WillCollide(field, c.Core.Combat.Player(), 0) {
				return
			}
			for _, char := range c.Core.Player.Chars() {
				c.applyBurstBuff(char)
			}
		}, 43+i)
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

func applyBurstShredCb(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	applyBurstShred(t)
}

func applyBurstShred(trg *enemy.Enemy) {
	trg.AddResistMod(enemy.ResistMod{
		Base:  modifier.NewBaseWithHitlag(burstShredKey, 240),
		Ele:   attributes.Anemo,
		Value: -0.3,
	})
}

func calcGadgetOffsets(iter int) (float64, float64) {
	x := 0.0
	switch iter % 3 {
	case 1:
		x = 5.19
	case 2:
		x = -5.19
	}
	y := 1.5
	switch iter % 3 {
	case 1, 2:
		y = 10.5
	}
	return x, y
}
