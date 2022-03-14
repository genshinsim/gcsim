package jean

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Jean, NewChar)
}

type char struct {
	*character.Tmpl
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

func (c *char) ActionFrames(a core.ActionType, p map[string]int) (int, int) {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 14 //frames from keqing lib
		case 1:
			f = 37 - 14
		case 2:
			f = 66 - 37
		case 3:
			f = 124 - 66
		case 4:
			f = 159 - 124
		}
		f = int(float64(f) / (1 + c.Stat(core.AtkSpd)))
		return f, f
	case core.ActionCharge:
		return 73, 73
	case core.ActionSkill:

		hold := p["hold"]
		//hold for p up to 5 seconds
		if hold > 300 {
			hold = 300
		}

		return 23 + hold, 23 + hold
	case core.ActionBurst:
		return 88, 88
	default:
		c.coretype.Log.NewEventBuildMsg(core.LogActionEvent, c.Index, "unknown action (invalid frames): ", a.String())
		return 0, 0
	}
}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  coretype.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       auto[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.AddTask(func() {
		snap := c.Snapshot(&ai)
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.4, false, coretype.TargettableEnemy), 0)

		//check for healing
		if c.Core.Rand.Float64() < 0.5 {
			heal := 0.15 * (snap.BaseAtk*(1+snap.Stats[core.ATKP]) + snap.Stats[core.ATK])
			c.Core.Health.Heal(core.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Wind Companion",
				Src:     heal,
				Bonus:   c.Stat(core.Heal),
			})
		}
	}, "jean-na", f-1)

	c.AdvanceNormalIndex()

	return f, a
}

// CA has no special interaction with her kit
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  coretype.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.4, false, coretype.TargettableEnemy), f-1, f-1)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gale Blade",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	if c.Base.Cons >= 1 && p["hold"] >= 60 {
		//add 40% dmg
		snap.Stats[core.DmgP] += .4
		c.coretype.Log.NewEvent("jean c1 adding 40% dmg", coretype.LogCharacterEvent, c.Index, "final dmg%", snap.Stats[core.DmgP])
	}

	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(1, false, coretype.TargettableEnemy), f-1)

	count := 2
	if c.Core.Rand.Float64() < .67 {
		count++
	}
	c.QueueParticle("Jean", count, core.Anemo, f+100)

	c.SetCD(core.ActionSkill, 360)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	//p is the number of times enemy enters or exits the field
	enter := p["enter"]
	delay, ok := p["delay"]
	if !ok {
		delay = 10
	}

	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dandelion Breeze",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Anemo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, coretype.TargettableEnemy), f-10)

	ai.Abil = "Dandelion Breeze (In/Out)"
	ai.Mult = burstEnter[c.TalentLvlBurst()]
	for i := 0; i < enter; i++ {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, coretype.TargettableEnemy), f+i*delay)
	}

	c.Core.AddStatus("jeanq", 630)

	if c.Base.Cons >= 4 {
		//add debuff to all target for ??? duration
		for _, t := range c.coretype.Targets {
			t.AddResMod("jeanc4", core.ResistMod{
				Duration: 600, //10 seconds
				Ele:      core.Anemo,
				Value:    -0.4,
			})
		}
	}

	//heal on cast
	hpplus := snap.Stats[core.Heal]
	atk := snap.BaseAtk*(1+snap.Stats[core.ATKP]) + snap.Stats[core.ATK]
	heal := burstInitialHealFlat[c.TalentLvlBurst()] + atk*burstInitialHealPer[c.TalentLvlBurst()]
	healDot := burstDotHealFlat[c.TalentLvlBurst()] + atk*burstDotHealPer[c.TalentLvlBurst()]

	c.AddTask(func() {
		c.Core.Health.Heal(core.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Dandelion Breeze",
			Src:     heal,
			Bonus:   hpplus,
		})
	}, "Jean Heal Initial", f)

	//duration is 10.5s
	for i := 60; i < 630; i++ {
		c.AddTask(func() {
			// c.Core.Log.NewEvent("jean q healing", core.LogCharacterEvent, c.Index, "+heal", hpplus, "atk", atk, "heal amount", healDot)
			c.Core.Health.Heal(core.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.ActiveChar,
				Message: "Dandelion Field",
				Src:     healDot,
				Bonus:   hpplus,
			})
		}, "Jean Tick", i)
	}

	c.SetCD(core.ActionBurst, 1200)
	c.AddTask(func() {
		c.Energy = 16 //jean a4
	}, "jean-burst-energy-consume", 46)

	return f, a
}

func (c *char) c6() {
	//reduce dmg by 35% if q active, ignoring the lingering affect
	// c.Sim.AddDRFunc(func() float64 {
	// 	if c.S.StatusActive("jeanq") {
	// 		return 0.35
	// 	}
	// 	return 0
	// })
	c.coretype.Log.NewEvent("jean c6 not implemented", coretype.LogCharacterEvent, c.Index)
}

func (c *char) ReceiveParticle(p core.Particle, isActive bool, partyCount int) {
	c.Tmpl.ReceiveParticle(p, isActive, partyCount)
	if c.Base.Cons >= 2 {
		//only pop this if jean is active
		if !isActive {
			return
		}
		for _, active := range c.Core.Chars {
			val := make([]float64, core.EndStatType)
			val[core.AtkSpd] = 0.15
			active.AddMod(coretype.CharStatMod{
				Key:    "jean-c2",
				Amount: func() ([]float64, bool) { return val, true },
				Expiry: c.Core.Frame + 900,
			})
			c.coretype.Log.NewEvent("c2 - adding atk spd", coretype.LogCharacterEvent, c.Index, "char", c.Index)
		}
	}
}
