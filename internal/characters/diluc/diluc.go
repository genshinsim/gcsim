package diluc

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("diluc", NewChar)
}

type char struct {
	*character.Tmpl
	eStarted    bool
	eStartFrame int
	eLastUse    int
	eCounter    int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 40
	c.EnergyMax = 40
	c.Weapon.Class = def.WeaponClassClaymore
	c.NormalHitNum = 4

	if c.Base.Cons >= 1 && s.Flags().HPMode {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}

	return &c, nil
}

func (c *char) c1() {
	c.Sim.AddOnAttackWillLand(func(t def.Target, ds *def.Snapshot) {
		if ds.ActorIndex != c.Index {
			return
		}
		if t.HP()/t.MaxHP() > .5 {
			ds.Stats[def.DmgP] += 0.15
			c.Log.Debugw("diluc c2 adding dmg", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "hp %", t.HP()/t.MaxHP(), "final dmg", ds.Stats[def.DmgP])
		}
	}, "diluc-c1")
}

func (c *char) c2() {
	stack := 0
	last := 0
	c.Sim.AddOnHurt(func(s def.Sim) {
		if last != 0 && c.Sim.Frame()-last < 90 {
			return
		}
		//last time is more than 10 seconds ago, reset stacks back to 0
		if c.Sim.Frame()-last > 600 {
			stack = 0
		}
		stack++
		if stack > 3 {
			stack = 3
		}
		val := make([]float64, def.EndStatType)
		val[def.ATKP] = 0.1 * float64(stack)
		val[def.AtkSpd] = 0.05 * float64(stack)
		c.AddMod(def.CharStatMod{
			Key:    "diluc-c2",
			Amount: func(a def.AttackTag) ([]float64, bool) { return val, true },
			Expiry: c.Sim.Frame() + 600,
		})
	})

}

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionCharge:
		return 0
	case def.ActionSkill:
		switch c.eCounter {
		case 1:
			return 52
		case 2:
			return 81
		default:
			return 45
		}
	case def.ActionBurst:
		return 65
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(def.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, f-1)
	c.AdvanceNormalIndex()

	return f
}

func (c *char) Skill(p map[string]int) int {

	f := c.ActionFrames(def.ActionSkill, p)

	if c.eCounter == 0 {
		c.eStarted = true
		c.eStartFrame = c.Sim.Frame()
	}
	c.eLastUse = c.Sim.Frame()

	orb := 1
	if c.Sim.Rand().Float64() < 0.33 {
		orb = 2
	}
	c.QueueParticle("Diluc", orb, def.Pyro, f+60)

	//actual skill cd starts immediately on first cast
	//times out after 4 seconds of not using
	//every hit applies pyro
	//apply attack speed

	d := c.Snapshot(
		"Searing Onslaught",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Pyro,
		25,
		skill[c.eCounter][c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	//check for c4 dmg increase
	if c.Base.Cons >= 4 {
		if c.Sim.Status("dilucc4") > 0 {
			d.Stats[def.DmgP] += 0.4
			c.Log.Debugw("diluc c4 adding dmg", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "final dmg", d.Stats[def.DmgP])
		}
	}

	c.QueueDmg(&d, f-5)

	//add a timer to activate c4
	if c.Base.Cons >= 4 {
		c.AddTask(func() {
			c.Sim.AddStatus("dilucc4", 120) //effect lasts 2 seconds
		}, "dilucc4", f+120) // 2seconds after cast
	}

	c.eCounter++
	if c.eCounter == 3 {
		//ability can go on cd now
		cd := 600 - (c.Sim.Frame() - c.eStartFrame)
		c.Log.Debugw("diluc skill going on cd", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "duration", cd)
		c.SetCD(def.ActionSkill, cd)
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

	// c.S.Status["dilucq"] = c.Sim.Frame() + 12*60
	c.Sim.AddStatus("dilucq", 720)
	f := c.ActionFrames(def.ActionBurst, p)

	d := c.Snapshot(
		"Dawn (Strike)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDiluc,
		def.StrikeTypeBlunt,
		def.Pyro,
		50,
		burstInitial[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll

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
	c.AddWeaponInfuse(def.WeaponInfusion{
		Key:    "diluc-fire-weapon",
		Ele:    def.Pyro,
		Tags:   []def.AttackTag{def.AttackTagNormal, def.AttackTagExtra, def.AttackTagPlunge},
		Expiry: c.Sim.Frame() + 852, //with a4
	})

	// add 20% pyro damage
	val := make([]float64, def.EndStatType)
	val[def.PyroP] = 0.2
	c.AddMod(def.CharStatMod{
		Key:    "diluc-fire-weapon",
		Amount: func(a def.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Sim.Frame() + 852,
	})

	c.Energy = 0
	c.SetCD(def.ActionBurst, 900)
	return f
}

func (c *char) Tick() {
	c.Tmpl.Tick()

	if c.eStarted {
		//check if 4 second has passed since last use
		if c.Sim.Frame()-c.eLastUse >= 240 {
			//if so, set ability to be on cd equal to 10s less started
			cd := 600 - (c.Sim.Frame() - c.eStartFrame)
			c.Log.Debugw("diluc skill going on cd", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "duration", cd, "last", c.eLastUse)
			c.SetCD(def.ActionSkill, cd)
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

func (c *char) ActionStam(a def.ActionType, p map[string]int) float64 {
	switch a {
	case def.ActionDash:
		return 18
	case def.ActionCharge:
		return 50
	default:
		c.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}
