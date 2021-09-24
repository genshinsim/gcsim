package kaeya

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("kaeya", NewChar)
}

type char struct {
	*character.Tmpl
	c4icd     int
	a4count   int
	icicleICD []int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5

	c.icicleICD = make([]int, 4)
	c.a4()
	// c.burstICD()

	if c.Base.Cons > 0 {
		c.c1()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 14 //frames from keqing lib
		case 1:
			f = 41 - 14
		case 2:
			f = 72 - 41
		case 3:
			f = 128 - 72
		case 4:
			f = 176 - 128
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 87
	case core.ActionSkill:
		return 58 //could be 52 if going into Q
	case core.ActionBurst:
		return 78
	default:
		c.Core.Log.Warnf("%v: unknown action, frames invalid", a)
		return 0
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}
}

func (c *char) a4() {
	c.Core.Events.Subscribe(core.OnReactionOccured, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		if ds.ActorIndex != c.Index {
			return false
		}
		if ds.AttackTag != core.AttackTagElementalArt {
			return false
		}
		if t.AuraType() != core.Frozen {
			return false
		}
		if c.a4count == 2 {
			return false
		}
		c.a4count++
		c.QueueParticle("kaeya", 1, core.Cryo, 100)
		c.Core.Log.Debugw("kaeya a4 proc", "event", core.LogEnergyEvent, "char", c.Index, "frame", c.Core.F, "final cr", ds.Stats[core.CR])
		return false
	}, "kaeya-a4")

}

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		auto[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	c.AdvanceNormalIndex()

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {

	f := c.ActionFrames(core.ActionCharge, p)

	d := c.Snapshot(
		"Charge 1",
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		charge[0][c.TalentLvlAttack()],
	)
	d2 := d.Clone()
	d2.Abil = "Charge 2"
	d2.Mult = charge[1][c.TalentLvlAttack()]

	c.QueueDmg(&d, f-15) //TODO: damage frame
	c.QueueDmg(&d2, f-5) //TODO: damage frame

	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Frostgnaw",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		50,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll
	c.a4count = 0

	//2 or 3 1:1 ratio
	count := 2
	if c.Core.Rand.Float64() < 0.67 {
		count = 3
	}
	c.QueueParticle("kaeya", count, core.Cryo, f+100)

	//add a2
	heal := .15 * (d.BaseAtk*(1+d.Stats[core.ATKP]) + d.Stats[core.ATK])
	c.AddTask(func() {
		c.Core.Health.HealActive(heal)
		//apply damage
		c.Core.Combat.ApplyDamage(&d)
	}, "Kaeya-Skill", 28) //TODO: assumed same as when cd starts

	c.SetCD(core.ActionSkill, 360+28) //+28 since cd starts 28 frames in
	return f
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)
	d := c.Snapshot(
		"Glacial Waltz",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Cryo,
		25,
		burst[c.TalentLvlBurst()],
	)

	//duration starts counting 49 frames in per kqm lib
	//hits around 13 times

	//each icicle takes 120frames to complete a rotation and has a internal cooldown of 0.5
	count := 3
	if c.Base.Cons == 6 {
		count++
	}
	offset := 120 / count

	for i := 0; i < count; i++ {

		//each icicle will start at i * offset (i.e. 0, 40, 80 OR 0, 30, 60, 90)
		//assume each icicle will last for 8 seconds
		//assume damage dealt every 120 (since only hitting at the front)
		//on icicle collision, it'll trigger an aoe dmg with radius 2
		//in effect, every target gets hit every time icicles rotate around
		for j := f + offset*i; j < f+480; j += 120 {
			x := d.Clone()
			x.Targets = core.TargetAll
			x.ExtraIndex = i
			c.QueueDmg(&x, j)
		}

	}

	if c.Base.Cons == 6 {
		c.Energy = 15
	} else {
		c.Energy = 0
	}

	// c.CD[def.BurstCD] = c.Sim.F + 900
	c.SetCD(core.ActionBurst, 900)
	return f
}

func (c *char) burstICD() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		if ds.Abil != "Glacial Waltz" {
			return false
		}
		//check icd
		if c.icicleICD[ds.ExtraIndex] > c.Core.F {
			ds.Cancelled = true
			return false
		}
		c.icicleICD[ds.ExtraIndex] = c.Core.F + 30
		return false
	}, "kaeya-burst-icd")
}
