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
		if hitOpponents >= 2 && c.shieldLevel < 3 {
			c.updateShield(3, defFactor)
		} else if hitOpponents >= 1 && c.shieldLevel < 2 {
			c.updateShield(2, defFactor)
		}
	}

	if c.Core.Player.Shields.Get(shield.ShieldXinyanSkill) == nil {
		c.Core.Tasks.Add(func() {
			c.updateShield(1, defFactor)
		}, skillHitmark)
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), skillHitmark, skillHitmark, cb, c.c4)

	c.SetCDWithDelay(action.ActionSkill, 18*60, 6)
	c.Core.QueueParticle("xinyan", 4, attributes.Pyro, skillHitmark+80)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
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

func (c *char) updateShield(level int, defFactor float64) {
	c.shieldLevel = level
	shieldhp := 0.0

	switch c.shieldLevel {
	case 1:
		shieldhp = shieldLevel1Flat[c.TalentLvlSkill()] + shieldLevel1[c.TalentLvlSkill()]*defFactor
	case 2:
		shieldhp = shieldLevel2Flat[c.TalentLvlSkill()] + shieldLevel2[c.TalentLvlSkill()]*defFactor
	case 3:
		shieldhp = shieldLevel3Flat[c.TalentLvlSkill()] + shieldLevel3[c.TalentLvlSkill()]*defFactor
		c.Core.Tasks.Add(c.shieldDot(), 2*60)
	}
	shd := c.newShield(shieldhp, shield.ShieldXinyanSkill, 12*60)
	c.Core.Player.Shields.Add(shd)
	c.Core.Log.NewEvent("update shield level", glog.LogCharacterEvent, c.Index).
		Write("level", c.shieldLevel).
		Write("expiry", shd.Expiry())
}
