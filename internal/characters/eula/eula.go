package eula

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("eula", NewChar)
}

type char struct {
	*character.Tmpl
	grimheartReset int
	burstCounter   int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.MaxEnergy = 80
	c.Weapon.Class = def.WeaponClassClaymore
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.a4()
	c.onExitField()

	c.S.AddSnapshotHook(func(ds *def.Snapshot) bool {
		if c.Sim.Status("eulaq") == 0 {
			return false
		}
		if ds.Actor != c.Base.Name {
			return false
		}
		if c.Sim.Status("Eula Burst ICD") > 0 {
			return false
		}
		switch ds.AttackTag {
		case def.AttackTagElementalArt:
		case def.AttackTagElementalBurst:
		case def.AttackTagNormal:
		default:
			return false
		}

		//add to counter
		c.burstCounter++
		c.Log.Debugw("eula burst add stack", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "stack count", c.burstCounter)
		//check for c6
		if c.Base.Cons == 6 && c.S.Rand.Float64() < 0.5 {
			c.burstCounter++
			c.Log.Debugw("eula c6 add additional stack", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "stack count", c.burstCounter)
		}
		c.S.Status["Eula Burst ICD"] = c.Sim.Frame() + 6
		return false
	}, "Eula Burst", def.PostDamageHook)

	return &c, nil
}

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionCharge:
		return 35 //TODO: no idea
	case def.ActionSkill:
		if p["hold"] == 0 {
			return 34
		}
		if c.Base.Cons >= 2 {
			return 34 //press and hold have same cd
		}
		return 80
	case def.ActionBurst:
		return 116 //ok
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) a4() {
	c.S.AddEventHook(func(s *def.Sim) bool {
		if s.ActiveChar != c.Base.Name {
			return false
		}
		//reset CD, add 1 stack
		v := c.Tags["grimheart"]
		if v < 2 {
			v++
		}
		c.Tags["grimheart"] = v

		c.Log.Debugw("eula a4 reset skill cd", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent)
		c.ResetActionCooldown(def.ActionSkill)

		return false
	}, "eula-a4", def.PostBurstHook)
}

var delay = [][]int{{11}, {25}, {36, 49}, {33}, {45, 63}}

