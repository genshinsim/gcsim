package eula

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("eula", NewChar)
}

type char struct {
	*character.Tmpl
	grimheartReset  int
	burstCounter    int
	burstCounterICD int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassClaymore
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.a4()
	c.onExitField()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if c.Core.Status.Duration("eulaq") == 0 {
			return false
		}
		if ds.ActorIndex != c.Index {
			return false
		}
		if c.burstCounterICD > c.Core.F {
			return false
		}
		switch ds.AttackTag {
		case core.AttackTagElementalArt:
		case core.AttackTagElementalBurst:
		case core.AttackTagNormal:
		default:
			return false
		}

		//add to counter
		c.burstCounter++
		c.Core.Log.Debugw("eula burst add stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "stack count", c.burstCounter)
		//check for c6
		if c.Base.Cons == 6 && c.Core.Rand.Float64() < 0.5 {
			c.burstCounter++
			c.Core.Log.Debugw("eula c6 add additional stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "stack count", c.burstCounter)
		}
		c.burstCounterICD = c.Core.F + 6
		return false
	}, "eula-burst-counter")
	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 29
		case 1:
			f = 25
		case 2:
			f = 65
		case 3:
			f = 33
		case 4:
			f = 88
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f, f
	case core.ActionCharge:
		return 35, 35 //TODO: no idea
	case core.ActionSkill:
		if p["hold"] == 0 {
			return 34, 34
		}
		if c.Base.Cons >= 2 {
			return 34, 34 //press and hold have same cd
		}
		return 80, 80
	case core.ActionBurst:
		return 116, 116 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0, 0
	}
}

func (c *char) a4() {
	c.Core.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index {
			return false
		}
		//reset CD, add 1 stack
		v := c.Tags["grimheart"]
		if v < 2 {
			v++
		}
		c.Tags["grimheart"] = v

		c.Core.Log.Debugw("eula a4 reset skill cd", "frame", c.Core.F, "event", core.LogCharacterEvent)
		c.ResetActionCooldown(core.ActionSkill)

		return false
	}, "eula-a4")
}

var delay = [][]int{{11}, {25}, {36, 49}, {33}, {45, 63}}

