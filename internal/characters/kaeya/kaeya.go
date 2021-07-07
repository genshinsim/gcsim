package kaeya

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("kaeya", NewChar)
}

type char struct {
	*character.Tmpl
	c4icd     int
	a4count   int
	icicleICD []int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.MaxEnergy = 60
	c.Weapon.Class = def.WeaponClassSword
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

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionCharge:
		return 87
	case def.ActionSkill:
		return 58 //could be 52 if going into Q
	case def.ActionBurst:
		return 78
	default:
		c.Log.Warnf("%v: unknown action, frames invalid", a)
		return 0
	}
}

func (c *char) ActionStam(a def.ActionType, p map[string]int) float64 {
	switch a {
	case def.ActionDash:
		return 18
	case def.ActionCharge:
		return 25
	default:
		c.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}
}

func (c *char) a4() {
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		if t.AuraType() != def.Frozen {
			return
		}
		if c.a4count == 2 {
			return
		}
		c.a4count++
		c.QueueParticle("kaeya", 1, def.Cryo, 100)
		c.Log.Debugw("kaeya a4 proc", "event", def.LogEnergyEvent, "char", c.Index, "frame", c.Sim.Frame(), "final cr", ds.Stats[def.CR])
	}, "kaeya-a4")

}

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(def.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
		25,
		auto[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)

	c.AdvanceNormalIndex()

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {

	f := c.ActionFrames(def.ActionCharge, p)

	d := c.Snapshot(
		"Charge 1",
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
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
	f := c.ActionFrames(def.ActionSkill, p)
	d := c.Snapshot(
		"Frostgnaw",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Cryo,
		50,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll
	c.a4count = 0

	//2 or 3 1:1 ratio
	count := 2
	if c.Sim.Rand().Float64() < 0.67 {
		count = 3
	}
	c.QueueParticle("kaeya", count, def.Cryo, f+100)

	//add a2
	heal := .15 * (d.BaseAtk*(1+d.Stats[def.ATKP]) + d.Stats[def.ATK])
	c.AddTask(func() {
		c.Sim.HealActive(heal)
		//apply damage
		c.Sim.ApplyDamage(&d)
	}, "Kaeya-Skill", 28) //TODO: assumed same as when cd starts

	c.SetCD(def.ActionSkill, 360+28) //+28 since cd starts 28 frames in
	return f
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(def.ActionBurst, p)
	d := c.Snapshot(
		"Glacial Waltz",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Cryo,
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
			x.Targets = def.TargetAll
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
	c.SetCD(def.ActionBurst, 900)
	return f
}

func (c *char) burstICD() {
	c.Sim.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.Abil != "Glacial Waltz" {
			return
		}
		//check icd
		if c.icicleICD[ds.ExtraIndex] > c.Sim.Frame() {
			ds.Cancelled = true
			return
		}
		c.icicleICD[ds.ExtraIndex] = c.Sim.Frame() + 30

	}, "kaeya-burst-icd")
}
