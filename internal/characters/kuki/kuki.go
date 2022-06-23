package kuki

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Kuki, NewChar)
}

type char struct {
	*character.Tmpl
	bellActiveUntil   int
	skillHealSnapshot core.Snapshot // Required as both on hit procs and continuous healing need to use this
	c1AoeMod          float64
	skilldur          int
	c4ICD             int
	c6icd             int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 4
	c.BurstCon = 5
	c.SkillCon = 3
	c.A1passive()
	c.A4passive()

	c.c1AoeMod = 2
	c.skilldur = 720
	c.c4ICD = 0
	c.c6icd = 0
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Cons >= 6 {
		c.c6()
	}

	// c.burstICD()
	return &c, nil
}
func (c *char) A1passive() {
	val := make([]float64, core.EndStatType)
	val[core.Heal] = .15
	c.AddMod(core.CharStatMod{
		Key:          "kuki-a1",
		Expiry:       -1,
		AffectedStat: core.Heal, // to avoid infinite loop when calling MaxHP
		Amount: func() ([]float64, bool) {

			if c.HP()/c.MaxHP() <= 0.5 {
				return val, true
			}
			return nil, false
		},
	})

}

func (c *char) A4passive() {
	//TODO: This assumes the dmg bonus works like Yae (multiplicative), however it can be flat (like ZL)
	m := make([]float64, core.EndStatType)
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "kuki-a4",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			// only trigger on elemental art damage
			if atk.Info.AttackTag != core.AttackTagElementalArt {
				return nil, false
			}
			m[core.DmgP] = c.Stat(core.EM) * 0.0025
			return m, true
		},
	})
	//This line only applies if healing is also multiplicative (and needs a skill check) instead I assumed it is flat

	// val := make([]float64, core.EndStatType)
	// c.AddMod(core.CharStatMod{
	// 	Key:    "kuki-a1",
	// 	Expiry: -1,
	// 	Amount: func() ([]float64, bool) {
	// 		val[core.Heal] = c.Stat(core.EM) * 0.0075
	// 		return val, true

	// 	},
	// })

}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
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
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(.3, false, core.TargettableEnemy), f-1, f-1)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge 1",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[0][c.TalentLvlAttack()],
	}
	//TODO: damage frame
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), f-15, f-15)
	ai.Abil = "Charge 2"
	ai.Mult = charge[1][c.TalentLvlAttack()]
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), f-5, f-5)

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	//remove some hp
	if 0.7*(c.HPCurrent/c.MaxHP()) > 0.2 {
		c.HPCurrent = 0.7 * c.HPCurrent
	} else if (c.HPCurrent / c.MaxHP()) > 0.2 { //check if below 20%
		c.HPCurrent = 0.2 * c.MaxHP()
	}
	//TODO: damage frame

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sanctifying Ring",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypePierce,
		Element:    core.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f)

	c.SetCD(core.ActionSkill, a+15*60)  // what's the diff between f and a again? Nice question Yakult
	c.AddTask(c.bellTick(), "bell", 90) //Assuming this executes every 90 frames-1.5s
	c.bellActiveUntil = c.Core.F + c.skilldur
	c.Core.Log.NewEvent("Bell activated", core.LogCharacterEvent, c.Index, "expected end", c.bellActiveUntil, "next expected tick", c.Core.F+90)

	c.Core.Status.AddStatus("kukibell", c.skilldur)

	return f, a
}

func (c *char) bellTick() func() {
	return func() {
		c.Core.Log.NewEvent("Bell ticking", core.LogCharacterEvent, c.Index)

		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Grass Ring of Sanctification",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagElementalArt,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypePierce,
			Element:    core.Electro,
			Durability: 25,
			Mult:       skilldot[c.TalentLvlSkill()],
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), 2, 2)

		c.Core.Health.Heal(core.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.ActiveChar,
			Message: "Grass Ring of Sanctification Healing",
			Src:     skillhealpp[c.TalentLvlSkill()]*c.MaxHP() + skillhealflat[c.TalentLvlSkill()],
			Bonus:   c.skillHealSnapshot.Stats[core.Heal],
		})

		c.Core.Log.NewEvent("Bell ticked", core.LogCharacterEvent, c.Index, "next expected tick", c.Core.F+90, "active", c.bellActiveUntil)
		//trigger damage
		//TODO: Check for snapshots

		//c.Core.Combat.QueueAttackEvent(&ae, 0)
		//check for orb
		//Particle check is 45% for particle
		if c.Core.Rand.Float64() < .45 {
			c.QueueParticle("Kuki", 1, core.Electro, 100) // TODO: idk the particle timing yet fml (or probability)
		}

		//queue up next hit only if next hit bell is still active
		if c.Core.F+90 <= c.bellActiveUntil {
			c.AddTask(c.bellTick(), "Kuki", 90)
		}
	}
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Gyoei Narukami Kariyama Rite",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    c.MaxHP() * burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	count := 7 //can be 11 at low HP
	if (c.HPCurrent / c.MaxHP()) <= 0.5 {
		count = 12
	}
	interval := 2 * 60 / 7

	for i := 0; i < count*interval; i += interval {

		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(c.c1AoeMod, false, core.TargettableEnemy), i)

	}

	c.ConsumeEnergy(55) //TODO: Check if she can be pre-funneled

	c.SetCDWithDelay(core.ActionBurst, 900, 55)
	return f, a
}