func (c *char) Attack(p map[string]int) (int, int) {
	//register action depending on number in chain
	//3 and 4 need to be registered as multi action

	f, a := c.ActionFrames(core.ActionAttack, p)

	//apply attack speed
	d := c.Snapshot(
		"Normal",
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Physical,
		25,
		0,
	)
	d.Targets = core.TargetAll

	for i, mult := range auto[c.NormalCounter] {
		x := d.Clone()
		x.Mult = mult[c.TalentLvlAttack()]
		x.Targets = core.TargetAll
		c.QueueDmg(&x, delay[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	//return animation cd
	//this also depends on which hit in the chain this is
	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	if p["hold"] == 0 {
		c.pressE()
		return f, a
	}
	c.holdE()
	return f, a
}

func (c *char) pressE() {
	//press e (60fps vid)
	//starts 343 cancel 378
	d := c.Snapshot(
		"Icetide Vortex",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Cryo,
		25,
		skillPress[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, 35)

	n := 1
	if c.Core.Rand.Float64() < .5 {
		n = 2
	}
	c.QueueParticle("eula", n, core.Cryo, 100)

	//add 1 stack to Grimheart
	v := c.Tags["grimheart"]
	if v < 2 {
		v++
	}
	c.Tags["grimheart"] = v
	c.Core.Log.Debugw("eula: grimheart stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "current count", v)
	c.grimheartReset = 18 * 60

	c.SetCD(core.ActionSkill, 240)
}

func (c *char) holdE() {
	//hold e
	//296 to 341, but cd starts at 322
	//60 fps = 108 frames cast, cd starts 62 frames in so need to + 62 frames to cd
	lvl := c.TalentLvlSkill()
	d := c.Snapshot(
		"Icetide Vortex (Hold)",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Cryo,
		25,
		skillHold[lvl],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, 80)

	//multiple brand hits
	v := c.Tags["grimheart"]

	d = c.Snapshot(
		"Icetide Vortex (Icewhirl)",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Cryo,
		25,
		icewhirl[lvl],
	)
	d.Targets = core.TargetAll

	for i := 0; i < v; i++ {
		x := d.Clone()
		c.QueueDmg(&x, 90+i*7) //we're basically forcing it so we get 3 stacks
	}

	//shred
	if v > 0 {
		d.OnHitCallback = func(t core.Target) {
			t.AddResMod("Icewhirl Cryo", core.ResistMod{
				Ele:      core.Cryo,
				Value:    -resRed[lvl],
				Duration: 7 * v * 60,
			})
			t.AddResMod("Icewhirl Physical", core.ResistMod{
				Ele:      core.Physical,
				Value:    -resRed[lvl],
				Duration: 7 * v * 60,
			})

		}
	}

	//A2
	if v == 2 {
		d := c.Snapshot(
			"Icetide (Lightfall)",
			core.AttackTagElementalBurst,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeBlunt,
			core.Physical,
			25,
			burstExplodeBase[c.TalentLvlBurst()]*0.5,
		)
		d.Targets = core.TargetAll
		c.QueueDmg(&d, 108) //we're basically forcing it so we get 3 stacks
	}

	n := 2
	if c.Core.Rand.Float64() < .5 {
		n = 3
	}
	c.QueueParticle("eula", n, core.Cryo, 100)

	//c1 add debuff
	if c.Base.Cons >= 1 && v > 0 {
		val := make([]float64, core.EndStatType)
		val[core.PhyP] = 0.3
		c.AddMod(core.CharStatMod{
			Key: "eula-c1",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return val, true
			},
			Expiry: c.Core.F + (6*v+6)*60, //TODO: check if this is right
		})
	}

	c.Tags["grimheart"] = 0
	c.SetCD(core.ActionSkill, 10*60+62)
}

//ult 365 to 415, 60fps = 120
//looks like ult charges for 8 seconds
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	c.Core.Status.AddStatus("eulaq", 7*60+f+1)

	c.burstCounter = 0
	if c.Base.Cons == 6 {
		c.burstCounter = 5
	}

	c.Core.Log.Debugw("eula burst started", "frame", c.Core.F, "event", core.LogCharacterEvent, "stacks", c.burstCounter, "expiry", c.Core.Status.Duration("eulaq"))

	lvl := c.TalentLvlBurst()
	//add initial damage

	d := c.Snapshot(
		"Glacial Illumination",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Cryo,
		50,
		burstInitial[lvl],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f-1)

	//add 1 stack to Grimheart
	v := c.Tags["grimheart"]
	if v < 2 {
		v++
	}
	c.Tags["grimheart"] = v
	c.Core.Log.Debugw("eula: grimheart stack", "frame", c.Core.F, "event", core.LogCharacterEvent, "current count", v)

	c.AddTask(func() {
		//check to make sure it hasn't already exploded due to exiting field
		if c.Core.Status.Duration("eulaq") > 0 {
			c.triggerBurst()
		}
	}, "Eula-Burst-Lightfall", 7*60+f) //after 8 seconds

	c.SetCD(core.ActionBurst, 20*60+f)
	//energy does not deplete until after animation
	c.AddTask(func() {
		c.Energy = 0
	}, "q", f)

	return f, a
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		if c.Core.Status.Duration("eulaq") > 0 {
			c.triggerBurst()
		}
		return false
	}, "eula-exit")
}

func (c *char) c4() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		if ds.Abil != "Glacial Illumination (Lightfall)" {
			return false
		}
		if !c.Core.Flags.DamageMode {
			return false
		}
		if t.HP()/t.MaxHP() < 0.5 {
			ds.Stats[core.DmgP] += 0.25
			c.Core.Log.Debugw("eula: c4 adding dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "final dmgp", ds.Stats[core.DmgP])
		}
		return false
	}, "eula-c4")
}

func (c *char) triggerBurst() {

	stacks := c.burstCounter
	if stacks > 30 {
		stacks = 30
	}

	d := c.Snapshot(
		"Glacial Illumination (Lightfall)",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Physical,
		50,
		burstExplodeBase[c.TalentLvlBurst()]+burstExplodeStack[c.TalentLvlBurst()]*float64(stacks),
	)
	d.Targets = core.TargetAll

	c.Core.Log.Debugw("eula burst triggering", "frame", c.Core.F, "event", core.LogCharacterEvent, "stacks", stacks, "mult", d.Mult)

	c.Core.Combat.ApplyDamage(&d)
	c.Core.Status.DeleteStatus("eulaq")
	c.burstCounter = 0
}

func (e *char) Tick() {
	e.Tmpl.Tick()
	e.grimheartReset--
	if e.grimheartReset == 0 {
		e.Tags["grimheart"] = 0
	}
}
