package lisa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var skillPressFrames []int
var skillHoldFrames []int

const (
	skillPressAnimation = 40
	skillPressHitmark   = 22
	skillHoldAnimation  = 143
	skillHoldHitmark    = 117
)

func init() {
	skillPressFrames = frames.InitAbilSlice(skillPressAnimation)

	skillPressFrames[action.ActionAttack] = 37
	skillPressFrames[action.ActionCharge] = 38
	skillPressFrames[action.ActionBurst] = 40
	skillPressFrames[action.ActionDash] = 35
	skillPressFrames[action.ActionJump] = 20
	skillPressFrames[action.ActionSwap] = 23

	skillHoldFrames = frames.InitAbilSlice(141)
	skillHoldFrames[action.ActionAttack] = 143
	skillHoldFrames[action.ActionCharge] = 125
	skillHoldFrames[action.ActionBurst] = 138
	skillHoldFrames[action.ActionDash] = 116
	skillHoldFrames[action.ActionJump] = 117
}

//p = 0 for no hold, p = 1 for hold
func (c *char) Skill(p map[string]int) action.ActionInfo {
	hold := p["hold"]
	if hold == 1 {
		return c.skillHold(p)
	}
	return c.skillPress(p)
}

//TODO: how long do stacks last?
func (c *char) skillPress(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Violet Arc",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagLisaElectro,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	cb := func(a combat.AttackCB) {
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		count := t.GetTag(conductiveTag)
		if count < 3 {
			t.SetTag(conductiveTag, count+1)
		}
	}

	c.Core.QueueAttack(ai,
		combat.NewDefSingleTarget(1, combat.TargettableEnemy),
		0,
		skillPressHitmark,
		cb,
	)

	if c.Core.Rand.Float64() < 0.5 {
		c.Core.QueueParticle("Lisa", 1, attributes.Electro, skillPressHitmark+100)
	}

	c.SetCDWithDelay(action.ActionSkill, 60, 17)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   20, //fastest cancel is at 20
		State:           action.SkillState,
	}
}

//After an extended casting time, calls down lightning from the heavens, dealing massive Electro DMG to all nearby opponents.
//Deals great amounts of extra damage to opponents based on the number of Conductive stacks applied to them, and clears their Conductive status.
func (c *char) skillHold(p map[string]int) action.ActionInfo {
	//no multiplier as that's target dependent
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Violet Arc (Hold)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 50,
	}

	//c2 add defense? no interruptions either way
	if c.Base.Cons >= 2 {
		//increase def for the duration of this abil in however many frames
		val := make([]float64, attributes.EndStatType)
		val[attributes.DEFP] = 0.25
		c.AddStatMod("lisa-c2", 126, attributes.NoStat,
			func() ([]float64, bool) { return val, true },
		)
	}

	clearStacks := func(a combat.AttackCB) {
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		//clear stacks
		t.SetTag(conductiveTag, 0)
	}

	count := 0
	var c1cb func(a combat.AttackCB)
	if c.Base.Cons > 0 {
		c1cb = func(a combat.AttackCB) {
			if count == 5 {
				return
			}
			count++
			c.AddEnergy("lisa-c1", 2)
		}
	}

	//[8:31 PM] ArchedNosi | Lisa Unleashed: yeah 4-5 50/50 with Hold
	//[9:13 PM] ArchedNosi | Lisa Unleashed: @gimmeabreak actually wait, xd i noticed i misread my sheet, Lisa Hold E always gens 5 orbs
	c.Core.QueueAttack(ai,
		combat.NewDefCircHit(3, false, combat.TargettableEnemy),
		0,
		skillHoldHitmark,
		clearStacks,
		c1cb,
	)

	// count := 4
	// if c.Core.Rand.Float64() < 0.5 {
	// 	count = 5
	// }
	c.Core.QueueParticle("Lisa", 5, attributes.Electro, skillHoldHitmark+100)

	// c.CD[def.SkillCD] = c.Core.F + 960 //16seconds, starts after 114 frames
	c.SetCDWithDelay(action.ActionSkill, 960, 114)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   20, //fastest cancel is at 20
		State:           action.SkillState,
	}
}
