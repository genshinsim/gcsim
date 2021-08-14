package eula

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("eula", NewChar)
}

type char struct {
	*character.Tmpl
	grimheartReset  int
	burstCounter    int
	burstCounterICD int
}

func NewChar(s core.Sim, log *zap.SugaredLogger, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
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

	c.Sim.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if c.Sim.Status("eulaq") == 0 {
			return
		}
		if ds.ActorIndex != c.Index {
			return
		}
		if c.burstCounterICD > c.Sim.Frame() {
			return
		}
		switch ds.AttackTag {
		case core.AttackTagElementalArt:
		case core.AttackTagElementalBurst:
		case core.AttackTagNormal:
		default:
			return
		}

		//add to counter
		c.burstCounter++
		c.Log.Debugw("eula burst add stack", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "stack count", c.burstCounter)
		//check for c6
		if c.Base.Cons == 6 && c.Sim.Rand().Float64() < 0.5 {
			c.burstCounter++
			c.Log.Debugw("eula c6 add additional stack", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "stack count", c.burstCounter)
		}
		c.burstCounterICD = c.Sim.Frame() + 6
	}, "eula-burst-counter")

	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
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
		return f
	case core.ActionCharge:
		return 35 //TODO: no idea
	case core.ActionSkill:
		if p["hold"] == 0 {
			return 34
		}
		if c.Base.Cons >= 2 {
			return 34 //press and hold have same cd
		}
		return 80
	case core.ActionBurst:
		return 116 //ok
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) a4() {
	c.Sim.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() != c.Index {
			return false
		}
		//reset CD, add 1 stack
		v := c.Tags["grimheart"]
		if v < 2 {
			v++
		}
		c.Tags["grimheart"] = v

		c.Log.Debugw("eula a4 reset skill cd", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent)
		c.ResetActionCooldown(core.ActionSkill)

		return false

	}, "eula-a4", core.PostBurstHook)
}

var delay = [][]int{{11}, {25}, {36, 49}, {33}, {45, 63}}

func (c *char) Attack(p map[string]int) int {
	//register action depending on number in chain
	//3 and 4 need to be registered as multi action

	f := c.ActionFrames(core.ActionAttack, p)

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
	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)
	if p["hold"] == 0 {
		c.pressE()
		return f
	}
	c.holdE()
	return f
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
	if c.Sim.Rand().Float64() < .5 {
		n = 2
	}
	c.QueueParticle("eula", n, core.Cryo, 100)

	//add 1 stack to Grimheart
	v := c.Tags["grimheart"]
	if v < 2 {
		v++
	}
	c.Tags["grimheart"] = v
	c.Log.Debugw("eula: grimheart stack", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "current count", v)
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
	if c.Sim.Rand().Float64() < .5 {
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
			Expiry: c.Sim.Frame() + (6*v+6)*60, //TODO: check if this is right
		})
	}

	c.Tags["grimheart"] = 0
	c.SetCD(core.ActionSkill, 10*60+62)
}

//ult 365 to 415, 60fps = 120
//looks like ult charges for 8 seconds
func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)
	c.Sim.AddStatus("eulaq", 7*60+f+1)

	c.burstCounter = 0
	if c.Base.Cons == 6 {
		c.burstCounter = 5
	}

	c.Log.Debugw("eula burst started", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "stacks", c.burstCounter, "expiry", c.Sim.Status("eulaq"))

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
	c.Log.Debugw("eula: grimheart stack", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "current count", v)

	c.AddTask(func() {
		//check to make sure it hasn't already exploded due to exiting field
		if c.Sim.Status("eulaq") > 0 {
			c.triggerBurst()
		}
	}, "Eula-Burst-Lightfall", 7*60+f) //after 8 seconds

	c.SetCD(core.ActionBurst, 20*60+f)
	//energy does not deplete until after animation
	c.AddTask(func() {
		c.Energy = 0
	}, "q", f)

	return f
}

func (c *char) onExitField() {
	c.Sim.AddEventHook(func(s core.Sim) bool {
		//trigger burst if burst is active
		if c.Sim.Status("eulaq") > 0 {
			c.triggerBurst()
		}
		return false
	}, "eula-exit", core.PostSwapHook)
}

func (c *char) c4() {
	c.Sim.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
		if ds.ActorIndex != c.Index {
			return
		}
		if ds.Abil != "Glacial Illumination (Lightfall)" {
			return
		}
		if !c.Sim.Flags().DamageMode {
			return
		}
		if t.HP()/t.MaxHP() < 0.5 {
			ds.Stats[core.DmgP] += 0.25
			c.Log.Debugw("eula: c4 adding dmg", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "final dmgp", ds.Stats[core.DmgP])
		}

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

	c.Log.Debugw("eula burst triggering", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "stacks", stacks, "mult", d.Mult)

	c.Sim.ApplyDamage(&d)
	c.Sim.DeleteStatus("eulaq")
	c.burstCounter = 0
}

func (e *char) Tick() {
	e.Tmpl.Tick()
	e.grimheartReset--
	if e.grimheartReset == 0 {
		e.Tags["grimheart"] = 0
	}
}
