package gorou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 31 // Initial Hit

func init() {
	burstFrames = frames.InitAbilSlice(56) // Q -> E
	burstFrames[action.ActionAttack] = 53  // Q -> N1
	burstFrames[action.ActionDash] = 42    // Q -> D
	burstFrames[action.ActionJump] = 43    // Q -> J
	burstFrames[action.ActionSwap] = 55    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// Initial Hit
	// A1/C6/Q duration all start on Initial Hit
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Juuga: Forward Unto Victory",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       burst[c.TalentLvlBurst()],
		}
		// A4 Part 2
		// Juuga: Forward Unto Victory: Skill DMG and Crystal Collapse DMG increased by 15.6% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[attributes.DEFP] + snap.Stats[attributes.DEF]) * 0.156

		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5),
			0,
		)

		// Q General's Glory:
		// Like the General's War Banner created by Inuzaka All-Round Defense, provides buffs to active characters
		// within the skill's AoE based on the number of Geo characters in the party. Also moves together with
		// your active character.
		c.eFieldSrc = c.Core.F
		c.Core.Tasks.Add(c.gorouSkillBuffField(c.Core.F), 17) // 17 so we get one last tick

		// If a General's War Banner created by Gorou currently exists on the field when this ability is used,
		// it will be destroyed. In addition, for the duration of General's Glory, Gorou's
		// Elemental Skill "Inuzaka All-Round Defense" will not create the General's War Banner.
		c.Core.Status.Delete(generalWarBannerKey)
		c.Core.Status.Add(generalGloryKey, generalGloryDuration) // field starts on Hitmark Initial

		// Generates 1 Crystal Collapse every 1.5s that deals AoE Geo DMG to 1 opponent within the skill's AoE.
		// Pulls 1 elemental shard in the skill's AoE to your active character's position every 1.5s (elemental
		// shards are created by Crystallize reactions).
		c.qFieldSrc = c.Core.F
		c.Core.Tasks.Add(c.gorouCrystalCollapse(c.Core.F), 90) // first crystal collapse is 1.5s after Hitmark Initial

		// A1: After using Juuga: Forward Unto Victory, all nearby party members' DEF is increased by 25% for 12s.
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(a1Key, 720),
				AffectedStat: attributes.DEFP,
				Amount: func() ([]float64, bool) {
					return c.a1Buff, true
				},
			})
		}

		// C4
		if c.Base.Cons >= 4 && c.geoCharCount > 1 {
			// TODO: not sure if this actually snapshots stats
			// ai := combat.AttackInfo{
			// 	Abil:      "Inuzaka All-Round Defense C4",
			// 	AttackTag: combat.AttackTagNone,
			// }
			c.healFieldStats, _ = c.Stats()
			c.Core.Tasks.Add(c.gorouBurstHealField(c.Core.F), 90)
		}

		// C6
		if c.Base.Cons >= 6 {
			c.c6()
		}
	}, burstHitmark)

	//TODO:  If Gorou falls, the effects of General's Glory will be cleared.

	c.c2Extension = 0

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

// recursive function for dealing damage
func (c *char) gorouCrystalCollapse(src int) func() {
	return func() {
		//do nothing if this has been overwritten
		if c.qFieldSrc != src {
			return
		}
		//do nothing if field expired
		if c.Core.Status.Duration(generalGloryKey) == 0 {
			return
		}
		//trigger damage
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Crystal Collapse",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       burstTick[c.TalentLvlBurst()],
		}
		//Juuga: Forward Unto Victory: Skill DMG and Crystal Collapse DMG increased by 15.6% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[attributes.DEFP] + snap.Stats[attributes.DEF]) * 0.156

		enemy := c.Core.Combat.ClosestEnemyWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 8), nil)
		if enemy != nil {
			c.Core.QueueAttackWithSnap(
				ai,
				snap,
				combat.NewCircleHitOnTarget(enemy, nil, 3.5),
				//TODO: skill damage frames
				1,
			)
		}

		//tick every 1.5s
		c.Core.Tasks.Add(c.gorouCrystalCollapse(src), 90)
	}
}

func (c *char) gorouBurstHealField(src int) func() {
	return func() {
		//do nothing if this has been overwritten
		if c.qFieldHealSrc != src {
			return
		}
		//do nothing if field expired
		if c.Core.Status.Duration(generalGloryKey) == 0 {
			return
		}
		//When General's Glory is in the "Impregnable" or "Crunch" states, it will also heal active characters
		//within its AoE by 50% of Gorou's own DEF every 1.5s.
		amt := c.Base.Def*(1+c.healFieldStats[attributes.DEFP]) + c.healFieldStats[attributes.DEF]
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Lapping Hound: Warm as Water",
			Src:     0.5 * amt,
			Bonus:   c.Stat(attributes.Heal),
		})

		//tick every 1.5s
		c.Core.Tasks.Add(c.gorouBurstHealField(src), 90)
	}
}
