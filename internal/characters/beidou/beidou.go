package beidou

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"github.com/genshinsim/gsim/pkg/shield"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("beidou", NewChar)
}

type char struct {
	*character.Tmpl
	burstSnapshot def.Snapshot
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = def.WeaponClassClaymore
	c.NormalHitNum = 5

	c.burstProc()
	c.a4()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
}

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
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
		atkspd := c.Stats[def.AtkSpd]
		if c.Sim.Status("beidoua4") > 0 {
			atkspd += 0.15
		}
		f = int(float64(f) / (1 + atkspd))
		return f
	case def.ActionCharge:
		f := 35 //frames from keqing lib
		atkspd := c.Stats[def.AtkSpd]
		if c.Sim.Status("beidoua4") > 0 {
			atkspd += 0.15
		}
		f = int(float64(f) / (1 + atkspd))
		return f
	case def.ActionSkill:
		return 41 //ok
	case def.ActionBurst:
		return 45 //ok
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
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
	mod := make([]float64, def.EndStatType)
	mod[def.DmgP] = .15

	c.AddMod(def.CharStatMod{
		Key:    "beidou-a4",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if a != def.AttackTagNormal && a != def.AttackTagExtra {
				return nil, false
			}
			if c.Sim.Status("beidoua4") == 0 {
				return nil, false
			}
			return mod, true
		},
	})
}

func (c *char) c4() {
	c.Sim.AddOnHurt(func(s def.Sim) {
		c.Log.Debugw("c4 triggered on damage", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "expiry", c.Sim.Frame()+600)
		c.Sim.AddStatus("beidouc4", 600)
	})

	mod := make([]float64, def.EndStatType)
	mod[def.DmgP] = .15

	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.Actor != c.Base.Name {
			return
		}
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		if c.Sim.Status("beidouc4") == 0 {
			return
		}

		c.Log.Debugw("c4 proc'd on attack", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "char", c.Index)
		d := c.Snapshot(
			"Beidou C4",
			def.AttackTagNone,
			def.ICDTagElementalBurst,
			def.ICDGroupDefault,
			def.StrikeTypeBlunt,
			def.Electro,
			25,
			0.2,
		)
		c.QueueDmg(&d, 1)
	}, "beidou-c4")
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
	counter := p["counter"]
	f := c.ActionFrames(def.ActionSkill, p)
	//0 for base dmg, 1 for 1x bonus, 2 for max bonus
	if counter >= 2 {
		counter = 2
		c.Sim.AddStatus("beidoua4", 600)
	}

	d := c.Snapshot(
		"Tidecaller (E)",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Electro,
		50,
		skillbase[c.TalentLvlSkill()]+skillbonus[c.TalentLvlSkill()]*float64(counter),
	)
	d.Targets = def.TargetAll
	c.QueueDmg(&d, f-1)

	//2 if no hit, 3 if 1 hit, 4 if perfect
	c.QueueParticle("beidou", 2+counter, def.Electro, 100)

	if counter > 0 {
		//add shield
		c.Sim.AddShield(&shield.Tmpl{
			Src:        c.Sim.Frame(),
			ShieldType: def.ShieldBeidouThunderShield,
			HP:         shieldPer[c.TalentLvlSkill()]*c.HPMax + shieldBase[c.TalentLvlSkill()],
			Ele:        def.Electro,
			Expires:    c.Sim.Frame() + 900, //15 sec
		})
	}

	c.SetCD(def.ActionSkill, 450)
	return f
}

func (c *char) Burst(p map[string]int) int {
	if c.Energy < c.EnergyMax {
		c.Log.Debugw("burst insufficient energy; skipping", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "character", c.Base.Name)
		return 0
	}

	f := c.ActionFrames(def.ActionSkill, p)
	d := c.Snapshot(
		"Stormbreaker (Q)",
		def.AttackTagElementalBurst,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Electro,
		100,
		burstonhit[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll
	c.QueueDmg(&d, f-1)

	c.Sim.AddStatus("beidouburst", 900)
	c.burstSnapshot = c.Snapshot(
		"Stormbreaker Proc (Q)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Electro,
		25,
		burstproc[c.TalentLvlBurst()],
	)

	if c.Base.Cons >= 1 {
		//create a shield
		c.Sim.AddShield(&shield.Tmpl{
			Src:        c.Sim.Frame(),
			ShieldType: def.ShieldBeidouThunderShield,
			HP:         .16 * c.HPMax,
			Ele:        def.Electro,
			Expires:    c.Sim.Frame() + 900, //15 sec
		})
	}

	if c.Base.Cons == 6 {
		for _, t := range c.Sim.Targets() {
			t.AddResMod("beidouc6", def.ResistMod{
				Duration: 900, //10 seconds
				Ele:      def.Electro,
				Value:    -0.1,
			})
		}
	}

	c.Energy = 0
	c.SetCD(def.ActionBurst, 1200)
	return f
}

func (c *char) burstProc() {
	icd := 0
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if ds.AttackTag != def.AttackTagNormal && ds.AttackTag != def.AttackTagExtra {
			return
		}
		if c.Sim.Status("beidouburst") == 0 {
			return
		}
		if icd > c.Sim.Frame() {
			c.Log.Debugw("beidou Q (active) on icd", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent)
			return
		}

		d := c.burstSnapshot.Clone()
		//on hit we have to chain
		d.OnHitCallback = c.chainQ(t.Index(), c.Sim.Frame(), 1)

		c.Log.Debugw("beidou Q proc'd", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "actor", ds.Actor, "attack tag", ds.AttackTag)
		c.QueueDmg(&d, 1)

		icd = c.Sim.Frame() + 60 // once per second

	}, "beidou-burst")
}

func (c *char) chainQ(index int, src int, count int) func(t def.Target) {
	if c.Base.Cons > 1 && count == 5 {
		return nil
	}
	if c.Base.Cons < 2 && count == 3 {
		return nil
	}
	//check number of targets, if target < 2 then no bouncing

	//figure out the next target
	l := len(c.Sim.Targets())
	if l < 2 {
		return nil
	}
	index++
	if index >= l {
		index = 0
	}

	//trigger dmg based on a clone of d
	return func(next def.Target) {
		// log.Printf("hit target %v, frame %v, done proc %v, queuing next index: %v\n", next.Index(), c.Sim.Frame(), count, index)
		d := c.burstSnapshot.Clone()
		d.Targets = index
		d.SourceFrame = c.Sim.Frame()
		d.OnHitCallback = c.chainQ(index, src, count+1)
		c.QueueDmg(&d, 1)
	}
}
