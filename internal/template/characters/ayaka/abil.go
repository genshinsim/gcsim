package ayaka

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

var hitmarks = [][]int{{8}, {10}, {16}, {8, 15, 22}, {27}}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	c1cbArgs := make([]core.AttackCBFunc, 0, 1)
	if c.Base.Cons >= 1 {
		c1cbArgs = append(c1cbArgs, c.c1cb)
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), hitmarks[c.NormalCounter][i], hitmarks[c.NormalCounter][i], c1cbArgs...)
	}

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		Abil:       "Charge",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       ca[c.TalentLvlAttack()],
	}

	cbArgs := make([]core.AttackCBFunc, 0, 1)
	if c.Base.Cons >= 1 {
		cbArgs = append(cbArgs, c.c1cb)
	}
	if c.Base.Cons >= 6 {
		cbArgs = append(cbArgs, c.c6cb)
	}

	for i := 0; i < 3; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f-2+i, f-2+i, cbArgs...)
	}

	return f, a
}

func (c *char) Dash(p map[string]int) (int, int) {
	f, ok := p["f"]
	a := f + 10
	if !ok {
		f, a = c.ActionFrames(core.ActionDash, p)
	}

	//no dmg attack at end of dash
	ai := core.AttackInfo{
		Abil:       "Dash",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNone,
		ICDTag:     core.ICDTagDash,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
	}

	//restore on hit, once per attack
	once := false
	cb := func(a core.AttackCB) {
		if once {
			return
		}

		c.Core.RestoreStam(10)
		val := make([]float64, core.EndStatType)
		val[core.CryoP] = 0.18
		//a1 increase normal + ca dmg by 30% for 6s
		c.AddMod(core.CharStatMod{
			Key:    "ayaka-a4",
			Expiry: c.Core.F + 600,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})
		once = true
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f, cb)
	//add cryo infuse
	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "ayaka-dash",
		Ele:    core.Cryo,
		Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.F + 300,
	})
	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		Abil:       "Hyouka",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}

	//2 or 3 1:1 ratio
	count := 4
	if c.Core.Rand.Float64() < 0.5 {
		count = 5
	}
	c.QueueParticle("ayaka", count, core.Cryo, f+100)

	//a1 increase normal + ca dmg by 30% for 6s
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.3
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "ayaka-a1",
		Expiry: c.Core.F + 360,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			return val, atk.Info.AttackTag == core.AttackTagNormal || atk.Info.AttackTag == core.AttackTagExtra
		},
	})

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(4, false, core.TargettableEnemy), 0, 33)

	c.SetCD(core.ActionSkill, 600)
	return f, a

}

func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	ai := core.AttackInfo{
		Abil:       "Soumetsu",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
	}

	c4cbArgs := make([]core.AttackCBFunc, 0, 1)
	if c.Base.Cons >= 4 {
		c4cbArgs = append(c4cbArgs, c.c4cb)
	}

	//5 second, 20 ticks, so once every 15 frames, bloom after 5 seconds
	ai.Mult = burstBloom[c.TalentLvlBurst()]
	ai.Abil = "Soumetsu (Bloom)"
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f, f+300, c4cbArgs...)

	// C2 mini-frostflake bloom
	var aiC2 core.AttackInfo
	if c.Base.Cons >= 2 {
		aiC2 = ai
		aiC2.Mult = burstBloom[c.TalentLvlBurst()] * .2
		aiC2.Abil = "C2 Mini-Frostflake Seki no To (Bloom)"
		// TODO: Not sure about the positioning/size...
		c.Core.Combat.QueueAttack(aiC2, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f+300, c4cbArgs...)
		c.Core.Combat.QueueAttack(aiC2, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f+300, c4cbArgs...)
	}

	for i := 0; i < 19; i++ {
		ai.Mult = burstCut[c.TalentLvlBurst()]
		ai.Abil = "Soumetsu (Cutting)"
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), f, f+i*15, c4cbArgs...)

		// C2 mini-frostflake cutting
		if c.Base.Cons >= 2 {
			aiC2.Mult = burstCut[c.TalentLvlBurst()] * .2
			aiC2.Abil = "C2 Mini-Frostflake Seki no To (Cutting)"
			// TODO: Not sure about the positioning/size...
			c.Core.Combat.QueueAttack(aiC2, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f+i*15, c4cbArgs...)
			c.Core.Combat.QueueAttack(aiC2, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f+i*15, c4cbArgs...)
		}
	}

	c.SetCD(core.ActionBurst, 20*60)
	c.ConsumeEnergy(8)

	return f, a
}

// Callback for Ayaka C1 that is attached to NA/CA hits
func (c *char) c1cb(a core.AttackCB) {
	// When Kamisato Ayaka's Normal or Charged Attacks deal Cryo DMG to opponents, it has a 50% chance of decreasing the CD of Kamisato Art: Hyouka by 0.3s. This effect can occur once every 0.1s.
	if a.AttackEvent.Info.Element != core.Cryo {
		return
	}
	if c.icdC1 > c.Core.F {
		return
	}
	if c.Core.Rand.Float64() < .5 {
		return
	}
	c.ReduceActionCooldown(core.ActionSkill, 18)
	c.icdC1 = c.Core.F + 6
}

// Callback for Ayaka C4 that is attached to Burst hits
func (c *char) c4cb(a core.AttackCB) {
	// Opponents damaged by Kamisato Art: Soumetsu's Frostflake Seki no To will have their DEF decreased by 30% for 6s.
	a.Target.AddDefMod("ayaka-c4", -0.3, 60*6)
}

// Callback for Ayaka C6 that is attached to CA hits
func (c *char) c6cb(a core.AttackCB) {
	if !c.c6CDTimerAvail {
		return
	}

	c.c6CDTimerAvail = false

	c.AddTask(func() {
		c.DeletePreDamageMod("ayaka-c6")

		c.AddTask(func() {
			c.c6CDTimerAvail = true
			c.c6AddBuff()
		}, "ayaka-c6-reset", 600)
	}, "ayaka-c6-clear", 30)
}

func (c *char) c6AddBuff() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 2.98
	c.AddPreDamageMod(core.PreDamageMod{
		Key: "ayaka-c6",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag != core.AttackTagExtra {
				return nil, false
			}
			return val, true
		},
		Expiry: -1,
	})
}
