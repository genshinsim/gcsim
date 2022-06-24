package sayu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	c.updateSkillFrames(hold)
	var cd, delay int
	if hold > 0 {
		if hold > 600 { // 10s
			hold = 600
		}

		// 18 = 15 anim start + 3 to start swirling
		// +2 frames for not proc the sacrificial by "Yoohoo Art: Fuuin Dash (Elemental DMG)"
		delay = 18 + hold + 2
		c.skillHold(p, hold)
		cd = int(6*60 + float64(hold)*0.5)
	} else {
		delay = 15
		c.skillPress(p)
		cd = 6 * 60
	}

	c.SetCDWithDelay(action.ActionSkill, cd, delay)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(c.skillFrames),
		AnimationLength: c.skillFrames[action.InvalidAction],
		CanQueueAfter:   c.skillFrames[action.InvalidAction],
		State:           action.SkillState,
	}
}

func (c *char) skillPress(p map[string]int) {

	c.c2Bonus = 0.033

	// Fuufuu Windwheel DMG
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Press)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 3)

	// Fuufuu Whirlwind Kick Press DMG
	ai = combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Press)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillPressEnd[c.TalentLvlSkill()],
	}
	snap = c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), 3+25)

	c.Core.QueueParticle("sayu-skill", 2, attributes.Anemo, c.skillFrames[action.InvalidAction]+73)

}

func (c *char) skillHold(p map[string]int, duration int) {

	c.eInfused = attributes.NoElement
	c.eDuration = c.Core.F + 18 + duration + 20
	c.infuseCheckLocation = combat.NewDefCircHit(0.1, true, combat.TargettablePlayer, combat.TargettableEnemy, combat.TargettableObject)
	c.c2Bonus = .0

	// ticks
	i := 0
	d := c.createSkillHoldSnapshot()
	for ; i <= duration; i += 30 { // 1 tick for sure
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackEvent(d, 0)

			if c.Base.Cons >= 2 && c.c2Bonus < 0.66 {
				c.c2Bonus += 0.033
				c.Core.Log.NewEvent("sayu c2 adding 3.3% dmg", glog.LogCharacterEvent, c.Index, "dmg bonus%", c.c2Bonus)
			}
		}, 18+i)

		if i%180 == 0 { // 3s
			c.Core.QueueParticle("sayu-skill-hold", 1, attributes.Anemo, 18+i+73)
		}
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Hold)",
		AttackTag:  combat.AttackTagElementalArtHold,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillHoldEnd[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), 18+duration+20)

	c.Core.QueueParticle("sayu-skill", 2, attributes.Anemo, c.skillFrames[action.InvalidAction]+73)
}

func (c *char) createSkillHoldSnapshot() *combat.AttackEvent {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Hold Tick)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	return (&combat.AttackEvent{
		Info:        ai,
		Pattern:     combat.NewDefCircHit(0.5, false, combat.TargettableEnemy),
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	})
}