func (c *char) Attack(p map[string]int) int {
	//register action depending on number in chain
	//3 and 4 need to be registered as multi action

	f := c.ActionFrames(def.ActionAttack, p)

	//apply attack speed
	d := c.Snapshot(
		"Normal",
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Physical,
		25,
		0,
	)

	for i, mult := range auto[c.NormalCounter] {
		x := d.Clone()
		x.Mult = mult[c.TalentLvlAttack()]
		x.Targets = def.TargetAll
		c.QueueDmg(&x, delay[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	//return animation cd
	//this also depends on which hit in the chain this is
	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)
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
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Cryo,
		25,
		skillPress[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, 35)

	n := 1
	if c.Sim.Rand().Float64() < .5 {
		n = 2
	}
	c.QueueParticle("eula", n, def.Cryo, 100)

	//add 1 stack to Grimheart
	v := c.Tags["grimheart"]
	if v < 2 {
		v++
	}
	c.Tags["grimheart"] = v
	c.Log.Debugw("eula: grimheart stack", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "current count", v)
	c.grimheartReset = 18 * 60

	c.SetCD(def.ActionSkill, 240)
}

func (c *char) holdE() {
	//hold e
	//296 to 341, but cd starts at 322
	//60 fps = 108 frames cast, cd starts 62 frames in so need to + 62 frames to cd
	lvl := c.TalentLvlSkill()
	d := c.Snapshot(
		"Icetide Vortex (Hold)",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Cryo,
		25,
		skillHold[lvl],
	)
	d.Targets = def.TargetAll
	c.QueueDmg(&d, 80)

	//multiple brand hits
	v := c.Tags["grimheart"]

	d = c.Snapshot(
		"Icetide Vortex (Icewhirl)",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Cryo,
		25,
		icewhirl[lvl],
	)
	d.Targets = def.TargetAll

	for i := 0; i < v; i++ {
		x := d.Clone()
		c.QueueDmg(&x, 90+i*7) //we're basically forcing it so we get 3 stacks
	}

	//A2
	if v == 2 {
		d := c.Snapshot(
			"Icetide (Lightfall)",
			def.AttackTagElementalBurst,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeBlunt,
			def.Physical,
			25,
			burstExplodeBase[c.TalentLvlBurst()]*0.5,
		)
		d.Targets = def.TargetAll
		c.QueueDmg(&d, 108) //we're basically forcing it so we get 3 stacks
	}

	n := 2
	if c.Sim.Rand().Float64() < .5 {
		n = 3
	}
	c.QueueParticle("eula", n, def.Cryo, 100)

	//add debuff per hit
	c.S.Target.AddResMod("Icewhirl Cryo", def.ResistMod{
		Ele:      def.Cryo,
		Value:    -resRed[lvl],
		Duration: 7 * v * 60,
	})
	c.S.Target.AddResMod("Icewhirl Physical", def.ResistMod{
		Ele:      def.Physical,
		Value:    -resRed[lvl],
		Duration: 7 * v * 60,
	})

	//c1 add debuff
	if c.Base.Cons >= 1 && v > 0 {
		val := make([]float64, def.EndStatType)
		val[def.PhyP] = 0.3
		c.AddMod(def.CharStatMod{
			Key:       "eula-c1",
			Mod:       val,
			Condition: func(ds def.Snapshot) bool { return true },
			Expiry:    c.Sim.Frame() + (6*v+6)*60, //TODO: check if this is right
		})
	}

	c.Tags["grimheart"] = 0
	c.SetCD(def.ActionSkill, 10*60+62)
}

func (c *char) skillResReduction() {
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {

	}, "eula-grimheart-debuff")
}

//ult 365 to 415, 60fps = 120
//looks like ult charges for 8 seconds
func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(def.ActionBurst, p)
	c.Sim.AddStatus("eulaq", 7*60+f+1)

	c.burstCounter = 0
	if c.Base.Cons == 6 {
		c.burstCounter = 5
	}

	c.Log.Debugw("eula burst started", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "stacks", c.burstCounter, "expiry", c.S.Status["eulaq"])

	lvl := c.TalentLvlBurst()
	//add initial damage

	d := c.Snapshot(
		"Glacial Illumination",
		def.AttackTagElementalBurst,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Cryo,
		50,
		burstInitial[lvl],
	)

	c.QueueDmg(&d, f-1)

	//50% chance of 2 to 3 balls
	count := 2
	if c.S.Rand.Float64() < .5 {
		count = 3
	}
	c.S.AddEnergyParticles("eula", count, def.Cryo, f+60) //guess

	//add 1 stack to Grimheart
	v := c.Tags["grimheart"]
	if v < 2 {
		v++
	}
	c.Tags["grimheart"] = v
	c.Log.Debugw("eula: grimheart stack", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "current count", v)

	c.S.AddTask(func(s *def.Sim) {
		//check to make sure it hasn't already exploded due to exiting field
		if c.Sim.Status("eulaq") > 0 {
			c.triggerBurst()
		}
	}, "Eula-Burst-Lightfall", 7*60+f) //after 8 seconds

	c.SetCD(def.ActionBurst, 20*60+f)
	//energy does not deplete until after animation
	c.AddTask(func() {
		c.Energy = 0
	}, "q", f)

	return f
}

func (c *char) onExitField() {
	c.S.AddEventHook(func(s *def.Sim) bool {
		//trigger burst if burst is active
		if c.Sim.Status("eulaq") > 0 {
			c.triggerBurst()
		}
		return false
	}, "eula-exit", def.PostSwapHook)
}

func (c *char) triggerBurst() {

	stacks := c.burstCounter
	if stacks > 30 {
		stacks = 30
	}

	d := c.Snapshot(
		"Glacial Illumination (Lightfall)",
		def.AttackTagElementalBurst,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Physical,
		50,
		burstExplodeBase[c.TalentLvlBurst()]+burstExplodeStack[c.TalentLvlBurst()]*float64(stacks),
	)

	hp := c.S.Target.HP / c.S.Target.MaxHP

	c.Log.Debugw("eula burst triggering", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "stacks", stacks, "mult", d.Mult, "target hp%", hp)

	//check target hp if < 25% and c4
	if c.Base.Cons >= 4 && c.S.Target.HPMode && hp < 0.5 {
		d.Stats[def.DmgP] += 0.25
		c.Log.Debugw("eula: c4 adding dmg", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "final dmgp", d.Stats[def.DmgP])
	}
	c.S.ApplyDamage(d)
	c.S.DeleteStatus("eulaq")
	c.burstCounter = 0
}

func (e *char) Tick() {
	e.Tmpl.Tick()
	e.grimheartReset--
	if e.grimheartReset == 0 {
		e.Tags["grimheart"] = 0
	}
}
