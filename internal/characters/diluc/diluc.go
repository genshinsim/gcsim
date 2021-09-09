package diluc

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("diluc", NewChar)
}

type char struct {
	*character.Tmpl
	eStarted    bool
	eStartFrame int
	eLastUse    int
	eCounter    int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 40
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 4

	if c.Base.Cons >= 1 && c.Core.Flags.DamageMode {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}

	return &c, nil
}

func (c *char) c1() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		if ds.ActorIndex == c.Index && t.HP()/t.MaxHP() > .5 {
			ds.Stats[core.DmgP] += 0.15
			c.Core.Log.Debugw("diluc c2 adding dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "hp %", t.HP()/t.MaxHP(), "final dmg", ds.Stats[core.DmgP])
		}
		return false
	}, "diluc-c1")

}

func (c *char) c2() {
	stack := 0
	last := 0
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if last != 0 && c.Core.F-last < 90 {
			return false
		}
		//last time is more than 10 seconds ago, reset stacks back to 0
		if c.Core.F-last > 600 {
			stack = 0
		}
		stack++
		if stack > 3 {
			stack = 3
		}
		val := make([]float64, core.EndStatType)
		val[core.ATKP] = 0.1 * float64(stack)
		val[core.AtkSpd] = 0.05 * float64(stack)
		c.AddMod(core.CharStatMod{
			Key:    "diluc-c2",
			Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
			Expiry: c.Core.F + 600,
		})
		return false
	}, "diluc-c2")

}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 24 //frames from keqing lib
		case 1:
			f = 53
		case 2:
			f = 38
		case 3:
			f = 65
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 0
	case core.ActionSkill:
		switch c.eCounter {
		case 1:
			return 52
		case 2:
			return 81
		default:
			return 45
		}
	case core.ActionBurst:
		return 65
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f-1)
	c.AdvanceNormalIndex()

	return f
}

func (c *char) Skill(p map[string]int) int {

	f := c.ActionFrames(core.ActionSkill, p)

	if c.eCounter == 0 {
		c.eStarted = true
		c.eStartFrame = c.Core.F
	}
	c.eLastUse = c.Core.F

	orb := 1
	if c.Core.Rand.Float64() < 0.33 {
		orb = 2
	}
	c.QueueParticle("Diluc", orb, core.Pyro, f+60)

	//actual skill cd starts immediately on first cast
	//times out after 4 seconds of not using
	//every hit applies pyro
	//apply attack speed

	d := c.Snapshot(
		"Searing Onslaught",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Pyro,
		25,
		skill[c.eCounter][c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	//check for c4 dmg increase
	if c.Base.Cons >= 4 {
		if c.Core.Status.Duration("dilucc4") > 0 {
			d.Stats[core.DmgP] += 0.4
			c.Core.Log.Debugw("diluc c4 adding dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "final dmg", d.Stats[core.DmgP])
		}
	}

	c.QueueDmg(&d, f-5)

	//add a timer to activate c4
	if c.Base.Cons >= 4 {
		c.AddTask(func() {
			c.Core.Status.AddStatus("dilucc4", 120) //effect lasts 2 seconds
		}, "dilucc4", f+120) // 2seconds after cast
	}

	c.eCounter++
	if c.eCounter == 3 {
		//ability can go on cd now
		cd := 600 - (c.Core.F - c.eStartFrame)
		c.Core.Log.Debugw("diluc skill going on cd", "frame", c.Core.F, "event", core.LogCharacterEvent, "duration", cd)
		c.SetCD(core.ActionSkill, cd)
		c.eStarted = false
		c.eStartFrame = -1
		c.eLastUse = -1
		c.eCounter = 0
	}
	//return animation cd
	//this also depends on which hit in the chain this is
	return f
}

func (c *char) Burst(p map[string]int) int {

	dot, ok := p["dot"]
	if !ok {
		dot = 2 //number of dot hits
	}
	if dot > 7 {
		dot = 7
	}
	explode, ok := p["explode"]
	if !ok {
		explode = 0 //if explode hits
	}

	// c.S.Status["dilucq"] = c.Core.F + 12*60
	c.Core.Status.AddStatus("dilucq", 720)
	f := c.ActionFrames(core.ActionBurst, p)

	d := c.Snapshot(
		"Dawn (Strike)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDiluc,
		core.StrikeTypeBlunt,
		core.Pyro,
		50,
		burstInitial[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, 100)

	//dot does damage every .2 seconds for 7 hits? so every 12 frames
	//dot does max 7 hits + explosion, roughly every 13 frame? blows up at 210 frames
	//first tick did 50 dur as well?
	for i := 1; i <= dot; i++ {
		x := d.Clone()
		x.Abil = "Dawn (Tick)"
		x.Mult = burstDOT[c.TalentLvlBurst()]
		c.QueueDmg(&x, 100+i+12)
	}

	if explode > 0 {
		x := d.Clone()
		x.Abil = "Dawn (Explode)"
		x.Mult = burstExplode[c.TalentLvlBurst()]
		c.QueueDmg(&x, 210)
	}

	//enhance weapon for 10.2 seconds
	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "diluc-fire-weapon",
		Ele:    core.Pyro,
		Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.F + 852, //with a4
	})

	// add 20% pyro damage
	val := make([]float64, core.EndStatType)
	val[core.PyroP] = 0.2
	c.AddMod(core.CharStatMod{
		Key:    "diluc-fire-weapon",
		Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Core.F + 852,
	})

	//c.Energy = 0  forcing every character to comsume energy after burts in the energy.go to make my life easier
	c.ConsumeEnergy(0, 0) //at 0,0 value acts the same as c.Energy = 0
	c.SetCD(core.ActionBurst, 900)
	return f
}

func (c *char) Tick() {
	c.Tmpl.Tick()

	if c.eStarted {
		//check if 4 second has passed since last use
		if c.Core.F-c.eLastUse >= 240 {
			//if so, set ability to be on cd equal to 10s less started
			cd := 600 - (c.Core.F - c.eStartFrame)
			c.Core.Log.Debugw("diluc skill going on cd", "frame", c.Core.F, "event", core.LogCharacterEvent, "duration", cd, "last", c.eLastUse)
			c.SetCD(core.ActionSkill, cd)
			//reset
			c.eStarted = false
			c.eStartFrame = -1
			c.eLastUse = -1
			c.eCounter = 0
		}
	}
}

// func (c *char) Snapshot(name string, a def.AttackTag, icd def.ICDTag, g def.ICDGroup, st def.StrikeType, e def.EleType, d float64, mult float64) def.Snapshot {
// 	ds := c.CharacterTemplate.Snapshot(name, a, icd, g, st, e, d, mult)
// 	if c.S.StatusActive("dilucq") {
// 		if ds.AttackTag == def.AttackTagNormal || ds.AttackTag == def.AttackTagExtra {
// 			ds.Element = def.Pyro
// 			ds.Stats[def.PyroP] += .2
// 		}
// 	}
// 	return ds
// }

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 50
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}
