package barbara

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// barbara skill - copied from bennett burst

func (c *char) Skill(p map[string]int) action.ActionInfo {

	f, a := c.ActionFrames(action.ActionSkill, p)

	//add field effect timer
	//assumes a4
	c.Core.Status.AddStatus("barbskill", 15*60+1)
	//hook for buffs; active right away after cast

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Let the Show Begin♪",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25, //TODO: what is 1A GU?
		Mult:       skill[c.TalentLvlSkill()],
	}
	//TODO: review barbara AOE size?
	c.Core.Combat.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 5, 5)
	c.Core.Combat.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 5, 35) // need to confirm timing of this

	stats, _ := c.SnapshotStats()
	hpplus := stats[attributes.Heal]
	heal := skillhp[c.TalentLvlSkill()] + skillhpp[c.TalentLvlSkill()]*c.MaxHP()
	//apply right away

	c.skillInitF = c.Core.F
	c.onSkillStackCount(c.Core.F)
	//add 1 tick each 5s
	//first tick starts at 0
	c.barbaraHealTick(heal, hpplus, c.Core.F)()
	ai.Abil = "Let the Show Begin♪ Wet Tick"
	ai.AttackTag = combat.AttackTagNone
	ai.Mult = 0
	c.barbaraWet(ai, c.Core.F)()
	if c.Base.Cons >= 2 {
		c.SetCD(action.ActionSkill, 32*60*0.85)
	} else {
		c.SetCD(action.ActionSkill, 32*60)
	}
	return f, a //todo fix field cast time
}

func (c *char) barbaraHealTick(healAmt float64, hpplus float64, skillInitF int) func() {
	return func() {
		//make sure it's not overwritten
		if c.skillInitF != skillInitF {
			return
		}
		//do nothing if buff expired
		if c.Core.Status.Duration("barbskill") == 0 {
			return
		}
		// c.Core.Log.NewEvent("barbara heal ticking", core.LogCharacterEvent, c.Index)
		c.Core.Health.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.ActiveChar,
			Message: "Melody Loop (Tick)",
			Src:     healAmt,
			Bonus:   hpplus,
		})

		// tick per 5 seconds
		c.AddTask(c.barbaraHealTick(healAmt, hpplus, skillInitF), "barbara-heal-tick", 5*60)
	}
}

func (c *char) barbaraWet(ai combat.AttackInfo, skillInitF int) func() {
	return func() {
		//make sure it's not overwritten
		if c.skillInitF != skillInitF {
			return
		}
		//do nothing if buff expired
		if c.Core.Status.Duration("barbskill") == 0 {
			return
		}
		c.Core.Log.NewEvent("barbara wet ticking", glog.LogCharacterEvent, c.Index)

		c.Core.Combat.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), -1, 5)

		// tick per 5 seconds
		c.AddTask(c.barbaraWet(ai, skillInitF), "barbara-wet", 5*60)
	}
}
