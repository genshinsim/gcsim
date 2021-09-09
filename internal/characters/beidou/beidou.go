package beidou

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/shield"
)

func init() {
	core.RegisterCharFunc("beidou", NewChar)
}

type char struct {
	*character.Tmpl
	burstSnapshot core.Snapshot
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
	c.CharZone = core.ZoneLiyue

	c.burstProc()
	c.a4()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 23 //frames from keqing lib
		case 1:
			f = 43
		case 2:
			f = 68
		case 3:
			f = 44
		case 4:
			f = 68
		}
		atkspd := c.Stats[core.AtkSpd]
		if c.Core.Status.Duration("beidoua4") > 0 {
			atkspd += 0.15
		}
		f = int(float64(f) / (1 + atkspd))
		return f
	case core.ActionCharge:
		f := 35 //frames from keqing lib
		atkspd := c.Stats[core.AtkSpd]
		if c.Core.Status.Duration("beidoua4") > 0 {
			atkspd += 0.15
		}
		f = int(float64(f) / (1 + atkspd))
		return f
	case core.ActionSkill:
		return 41 //ok
	case core.ActionBurst:
		return 45 //ok
	default:
		c.Core.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

/**
Counterattacking with Tidecaller at the precise moment when the character is hit grants the maximum DMG Bonus.

Gain the following effects for 10s after unleashing Tidecaller with its maximum DMG Bonus:
• DMG dealt by Normal and Charged Attacks is increased by 15%. ATK SPD of Normal and Charged Attacks is increased by 15%.
• Greatly reduced delay before unleashing Charged Attacks.

c1
When Stormbreaker is used:
Creates a shield that absorbs up to 16% of Beidou's Max HP for 15s.
This shield absorbs Electro DMG 250% more effectively.

c2
Stormbreaker's arc lightning can jump to 2 additional targets.

c3
Within 10s of taking DMG, Beidou's Normal Attacks gain 20% additional Electro DMG.

c6
During the duration of Stormbreaker, the Electro RES of surrounding opponents is decreased by 15%.
**/

func (c *char) a4() {
	mod := make([]float64, core.EndStatType)
	mod[core.DmgP] = .15

	c.AddMod(core.CharStatMod{
		Key:    "beidou-a4",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a != core.AttackTagNormal && a != core.AttackTagExtra {
				return nil, false
			}
			if c.Core.Status.Duration("beidoua4") == 0 {
				return nil, false
			}
			return mod, true
		},
	})
}

func (c *char) c4() {
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.Index {
			return false
		}
		c.Core.Status.AddStatus("beidouc4", 600)
		c.Core.Log.Debugw("c4 triggered on damage", "frame", c.Core.F, "event", core.LogCharacterEvent, "expiry", c.Core.F+600)
		return false
	}, "beidouc4")

	mod := make([]float64, core.EndStatType)
	mod[core.DmgP] = .15

	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.Actor != c.Base.Name {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Core.Status.Duration("beidouc4") == 0 {
			return false
		}

		c.Core.Log.Debugw("c4 proc'd on attack", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index)
		d := c.Snapshot(
			"Beidou C4",
			core.AttackTagNone,
			core.ICDTagElementalBurst,
			core.ICDGroupDefault,
			core.StrikeTypeBlunt,
			core.Electro,
			25,
			0.2,
		)
		c.QueueDmg(&d, 1)
		return false
	}, "beidou-c4")

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
	counter := p["counter"]
	f := c.ActionFrames(core.ActionSkill, p)
	//0 for base dmg, 1 for 1x bonus, 2 for max bonus
	if counter >= 2 {
		counter = 2
		c.Core.Status.AddStatus("beidoua4", 600)
	}

	d := c.Snapshot(
		"Tidecaller (E)",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Electro,
		50,
		skillbase[c.TalentLvlSkill()]+skillbonus[c.TalentLvlSkill()]*float64(counter),
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f-1)

	//2 if no hit, 3 if 1 hit, 4 if perfect
	c.QueueParticle("beidou", 2+counter, core.Electro, 100)

	if counter > 0 {
		//add shield
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldBeidouThunderShield,
			HP:         shieldPer[c.TalentLvlSkill()]*c.HPMax + shieldBase[c.TalentLvlSkill()],
			Ele:        core.Electro,
			Expires:    c.Core.F + 900, //15 sec
		})
	}

	c.SetCD(core.ActionSkill, 450)
	return f
}

func (c *char) Burst(p map[string]int) int {
	if c.Energy < c.EnergyMax {
		c.Core.Log.Debugw("burst insufficient energy; skipping", "frame", c.Core.F, "event", core.LogCharacterEvent, "character", c.Base.Name)
		return 0
	}

	f := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Stormbreaker (Q)",
		core.AttackTagElementalBurst,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		100,
		burstonhit[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f-1)

	c.Core.Status.AddStatus("beidouburst", 900)
	c.burstSnapshot = c.Snapshot(
		"Stormbreaker Proc (Q)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		burstproc[c.TalentLvlBurst()],
	)

	if c.Base.Cons >= 1 {
		//create a shield
		c.Core.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: core.ShieldBeidouThunderShield,
			HP:         .16 * c.HPMax,
			Ele:        core.Electro,
			Expires:    c.Core.F + 900, //15 sec
		})
	}

	if c.Base.Cons == 6 {
		for _, t := range c.Core.Targets {
			t.AddResMod("beidouc6", core.ResistMod{
				Duration: 900, //10 seconds
				Ele:      core.Electro,
				Value:    -0.1,
			})
		}
	}

	//c.Energy = 0  forcing every character to comsume energy after burts in the energy.go to make my life easier
	c.ConsumeEnergy(0, 0) //at 0,0 value acts the same as c.Energy = 0
	c.SetCD(core.ActionBurst, 1200)
	return f
}

func (c *char) burstProc() {
	icd := 0
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Core.Status.Duration("beidouburst") == 0 {
			return false
		}
		if icd > c.Core.F {
			c.Core.Log.Debugw("beidou Q (active) on icd", "frame", c.Core.F, "event", core.LogCharacterEvent)
			return false
		}

		d := c.burstSnapshot.Clone()
		//on hit we have to chain
		d.OnHitCallback = c.chainQ(t.Index(), c.Core.F, 1)

		c.Core.Log.Debugw("beidou Q proc'd", "frame", c.Core.F, "event", core.LogCharacterEvent, "actor", ds.Actor, "attack tag", ds.AttackTag)
		c.QueueDmg(&d, 1)

		icd = c.Core.F + 60 // once per second
		return false
	}, "beidou-burst")
}

func (c *char) chainQ(index int, src int, count int) func(t core.Target) {
	if c.Base.Cons > 1 && count == 5 {
		return nil
	}
	if c.Base.Cons < 2 && count == 3 {
		return nil
	}
	//check number of targets, if target < 2 then no bouncing

	//figure out the next target
	l := len(c.Core.Targets)
	if l < 2 {
		return nil
	}
	index++
	if index >= l {
		index = 0
	}

	//trigger dmg based on a clone of d
	return func(next core.Target) {
		// log.Printf("hit target %v, frame %v, done proc %v, queuing next index: %v\n", next.Index(), c.Core.F, count, index)
		d := c.burstSnapshot.Clone()
		d.Targets = index
		d.SourceFrame = c.Core.F
		d.OnHitCallback = c.chainQ(index, src, count+1)
		c.QueueDmg(&d, 1)
	}
}
