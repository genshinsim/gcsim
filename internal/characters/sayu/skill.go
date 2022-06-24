package sayu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillPressFrames []int
var skillHoldFrames []int

const skillPressHitmark = 41
const skillHoldHitmark = 79

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(skillPressHitmark)

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(skillHoldHitmark)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	if hold > 0 {
		if hold > 600 { // 10s
			hold = 600
		}
		return c.skillHold(p, hold)
	}
	return c.skillPress(p)
}

func (c *char) skillPress(p map[string]int) action.ActionInfo {

	c.c2Bonus = 0.033

	// Fuufuu Windwheel DMG
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Press)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagSayuSkillAnemo,
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
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), 28)

	c.Core.QueueParticle("sayu-skill", 2, attributes.Anemo, skillPressHitmark+73)

	c.SetCDWithDelay(action.ActionSkill, 6*60, 15)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.InvalidAction],
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int, duration int) action.ActionInfo {

	c.eInfused = attributes.NoElement
	c.eInfusedTag = combat.ICDTagNone
	c.eDuration = c.Core.F + 18 + duration + 20
	c.infuseCheckLocation = combat.NewDefCircHit(0.1, true, combat.TargettablePlayer, combat.TargettableEnemy, combat.TargettableObject)
	c.c2Bonus = .0

	// ticks
	d := c.createSkillHoldSnapshot()
	c.Core.Tasks.Add(c.absorbCheck(c.Core.F, 0, int(duration/12)), 18)

	for i := 0; i <= duration; i += 30 { // 1 tick for sure
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

	c.Core.QueueParticle("sayu-skill", 2, attributes.Anemo, skillHoldHitmark+73)

	// 18 = 15 anim start + 3 to start swirling
	// +2 frames for not proc the sacrificial by "Yoohoo Art: Fuuin Dash (Elemental DMG)"
	c.SetCDWithDelay(action.ActionSkill, int(6*60+float64(duration)*0.5), 18+duration+2)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return skillHoldFrames[next] + duration },
		AnimationLength: skillHoldFrames[action.InvalidAction] + duration,
		CanQueueAfter:   skillHoldFrames[action.InvalidAction] + duration,
		State:           action.SkillState,
	}
}

func (c *char) createSkillHoldSnapshot() *combat.AttackEvent {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Fuuin Dash (Hold Tick)",
		AttackTag:  combat.AttackTagElementalArtHold,
		ICDTag:     combat.ICDTagSayuSkillAnemo,
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

func (c *char) absorbCheck(src, count, max int) func() {
	return func() {
		if count == max {
			return
		}

		c.eInfused = c.Core.Combat.AbsorbCheck(c.infuseCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)
		if c.eInfused != attributes.NoElement {
			switch c.eInfused {
			case attributes.Pyro:
				c.eInfusedTag = combat.ICDTagSayuSkillPyro
			case attributes.Hydro:
				c.eInfusedTag = combat.ICDTagSayuSkillHydro
			case attributes.Electro:
				c.eInfusedTag = combat.ICDTagSayuSkillElectro
			case attributes.Cryo:
				c.eInfusedTag = combat.ICDTagSayuSkillCryo
			}
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index,
				"sayu infused ", c.eInfused.String(),
			)
			return
		}
		c.Core.Tasks.Add(c.absorbCheck(src, count+1, max), 12)
	}
}
