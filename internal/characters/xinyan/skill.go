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

const skillHitmark = 15
const skillShieldStart = 28

func init() {
	skillFrames = frames.InitAbilSlice(62) // E -> Swap
	skillFrames[action.ActionAttack] = 53  // E -> N1
	skillFrames[action.ActionBurst] = 54   // E -> Q
	skillFrames[action.ActionDash] = 54    // E -> D
	skillFrames[action.ActionJump] = 53    // E -> J
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Sweeping Fervor",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.09 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	snap := c.Snapshot(&ai)

	defFactor := snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]

	hitOpponents := 0
	cb := func(_ combat.AttackCB) {
		hitOpponents++

		// including a1
		c.QueueCharTask(func() {
			if hitOpponents >= c.shieldLevel3Requirement && c.shieldLevel < 3 {
				c.updateShield(3, defFactor)
			} else if hitOpponents >= c.shieldLevel2Requirement && c.shieldLevel < 2 {
				c.updateShield(2, defFactor)
			}
		}, skillShieldStart-skillHitmark)
	}

	if c.Core.Player.Shields.Get(shield.ShieldXinyanSkill) == nil {
		c.QueueCharTask(func() {
			c.updateShield(1, defFactor)
		}, skillShieldStart)
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1}, 3),
		skillHitmark,
		skillHitmark,
		cb,
		c.c4,
	)

	c.SetCDWithDelay(action.ActionSkill, 18*60, 13)
	c.Core.QueueParticle("xinyan", 4, attributes.Pyro, skillHitmark+c.ParticleDelay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) shieldDot(src int) func() {
	return func() {
		if c.Core.Player.Shields.Get(shield.ShieldXinyanSkill) == nil {
			return
		}
		if c.shieldLevel != 3 {
			return
		}
		if c.shieldTickSrc != src {
			return
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Sweeping Fervor (DoT)",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       skillDot[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3), 1, 1)

		c.Core.Tasks.Add(c.shieldDot(src), 2*60)
	}
}

func (c *char) updateShield(level int, defFactor float64) {
	previousLevel := c.shieldLevel
	c.shieldLevel = level
	shieldhp := 0.0

	switch c.shieldLevel {
	case 1:
		shieldhp = shieldLevel1Flat[c.TalentLvlSkill()] + shieldLevel1[c.TalentLvlSkill()]*defFactor
	case 2:
		shieldhp = shieldLevel2Flat[c.TalentLvlSkill()] + shieldLevel2[c.TalentLvlSkill()]*defFactor
	case 3:
		shieldhp = shieldLevel3Flat[c.TalentLvlSkill()] + shieldLevel3[c.TalentLvlSkill()]*defFactor
		c.shieldTickSrc = c.Core.F
		c.Core.Tasks.Add(c.shieldDot(c.Core.F), 2*60)
	}
	shd := c.newShield(shieldhp, shield.ShieldXinyanSkill, 12*60)
	c.Core.Player.Shields.Add(shd)
	c.Core.Log.NewEvent("update shield level", glog.LogCharacterEvent, c.Index).
		Write("previousLevel", previousLevel).
		Write("level", c.shieldLevel).
		Write("expiry", shd.Expiry())
}
