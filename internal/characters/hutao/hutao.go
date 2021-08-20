package hutao

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("hutao", NewChar)
}

type char struct {
	*character.Tmpl
	paraParticleICD int
	// chargeICDCounter   int
	// chargeCounterReset int
	ppBonus    float64
	tickActive bool
	c6icd      int
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
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 6

	c.ppHook()
	c.onExitField()
	c.a4()

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

/**
[11:32 PM] sakuno | yanfei is my new maid: @gimmeabreak
https://www.youtube.com/watch?v=3aCiH2U4BjY

framecounts for 7 attempts of N2CJ (no hitlag):
83, 85, 88, 89, 77, 82, 84

first 3 not from the uploaded recording (as a n1cd player i cud barely pull it off :monkaS: )
YouTube
**/

//var normalFrames = []int{13, 16, 25, 36, 44, 39}               // from kqm lib
var normalFrames = []int{10, 13, 22, 33, 41, 36} // from kqm lib, -3 for hit lag
//var dmgFrame = [][]int{{13}, {16}, {25}, {36}, {26, 44}, {39}} // from kqm lib
var dmgFrame = [][]int{{10}, {13}, {22}, {33}, {23, 41}, {36}} // from kqm lib - 3 for hit lag

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := normalFrames[c.NormalCounter]
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionCharge:
		return 9 //rough.. 11, -2 for hit lag
	case core.ActionSkill:
		return 42 // from kqm lib
	case core.ActionBurst:
		return 130 // from kqm lib
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Core.Status.Duration("paramita") > 0 && c.Base.Cons >= 1 {
			return 0
		}
		return 25
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) a4() {
	val := make([]float64, core.EndStatType)
	val[core.PyroP] = 0.33
	c.AddMod(core.CharStatMod{
		Key:    "hutao-a4",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if c.Core.Status.Duration("paramita") == 0 {
				return nil, false
			}
			if c.HPCurrent/c.HPMax <= 0.5 {
				return val, true
			}
			return nil, false
		},
	})
}

func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		c.checkc6()
		return false
	}, "hutao-c6")
}

