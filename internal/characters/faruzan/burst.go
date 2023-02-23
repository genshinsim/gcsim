package faruzan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
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
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.5}, 6.3),
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

	player := c.Core.Combat.Player()
	playerPos := player.Pos()
	direction := player.Direction()
	gadgetPositions := []combat.Point{
		combat.CalcOffsetPoint(playerPos, combat.Point{X: 5.19, Y: 10.5}, direction),
		combat.CalcOffsetPoint(playerPos, combat.Point{X: -5.19, Y: 10.5}, direction),
		combat.CalcOffsetPoint(playerPos, combat.Point{Y: 1.5}, direction),
	}
	count := 0
	for i := 137; i <= duration; i += 120 {
		gadgetPos := gadgetPositions[count%3].Pos()
		c.Core.Tasks.Add(func() {
			if c.burstSrc != currSrc {
				return
			}
			enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(gadgetPos, nil, 6), nil)
			for _, e := range enemies {
				applyBurstShred(e)
			}
		}, 43+i)
		count += 1
	}

	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 40)
	buffFunc := func() {
		if c.burstSrc != currSrc {
			return
		}
		if !c.Core.Combat.Player().IsWithinArea(burstArea) {
			return
		}
		for _, char := range c.Core.Player.Chars() {
			c.applyBurstBuff(char)
		}
	}

	// In-game refreshes 0.1s. We give buff every 239f to reduce spam.
	for i := 0; i <= duration; i += 239 {
		c.Core.Tasks.Add(buffFunc, 43+i)
	}

	// Last refresh to account for 0.1s tick period
	c.Core.Tasks.Add(buffFunc, 43+duration)

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
	applyBurstShred(a.Target)
}

func applyBurstShred(trg combat.Target) {
	t, ok := trg.(*enemy.Enemy)
	if !ok {
		return
	}
	t.AddResistMod(combat.ResistMod{
		Base:  modifier.NewBaseWithHitlag(burstShredKey, 240),
		Ele:   attributes.Anemo,
		Value: -0.3,
	})
}
