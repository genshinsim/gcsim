package jean

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Jean, NewChar)
}

type char struct {
	*character.Tmpl
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
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
	c.InitCancelFrames()

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	if c.Base.Cons == 6 {
		c.c6()
	}
}

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       auto[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.AddTask(func() {
		snap := c.Snapshot(&ai)
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.4, false, core.TargettableEnemy), 0)

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
	}, "jean-na", f)

	c.AdvanceNormalIndex()

	return f, a
}

// CA has no special interaction with her kit
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.4, false, core.TargettableEnemy), f, f)

	return f, a
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 20
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
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
		c.Core.Log.NewEvent("jean c1 adding 40% dmg", core.LogCharacterEvent, c.Index, "final dmg%", snap.Stats[core.DmgP])
	}

	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(1, false, core.TargettableEnemy), f)

	count := 2
	if c.Core.Rand.Float64() < 2.0/3.0 {
		count++
	}
	c.QueueParticle("Jean", count, core.Anemo, f+100)

	c.SetCDWithDelay(core.ActionSkill, 360, f-2)
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	//p is the number of times enemy enters or exits the field
	enter := p["enter"]
	if enter < 1 {
		enter = 1
	}
	delay, ok := p["enter_delay"]
	if !ok {
		delay = 600 / enter
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

	//initial hit at 40f
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), 40)

	ai.Abil = "Dandelion Breeze (In/Out)"
	ai.Mult = burstEnter[c.TalentLvlBurst()]
	//first enter is at frame 55
	for i := 0; i < enter; i++ {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(5, false, core.TargettableEnemy), 55+i*delay)
	}

	c.Core.Status.AddStatus("jeanq", 600+f)

	if c.Base.Cons >= 4 {
		//add debuff to all target for ??? duration
		for _, t := range c.Core.Targets {
			t.AddResMod("jeanc4", core.ResistMod{
				Duration: 600 + f, //10 seconds + animation
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

	player, ok := c.Core.Targets[0].(*player.Player)
	if !ok {
		panic("target 0 should be Player but is not!!")
	}

	//attack self
	selfSwirl := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dandelion Breeze (Self Swirl)",
		Element:    core.Anemo,
		Durability: 25,
	}

	//duration is 10.5s, first tick start at frame 100, + 60 each
	for i := 100; i < 100+630; i += 60 {
		c.AddTask(func() {
			// c.Core.Log.NewEvent("jean q healing", core.LogCharacterEvent, c.Index, "+heal", hpplus, "atk", atk, "heal amount", healDot)
			c.Core.Health.Heal(core.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.ActiveChar,
				Message: "Dandelion Field",
				Src:     healDot,
				Bonus:   hpplus,
			})

			ae := core.AttackEvent{
				Info:        selfSwirl,
				Pattern:     core.NewDefSingleTarget(0, player.TargetType),
				SourceFrame: c.Core.F,
			}
			c.Core.Log.NewEvent("jean self swirling", core.LogCharacterEvent, c.Index)
			player.ReactWithSelf(&ae)
		}, "Jean Tick", i)
	}

	c.SetCDWithDelay(core.ActionBurst, 1200, 38)
	c.AddTask(func() {
		c.Energy = 16 //jean a4
	}, "jean-burst-energy-consume", 41)

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
	c.Core.Log.NewEvent("jean c6 not implemented", core.LogCharacterEvent, c.Index)
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
			active.AddMod(core.CharStatMod{
				Key:    "jean-c2",
				Amount: func() ([]float64, bool) { return val, true },
				Expiry: c.Core.F + 900,
			})
			c.Core.Log.NewEvent("c2 - adding atk spd", core.LogCharacterEvent, c.Index, "char", c.Index)
		}
	}
}