func (c *char) checkc6() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Core.F < c.c6icd && c.c6icd != 0 {
		return
	}
	//check if hp less than 25%
	if c.HPCurrent/c.HPMax > .25 {
		return
	}
	//if dead, revive back to 1 hp
	if c.HPCurrent == -1 {
		c.HPCurrent = 1
	}
	//increase crit rate to 100%
	val := make([]float64, core.EndStatType)
	val[core.CR] = 1
	c.AddMod(core.CharStatMod{
		Key:    "hutao-c6",
		Amount: func(a core.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Core.F + 600,
	})

	c.c6icd = c.Core.F + 3600
}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(core.ActionAttack, p)
	hits := len(attack[c.NormalCounter])
	//check for particles
	c.ppParticles()

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

	for i := 0; i < hits; i++ {
		x := d.Clone()
		x.Mult = attack[c.NormalCounter][i][c.TalentLvlAttack()]
		c.QueueDmg(&x, dmgFrame[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {

	f := c.ActionFrames(core.ActionCharge, p)

	if c.Core.Status.Duration("paramita") > 0 {
		//[3:56 PM] Isu: My theory is that since E changes attack animations, it was coded
		//to not expire during any attack animation to simply avoid the case of potentially
		//trying to change animations mid-attack, but not sure how to fully test that
		//[4:41 PM] jstern25| ₼WHO_SUPREMACY: this mostly checks out
		//her e can't expire during q as well
		if f > c.Core.Status.Duration("paramita") {
			c.Core.Status.AddStatus("paramita", f)
			// c.S.Status["paramita"] = c.Core.F + f //extend this to barely cover the burst
		}

		c.applyBB()
		//charge land 182, tick 432, charge 632, tick 675
		//charge land 250, tick 501, charge 712, tick 748

		//e cast at 123, animation ended 136 should end at 664 if from cast or 676 if from animation end, tick at 748 still buffed?
	}

	//check for particles
	//TODO: assuming charge can generate particles as well
	c.ppParticles()

	d := c.Snapshot(
		"Charge Attack",
		core.AttackTagExtra,
		core.ICDTagExtraAttack,
		core.ICDGroupPole,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		charge[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-5)

	return f
}

func (c *char) ppParticles() {
	if c.Core.Status.Duration("paramita") > 0 {
		if c.paraParticleICD < c.Core.F {
			c.paraParticleICD = c.Core.F + 300 //5 seconds
			count := 2
			if c.Core.Rand.Float64() < 0.5 {
				count = 3
			}
			c.QueueParticle("Hutao", count, core.Pyro, dmgFrame[c.NormalCounter][0])
		}
	}
}

func (c *char) applyBB() {
	c.Core.Log.Debugw("Applying Blood Blossom", "frame", c.Core.F, "event", core.LogCharacterEvent, "current dur", c.Core.Status.Duration("htbb"))
	//check if blood blossom already active, if active extend duration by 8 second
	//other wise start first tick func
	if !c.tickActive {
		//TODO: does BB tick immediately on first application?
		c.AddTask(c.bbtickfunc(c.Core.F), "bb", 240)
		c.tickActive = true
		c.Core.Log.Debugw("Blood Blossom applied", "frame", c.Core.F, "event", core.LogCharacterEvent, "expected end", c.Core.F+570, "next expected tick", c.Core.F+240)
	}
	// c.CD["bb"] = c.Core.F + 570 //TODO: no idea how accurate this is, does this screw up the ticks?
	c.Core.Status.AddStatus("htbb", 570)
	c.Core.Log.Debugw("Blood Blossom duration extended", "frame", c.Core.F, "event", core.LogCharacterEvent, "new expiry", c.Core.Status.Duration("htbb"))
}

func (c *char) bbtickfunc(src int) func() {
	return func() {
		c.Core.Log.Debugw("Blood Blossom checking for tick", "frame", c.Core.F, "event", core.LogCharacterEvent, "cd", c.Core.Status.Duration("htbb"), "src", src)
		if c.Core.Status.Duration("htbb") == 0 {
			c.tickActive = false
			return
		}
		//queue up one damage instance
		d := c.Snapshot(
			"Blood Blossom",
			core.AttackTagElementalArt,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Pyro,
			25,
			bb[c.TalentLvlSkill()],
		)

		//if cons 2, add flat dmg
		if c.Base.Cons >= 2 {
			d.FlatDmg += c.HPMax * 0.1
		}
		c.Core.Combat.ApplyDamage(&d)
		c.Core.Log.Debugw("Blood Blossom ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "next expected tick", c.Core.F+240, "dur", c.Core.Status.Duration("htbb"), "src", src)
		//only queue if next tick buff will be active still
		// if c.Core.F+240 > c.CD["bb"] {
		// 	return
		// }
		//queue up next instance
		c.AddTask(c.bbtickfunc(src), "bb", 240)

	}
}

func (c *char) Skill(p map[string]int) int {
	//increase based on hp at cast time
	//drains hp
	c.Core.Status.AddStatus("paramita", 520+20) //to account for animation
	c.Core.Log.Debugw("Paramita acivated", "frame", c.Core.F, "event", core.LogCharacterEvent, "expiry", c.Core.F+540+20)
	//figure out atk buff
	c.ppBonus = ppatk[c.TalentLvlSkill()] * c.HPMax
	max := (c.Base.Atk + c.Weapon.Atk) * 4
	if c.ppBonus > max {
		c.ppBonus = max
	}

	//remove some hp
	c.HPCurrent = 0.7 * c.HPCurrent
	c.checkc6()

	c.SetCD(core.ActionSkill, 960)
	return c.ActionFrames(core.ActionSkill, p)
}

func (c *char) ppHook() {
	val := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key:    "hutao-paramita",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if c.Core.Status.Duration("paramita") == 0 {
				return nil, false
			}
			val[core.ATK] = c.ppBonus
			return val, true
		},
	})
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		c.Core.Status.DeleteStatus("paramita")
		return false
	}, "hutao-exit")
}

func (c *char) Burst(p map[string]int) int {
	low := (c.HPCurrent / c.HPMax) <= 0.5
	mult := burst[c.TalentLvlBurst()]
	regen := regen[c.TalentLvlBurst()]
	if low {
		mult = burstLow[c.TalentLvlBurst()]
		regen = regenLow[c.TalentLvlBurst()]
	}
	targets := p["targets"]
	//regen for p+1 targets, max at 5; if not specified then p = 1
	count := 1
	if targets > 0 {
		count = targets
	}
	if count > 5 {
		count = 5
	}
	c.HPCurrent += c.HPMax * float64(count) * regen

	f := c.ActionFrames(core.ActionBurst, p)

	//[2:28 PM] Aluminum | Harbinger of Jank: I think the idea is that PP won't fall off before dmg hits, but other buffs aren't snapshot
	//[2:29 PM] Isu: yes, what Aluminum said. PP can't expire during the burst animation, but any other buff can
	if f > c.Core.Status.Duration("paramita") && c.Core.Status.Duration("paramita") > 0 {
		c.Core.Status.AddStatus("paramita", f) //extend this to barely cover the burst
	}

	if c.Core.Status.Duration("paramita") > 0 && c.Base.Cons >= 2 {
		c.applyBB()
	}

	c.AddTask(func() {
		//TODO: apparently damage is based on stats on contact, not at cast
		d := c.Snapshot(
			"Spirit Soother",
			core.AttackTagElementalBurst,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Pyro,
			50,
			mult,
		)
		d.Targets = core.TargetAll
		c.Core.Combat.ApplyDamage(&d)
	}, "Hutao Burst", f-5) //random 5 frame

	c.Energy = 0
	c.SetCD(core.ActionBurst, 900)
	return f
}

func (c *char) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	if c.Core.Status.Duration("paramita") > 0 {
		switch ds.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
		default:
			return ds
		}
		ds.Element = core.Pyro
	}
	return ds
}
