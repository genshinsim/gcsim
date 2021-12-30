package beidou

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/shield"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, f-1)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	counter := p["counter"]
	f, a := c.ActionFrames(core.ActionSkill, p)
	//0 for base dmg, 1 for 1x bonus, 2 for max bonus
	if counter >= 2 {
		counter = 2
		c.Core.Status.AddStatus("beidoua4", 600)
	}
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Tidecaller (E)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Electro,
		Durability: 50,
		Mult:       skillbase[c.TalentLvlSkill()] + skillbonus[c.TalentLvlSkill()]*float64(counter),
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, f-1)

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
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	if c.Energy < c.EnergyMax {
		c.Core.Log.Debugw("burst insufficient energy; skipping", "frame", c.Core.F, "event", core.LogCharacterEvent, "character", c.Base.Key.String())
		return 0, 0
	}

	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stormbreaker (Q)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Electro,
		Durability: 100,
		Mult:       burstonhit[c.TalentLvlBurst()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 0, f-1)

	c.Core.Status.AddStatus("beidouburst", 900)

	procAI := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stormbreak Proc (Q)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       burstproc[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)
	c.burstAtk = &core.AttackEvent{
		Info:     procAI,
		Snapshot: snap,
	}

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
				Value:    -0.15,
			})
		}
	}

	c.ConsumeEnergy(11)
	c.SetCD(core.ActionBurst, 1200)
	return f, a
}

func (c *char) burstProc() {
	icd := 0
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if c.Core.Status.Duration("beidouburst") == 0 {
			return false
		}
		if icd > c.Core.F {
			c.Core.Log.Debugw("beidou Q (active) on icd", "frame", c.Core.F, "event", core.LogCharacterEvent)
			return false
		}

		//trigger a chain of attacks starting at the first target
		atk := *c.burstAtk
		atk.SourceFrame = c.Core.F
		atk.Pattern = core.NewDefSingleTarget(t.Index(), core.TargettableEnemy)
		cb := c.chain(c.Core.F, 1)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.Combat.QueueAttackEvent(&atk, 1)

		c.Core.Log.Debugw("beidou Q proc'd", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", ae.Info.ActorIndex, "attack tag", ae.Info.AttackTag)

		icd = c.Core.F + 60 // once per second
		return false
	}, "beidou-burst")
}

func (c *char) chain(src int, count int) core.AttackCBFunc {

	if c.Base.Cons > 1 && count == 5 {
		return nil
	}
	if c.Base.Cons < 2 && count == 3 {
		return nil
	}
	return func(a core.AttackCB) {
		//on hit figure out the next target
		trgs := c.Core.EnemyExcl(a.Target.Index())
		if len(trgs) == 0 {
			//do nothing if no other target other than this one
			return
		}
		//otherwise pick a random one
		next := c.Core.Rand.Intn(len(trgs))
		//queue an attack vs next target
		atk := *c.burstAtk
		atk.SourceFrame = src
		atk.Pattern = core.NewDefSingleTarget(trgs[next], core.TargettableEnemy)
		cb := c.chain(src, count+1)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.Combat.QueueAttackEvent(&atk, 1)

	}
}
