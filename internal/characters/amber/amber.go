package amber

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("amber", NewChar)
}

type char struct {
	*character.Tmpl
	bunnies      []bunny
	eCharge      int
	eChargeMax   int
	eNextRecover int
	eTickSrc     int
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
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.eChargeMax = 1
	if c.Base.Cons >= 4 {
		c.eChargeMax = 2
	}
	c.eCharge = c.eChargeMax

	if c.Base.Cons >= 2 {
		c.overloadExplode()
	}

	c.bunnies = make([]bunny, 0, 2)

	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 15 //frames from keqing lib
		case 1:
			f = 33 - 15
		case 2:
			f = 72 - 33
		case 3:
			f = 113 - 72
		case 4:
			f = 155 - 113
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionAim:
		return 94 //kqm
	case core.ActionBurst:
		return 74 //swap canceled
	case core.ActionSkill:
		return 35 //no cancel
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) Attack(p map[string]int) int {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, travel+f)

	if c.Base.Cons >= 1 {
		x := d.Clone()
		x.Mult = .2 * d.Mult
		c.QueueDmg(&x, travel+f)
	}

	c.AdvanceNormalIndex()

	return f
}

func (c *char) Aimed(p map[string]int) int {
	f := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	b := p["bunny"]

	if c.Base.Cons >= 2 && b != 0 {
		//explode the first bunny
		c.AddTask(func() {
			c.manualExplode()
		}, "bunny", travel+f)

		//also don't do any dmg since we're shooting at bunny

		return f
	}

	d := c.Snapshot(
		"Aim (Charged)",
		core.AttackTagExtra,
		core.ICDTagExtraAttack,
		core.ICDGroupAmber,
		core.StrikeTypePierce,
		core.Pyro,
		50,
		aim[c.TalentLvlAttack()],
	)

	d.HitWeakPoint = true
	d.AnimationFrames = f

	//add 15% since 360noscope

	c.AddTask(func() {
		val := make([]float64, core.EndStatType)
		val[core.ATKP] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "a2",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return val, true
			},
			Expiry: c.Core.F + 600,
		})
	}, "aim", f+travel)

	c.QueueDmg(&d, travel+f)

	return f
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)
	hold := p["hold"]

	c.AddTask(func() {
		c.makeBunny()
	}, "new-bunny", f+hold)

	c.overloadExplode()

	if c.Base.Cons < 4 {
		c.SetCD(core.ActionSkill, 900)
		return f + hold
	}

	switch c.eCharge {
	case c.eChargeMax:
		c.Log.Debugw("amber bunny at max charge, queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "recover at", c.Core.F+721)
		c.eNextRecover = c.Core.F + 721
		c.AddTask(c.recoverCharge(c.Core.F), "charge", 720)
		c.eTickSrc = c.Core.F
	case 1:
		c.SetCD(core.ActionSkill, c.eNextRecover)
	}
	c.eCharge--

	return f + hold
}

func (c *char) recoverCharge(src int) func() {
	return func() {
		if c.eTickSrc != src {
			c.Log.Debugw("amber bunny recovery function ignored, src diff", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src, "new src", c.eTickSrc)
			return
		}
		c.eCharge++
		c.Log.Debugw("amber bunny recovering a charge", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src, "total charge", c.eCharge)
		c.SetCD(core.ActionSkill, 0)
		if c.eCharge >= c.eChargeMax {
			//fully charged
			return
		}
		//other wise restore another charge
		c.Log.Debugw("amber bunny queuing next recovery", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src, "recover at", c.Core.F+720)
		c.eNextRecover = c.Core.F + 721
		c.AddTask(c.recoverCharge(src), "charge", 720)

	}
}

type bunny struct {
	ds  core.Snapshot
	src int
}

func (c *char) makeBunny() {
	b := bunny{}
	b.src = c.Core.F
	b.ds = c.Snapshot(
		"Baron Bunny",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Pyro,
		50,
		bunnyExplode[c.TalentLvlSkill()],
	)
	b.ds.Targets = core.TargetAll

	c.bunnies = append(c.bunnies, b)

	//ondeath explodes
	//duration is 8.2 sec
	c.AddTask(func() {
		c.explode(b.src)
	}, "bunny", 492)
}

func (c *char) explode(src int) {
	n := 0
	c.Log.Debugw("amber exploding bunny", "frame", c.Core.F, "event", core.LogCharacterEvent, "src", src)
	for _, v := range c.bunnies {
		if v.src == src {

			c.QueueDmg(&v.ds, 1)
			//4 orbs
			c.QueueParticle("amber", 4, core.Pyro, 100)
		} else {
			c.bunnies[n] = v
			n++
		}
	}

	c.bunnies = c.bunnies[:n]
}

func (c *char) manualExplode() {
	if len(c.bunnies) > 0 {
		ds := c.bunnies[0].ds
		ds.Mult = ds.Mult + 2
		c.QueueDmg(&ds, 1)
		c.QueueParticle("amber", 4, core.Pyro, 100)
	}
	c.bunnies = c.bunnies[1:]
}

func (c *char) overloadExplode() {
	//explode all bunnies on overload
	c.Core.Events.Subscribe(core.OnTransReaction, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if len(c.bunnies) == 0 {
			return false
		}
		//TODO: only amber trigger?
		if ds.ActorIndex != c.Index {
			return false
		}
		//TODO: does it have to be charge shot trigger only??
		if ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if ds.ReactionType == core.Overload {
			for _, v := range c.bunnies {
				ds := v.ds
				ds.Mult = ds.Mult + 2
				c.QueueDmg(&ds, 1)
				c.QueueParticle("amber", 4, core.Pyro, 100)
			}
			c.bunnies = make([]bunny, 0, 2)
		}
		return false
	}, "bunnyer-overload")
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)

	//2sec duration, tick every .4 sec in zone 1
	//2sec duration, tick every .6 sec in zone 2
	//2sec duration, tick every .2 sec in zone 3
	d := c.Snapshot(
		"Fiery Rain",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupAmber,
		core.StrikeTypePierce,
		core.Pyro,
		25,
		burstTick[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	d.Stats[core.CR] += 0.1 // a2

	for i := f + 24; i < 120+f; i += 24 {
		x := d.Clone()
		c.QueueDmg(&x, i)
	}

	for i := f + 36; i < 120+f; i += 36 {
		x := d.Clone()
		c.QueueDmg(&x, i)
	}

	for i := f + 12; i < 120+f; i += 12 {
		x := d.Clone()
		c.QueueDmg(&x, i)
	}

	if c.Base.Cons == 6 {
		for _, active := range c.Core.Chars {
			val := make([]float64, core.EndStatType)
			val[core.ATKP] = 0.15
			active.AddMod(core.CharStatMod{
				Key:    "amber-c6",
				Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
				Expiry: c.Core.F + 900,
			})
			c.Log.Debugw("c6 - adding atk %", "frame", c.Core.F, "event", core.LogCharacterEvent, "character", c.Name())
		}
	}

	c.Energy = 0
	c.SetCD(core.ActionBurst, 720)
	return f
}
