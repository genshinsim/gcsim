package gorou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 74

func init() {
	burstFrames = frames.InitAbilSlice(74)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c.Core.Tasks.Add(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Juuga: Forward Unto Victory",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Geo,
			StrikeType: combat.StrikeTypeBlunt,
			//TODO: don't know the gauge of this
			Durability: 25,
			Mult:       burst[c.TalentLvlSkill()],
		}
		//Juuga: Forward Unto Victory: Skill DMG and Crystal Collapse DMG increased by 15.6% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[attributes.DEFP] + snap.Stats[attributes.DEF]) * 0.156

		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewDefCircHit(5, false, combat.TargettableEnemy),
			//TODO: skill damage frames
			0,
		)
	}, burstHitmark+10)

	//Like the General's War Banner created by Inuzaka All-Round Defense, provides buffs to active characters
	//within the skill's AoE based on the number of Geo characters in the party. Also moves together with
	//your active character.
	c.eFieldSrc = c.Core.F
	c.Core.Tasks.Add(c.gorouSkillBuffField(c.Core.F), 59) //59 so we get one last tick

	//If a General's War Banner created by Gorou currently exists on the field when this ability is used,
	//it will be destroyed. In addition, for the duration of General's Glory, Gorou's
	//Elemental Skill "Inuzaka All-Round Defense" will not create the General's War Banner.
	c.Core.Status.Delete(generalWarBannerKey)
	c.Core.Status.Add(generalGloryKey, generalGloryDuration)

	//Generates 1 Crystal Collapse every 1.5s that deals AoE Geo DMG to 1 opponent within the skill's AoE.
	//Pulls 1 elemental shard in the skill's AoE to your active character's position every 1.5s (elemental
	//shards are created by Crystallize reactions).
	c.qFieldSrc = c.Core.F
	c.Core.Tasks.Add(c.gorouCrystalCollapse(c.Core.F), 90) //every 90s?

	//TODO:  If Gorou falls, the effects of General's Glory will be cleared.

	//A1: After using Juuga: Forward Unto Victory, all nearby party members' DEF is increased by 25% for 12s.
	m := make([]float64, attributes.EndStatType)
	m[attributes.DEFP] = .25
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(heedlessKey, 720),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	//c6 check
	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.c2Extension = 0

	c.SetCDWithDelay(action.ActionBurst, 20*60, 8)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}

//recursive function for dealing damage
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
			Element:    attributes.Geo,
			StrikeType: combat.StrikeTypeBlunt,
			//TODO: don't know the gauge of this
			Durability: 25,
			Mult:       burstTick[c.TalentLvlSkill()],
		}
		//Juuga: Forward Unto Victory: Skill DMG and Crystal Collapse DMG increased by 15.6% of DEF
		snap := c.Snapshot(&ai)
		ai.FlatDmg = (snap.BaseDef*snap.Stats[attributes.DEFP] + snap.Stats[attributes.DEF]) * 0.156

		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewDefCircHit(5, false, combat.TargettableEnemy),
			//TODO: skill damage frames
			1,
		)

		//tick every 1.5s
		c.Core.Tasks.Add(c.gorouCrystalCollapse(src), 90)
	}
}
