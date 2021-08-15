package xingqiu

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

type char struct {
	*character.Tmpl
	numSwords          int
	nextRegen          bool
	burstCounter       int
	burstICDResetTimer int //if c.S.F > this, then reset counter to = 0
	orbitalActive      bool
	burstSwordICD      int
}

func init() {
	core.RegisterCharFunc("xingqiu", NewChar)
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
	c.Weapon.Class = core.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5

	a4 := make([]float64, core.EndStatType)
	a4[core.HydroP] = 0.2
	c.AddMod(core.CharStatMod{
		Key: "a4",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return a4, true
		},
		Expiry: -1,
	})
	c.burstHook()

	/** c6
	Activating 2 of Guhua Sword: Raincutter's sword rain attacks greatly increases the DMG of the third.
	Xingqiu regenerates 3 Energy when sword rain attacks hit opponents.
	**/

	return &c, nil
}

var delay = [][]int{{8}, {24}, {24, 43}, {36}, {43, 78}}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 9
		case 1:
			f = 25
		case 2:
			f = 44
		case 3:
			f = 37
		case 4:
			f = 79
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 63
	case core.ActionSkill:
		return 77 //should be 82
	case core.ActionBurst:
		return 39 //ok
	default:
		c.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0
	}
}

func (c *char) Attack(p map[string]int) int {
	//apply attack speed
	f := c.ActionFrames(core.ActionAttack, p)

	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		0,
	)

	for i, mult := range attack[c.NormalCounter] {
		x := d.Clone()
		x.Mult = mult[c.TalentLvlAttack()]
		c.QueueDmg(&x, delay[c.NormalCounter][i])
	}

	//add a 75 frame attackcounter reset
	c.AdvanceNormalIndex()
	//return animation cd
	//this also depends on which hit in the chain this is
	return f
}

