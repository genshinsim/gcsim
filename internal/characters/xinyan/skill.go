package xinyan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var skillFrames []int

const skillHitmark = 65

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sweeping Fervor",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	defFactor := snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]

	hitOpponents := 0
	cb := func(a combat.AttackCB) {
		hitOpponents++

		// including a1
		updateShield := false
		shieldhp := 0.0
		if hitOpponents >= 2 && c.shieldLevel < 3 {
			shieldhp = shieldLevel3Flat[c.TalentLvlSkill()] + shieldLevel3[c.TalentLvlSkill()]*defFactor
			c.shieldLevel = 3
			updateShield = true
			c.Core.Tasks.Add(c.shieldDot(), 2*60)
		} else if hitOpponents >= 1 && c.shieldLevel < 2 {
			shieldhp = shieldLevel2Flat[c.TalentLvlSkill()] + shieldLevel2[c.TalentLvlSkill()]*defFactor
			c.shieldLevel = 2
			updateShield = true
		}

		if updateShield {
			shd := c.newShield(shieldhp, shield.ShieldXinyanSkill, 12*60)
			c.Core.Player.Shields.Add(shd)
			c.Core.Log.NewEvent("update shield level", glog.LogCharacterEvent, c.Index, "level", c.shieldLevel, "expiry", shd.Expiry())
		}
	}

	if c.Core.Player.Shields.Get(shield.ShieldXinyanSkill) == nil {
		shieldhp := shieldLevel1Flat[c.TalentLvlSkill()] + shieldLevel1[c.TalentLvlSkill()]*defFactor
		c.Core.Player.Shields.Add(c.newShield(shieldhp, shield.ShieldXinyanSkill, 12*60))
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), skillHitmark, skillHitmark, cb)

	c.SetCDWithDelay(action.ActionSkill, 18*60, 6)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		Post:            skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) shieldDot() func() {
	return func() {
		if c.Core.Player.Shields.Get(shield.ShieldXinyanSkill) == nil {
			return
		}
		if c.shieldLevel != 3 {
			return
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Sweeping Fervor (DoT)",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       skillDot[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 1, 1)

		c.Core.Tasks.Add(c.shieldDot(), 2*60)
	}
}
