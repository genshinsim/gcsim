package xingqiu

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl
	numSwords    int
	nextRegen    bool
	burstCounter int
	burstTickSrc int
	// burstICDResetTimer int //if c.S.F > this, then reset counter to = 0
	orbitalActive bool
	burstSwordICD int
}

func init() {
	core.RegisterCharFunc(core.Xingqiu, NewChar)
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5
	c.CharZone = core.ZoneLiyue

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.burstStateHook()

	c.AddMod(core.CharStatMod{
		Key: "a4",
		Amount: func() ([]float64, bool) {
			a4 := make([]float64, core.EndStatType)
			a4[core.HydroP] = 0.2
			return a4, true
		},
		Expiry: -1,
	})
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
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
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 63, 63
	case core.ActionSkill:
		return 77, 77 //should be 82
	case core.ActionBurst:
		return 39, 39 //ok
	default:
		c.Core.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

var hitmarks = [][]int{{8}, {24}, {24, 43}, {36}, {43, 78}}

func (c *char) Attack(p map[string]int) (int, int) {
	//apply attack speed
	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Abil = fmt.Sprintf("Normal %v", c.NormalCounter)
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), hitmarks[c.NormalCounter][i], hitmarks[c.NormalCounter][i])
	}

	//add a 75 frame attackcounter reset
	c.AdvanceNormalIndex()
	//return animation cd
	//this also depends on which hit in the chain this is
	return f, a
}

func (c *char) orbitalfunc(src int) func() {
	return func() {
		c.Core.Log.NewEvent("orbital checking tick", core.LogCharacterEvent, c.Index, "expiry", c.Core.Status.Duration("xqorb"), "src", src)
		if c.Core.Status.Duration("xqorb") == 0 {
			c.orbitalActive = false
			return
		}

		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Xingqiu Skill (Orbital)",
			AttackTag:  core.AttackTagNone,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Hydro,
			Durability: 25,
		}
		c.Core.Log.NewEvent("orbital ticked", core.LogCharacterEvent, c.Index, "next expected tick", c.Core.F+150, "expiry", c.Core.Status.Duration("xqorb"), "src", src)

		//queue up next instance
		c.AddTask(c.orbitalfunc(src), "xq-skill-orbital", 135)

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), -1, 1)
	}
}

func (c *char) applyOrbital(duration int) {
	f := c.Core.F
	c.Core.Log.NewEvent("Applying orbital", core.LogCharacterEvent, c.Index, "current status", c.Core.Status.Duration("xqorb"))
	//check if orbitals already active, if active extend duration
	//other wise start first tick func
	if !c.orbitalActive {
		c.AddTask(c.orbitalfunc(f), "xq-skill-orbital", 14)
		c.orbitalActive = true
		c.Core.Log.NewEvent("orbital applied", core.LogCharacterEvent, c.Index, "expected end", f+900, "next expected tick", f+40)
	}

	c.Core.Status.AddStatus("xqorb", duration)
	c.Core.Log.NewEvent("orbital duration extended", core.LogCharacterEvent, c.Index, "new expiry", c.Core.Status.Duration("xqorb"))
}

var rainscreenDelay = [2]int{19, 35}

func (c *char) Skill(p map[string]int) (int, int) {
	//applies wet to self 30 frame after cast, sword applies wet every 2.5seconds, so should be 7 seconds

	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Guhua Sword: Fatal Rainscreen",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
	}

	for i, v := range rainscreen {
		ai.Mult = v[c.TalentLvlSkill()]
		if c.Base.Cons >= 4 {
			//check if ult is up, if so increase multiplier
			if c.Core.Status.Duration("xqburst") > 0 {
				ai.Mult = ai.Mult * 1.5
			}
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), rainscreenDelay[i], rainscreenDelay[i])
	}

	// Orbitals spawn in 1 frame before the second hit connects going by the "Wet" text
	c.AddTask(func() {
		orbital := p["orbital"]
		if orbital == 1 {
			c.applyOrbital(15 * 60)
		}
	}, "xingqiu-spawn-orbitals", 34)

	c.QueueParticle(c.Base.Key.String(), 5, core.Hydro, 100)

	//should last 15s, cd 21s
	c.SetCD(core.ActionSkill, 21*60)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	//apply hydro every 3rd hit
	//triggered on normal attack
	//also applies hydro on cast if p=1

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
	c.Core.Log.NewEvent("Xingqiu burst activated", core.LogCharacterEvent, c.Index, "expiry", c.Core.F+dur)

	orbital := p["orbital"]

	if orbital == 1 {
		c.applyOrbital(dur)
	}

	c.burstCounter = 0
	c.numSwords = 2

	// c.CD[combat.BurstCD] = c.S.F + 20*60
	c.SetCDWithDelay(core.ActionBurst, 20*60, 7)
	c.ConsumeEnergy(7)
	return f, a
}