func (c *char) orbitalfunc(src int) func() {
	return func() {
		c.Log.Debugw("orbital checking tick", "frame", c.Core.F, "event", core.LogCharacterEvent, "expiry", c.Core.Status.Duration("xqorb"), "src", src)
		if c.Core.Status.Duration("xqorb") == 0 {
			c.orbitalActive = false
			return
		}
		//queue up one damage instance
		d := c.Snapshot(
			"Xingqiu Skill (Orbital)",
			core.AttackTagNone,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Hydro,
			25,
			0,
		)
		d.Targets = core.TargetAll

		c.QueueDmg(&d, 1)
		c.Log.Debugw("orbital ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "next expected tick", c.Core.F+150, "expiry", c.Core.Status.Duration("xqorb"), "src", src)
		//queue up next instance
		c.AddTask(c.orbitalfunc(src), "xq-skill-orbital", 150)
	}
}

func (c *char) applyOrbital() {
	f := c.Core.F
	c.Log.Debugw("Applying orbital", "frame", f, "event", core.LogCharacterEvent, "current status", c.Core.Status.Duration("xqorb"))
	//check if blood blossom already active, if active extend duration by 8 second
	//other wise start first tick func
	if !c.orbitalActive {
		//TODO: does BB tick immediately on first application?
		c.AddTask(c.orbitalfunc(f), "xq-skill-orbital", 40)
		c.orbitalActive = true
		c.Log.Debugw("orbital applied", "frame", f, "event", core.LogCharacterEvent, "expected end", f+900, "next expected tick", f+40)
	}
	c.Core.Status.AddStatus("xqorb", 900)
	c.Log.Debugw("orbital duration extended", "frame", f, "event", core.LogCharacterEvent, "new expiry", c.Core.Status.Duration("xqorb"))
}

func (c *char) Skill(p map[string]int) int {
	//applies wet to self 30 frame after cast, sword applies wet every 2.5seconds, so should be 7 seconds
	orbital := p["orbital"]
	if orbital == 1 {
		c.applyOrbital()
	}

	f := c.ActionFrames(core.ActionSkill, p)

	d := c.Snapshot(
		"Guhua Sword: Fatal Rainscreen",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Hydro,
		25,
		rainscreen[0][c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll
	if c.Base.Cons >= 4 {
		//check if ult is up, if so increase multiplier
		if c.Core.Status.Duration("xqburst") > 0 {
			d.Mult = d.Mult * 1.5
		}
	}
	d2 := d.Clone()
	d2.Mult = rainscreen[1][c.TalentLvlSkill()]

	c.QueueDmg(&d, 19)
	c.QueueDmg(&d2, 39)

	c.QueueParticle(c.Base.Name, 5, core.Hydro, 100)

	//should last 15s, cd 21s
	c.SetCD(core.ActionSkill, 21*60)
	return f
}

func (c *char) burstHook() {
	c.Core.Events.Subscribe(core.PostAttack, func(args ...interface{}) bool {
		//check if buff is up
		if c.Core.Status.Duration("xqburst") <= 0 {
			return false
		}
		//check if off ICD
		if c.burstSwordICD > c.Core.F {
			return false
		}

		const delay = 5 //wait 5 frames into attack animation

		//trigger swords, only first sword applies hydro
		for i := 0; i < c.numSwords; i++ {

			wave := i

			d := c.Snapshot(
				"Guhua Sword: Raincutter",
				core.AttackTagElementalBurst,
				core.ICDTagElementalBurst,
				core.ICDGroupDefault,
				core.StrikeTypePierce,
				core.Hydro,
				25,
				burst[c.TalentLvlBurst()],
			)
			d.Targets = 0 //only hit main target
			d.OnHitCallback = func(t core.Target) {
				//check energy
				if c.nextRegen && wave == 0 {
					c.AddEnergy(3)
				}
				//check c2
				if c.Base.Cons >= 2 {
					t.AddResMod("xingqiu-c2", core.ResistMod{
						Ele:      core.Hydro,
						Value:    -0.15,
						Duration: 4 * 60,
					})
				}
			}

			c.QueueDmg(&d, delay+20+i)

			c.burstCounter++
		}

		//figure out next wave # of swords
		switch c.numSwords {
		case 2:
			c.numSwords = 3
			c.nextRegen = false
		case 3:
			if c.Base.Cons == 6 {
				c.numSwords = 5
				c.nextRegen = true
			} else {
				c.numSwords = 2
				c.nextRegen = false
			}
		case 5:
			c.numSwords = 2
			c.nextRegen = false
		}

		//estimated 1 second ICD
		c.burstSwordICD = c.Core.F + 60

		return false
	}, "xq-burst")
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionBurst, p)
	//apply hydro every 3rd hit
	//triggered on normal attack
	//also applies hydro on cast if p=1
	orbital := p["orbital"]

	if orbital == 1 {
		c.applyOrbital()
	}
	//how we doing that?? trigger 0 dmg?

	/**
	The number of Hydro Swords summoned per wave follows a specific pattern, usually alternating between 2 and 3 swords.
	At C6, this is upgraded and follows a pattern of 2 → 3 → 5… which then repeats.

	There is an approximately 1 second interval between summoned Hydro Sword waves, so that means a theoretical maximum of 15 or 18 waves.

	Each wave of Hydro Swords is capable of applying one (1) source of Hydro status, and each individual sword is capable of getting a crit.
	**/

	/** c2
	Extends the duration of Guhua Sword: Raincutter by 3s.
	Decreases the Hydro RES of opponents hit by sword rain attacks by 15% for 4s.
	**/
	dur := 15
	if c.Base.Cons >= 2 {
		dur += 3
	}
	dur = dur * 60
	c.Core.Status.AddStatus("xqburst", dur)
	c.Log.Debugw("Xingqiu burst activated", "frame", c.Core.F, "event", core.LogCharacterEvent, "expiry", c.Core.F+dur)

	c.burstCounter = 0
	c.numSwords = 2

	// c.CD[combat.BurstCD] = c.S.F + 20*60
	c.SetCD(core.ActionBurst, 20*60)
	c.Energy = 0
	return f
}
