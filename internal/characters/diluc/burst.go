package diluc

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 100

func init() {
	burstFrames = frames.InitAbilSlice(140) // Q -> D
	burstFrames[action.ActionAttack] = 139  // Q -> N1
	burstFrames[action.ActionSkill] = 139   // Q -> E
	burstFrames[action.ActionJump] = 139    // Q -> J
	burstFrames[action.ActionSwap] = 138    // Q -> Swap
}

const burstBuffKey = "diluc-q"

func (c *char) Burst(p map[string]int) action.Info {
	// A4:
	// The Pyro Infusion provided by Dawn lasts for 4s longer.
	duration := 480
	hasA4 := c.Base.Ascension >= 4
	if hasA4 {
		duration += 240
	}

	// infusion starts when burst starts and ends when burst comes off CD - check any diluc video
	c.AddStatus(burstBuffKey, duration, true)

	// A4:
	// Additionally, Diluc gains 20% Pyro DMG Bonus during the duration of this effect.
	if hasA4 {
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(burstBuffKey, duration),
			AffectedStat: attributes.PyroP,
			Amount: func() ([]float64, bool) {
				return c.a4buff, true
			},
		})
	}

	// Snapshot occurs late in the animation when it is released from the claymore
	// For our purposes, snapshot upon damage proc
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Dawn (Strike)",
			AttackTag:          attacks.AttackTagElementalBurst,
			ICDTag:             attacks.ICDTagElementalBurst,
			ICDGroup:           attacks.ICDGroupDiluc,
			StrikeType:         attacks.StrikeTypeBlunt,
			Element:            attributes.Pyro,
			Durability:         50,
			Mult:               burstInitial[c.TalentLvlBurst()],
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.09 * 60,
			CanBeDefenseHalted: true,
		}

		c.Core.QueueAttack(
			ai,
			combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1}, 16, 6),
			0,
			1,
		)

		ai.StrikeType = attacks.StrikeTypeDefault
		// both initial hit, DoT and explosion all have 50 durability
		ai.Abil = "Dawn (Tick)"
		ai.Mult = burstDOT[c.TalentLvlBurst()]

		// only initial hit has hitlag
		ai.HitlagHaltFrames = 0
		ai.CanBeDefenseHalted = false

		// DoT and explosion dmg
		// - gadget spawns at Y: 1m and lives for ~1.7s until it explodes
		// - moves at 14 m/s with dmg happening every 0.2s, so it moves at 2.8m per attack
		// - 1.7s / (0.2 s/attack) ~= 8 attacks total before explosion
		initialPos := c.Core.Combat.Player().Pos()
		initialDirection := c.Core.Combat.Player().Direction()
		for i := 0; i < 8; i++ {
			nextPos := geometry.CalcOffsetPoint(initialPos, geometry.Point{Y: 1 + 2.8*float64(i)}, initialDirection)
			c.Core.QueueAttack(
				ai,
				combat.NewBoxHit(c.Core.Combat.Player(), nextPos, geometry.Point{Y: -5}, 16, 8),
				0,
				(i+1)*12,
			)
		}

		ai.Abil = "Dawn (Explode)"
		ai.Mult = burstExplode[c.TalentLvlBurst()]
		// 1m + 14 m/s * 1.7s
		finalPos := geometry.CalcOffsetPoint(initialPos, geometry.Point{Y: 1 + 14*1.7}, initialDirection)
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(c.Core.Combat.Player(), finalPos, geometry.Point{Y: -6}, 16, 10),
			0,
			1.7*60,
		)
	}, burstHitmark)

	c.ConsumeEnergy(21)
	c.SetCDWithDelay(action.ActionBurst, 720, 14)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
