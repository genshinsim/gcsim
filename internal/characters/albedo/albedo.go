package albedo

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Albedo, NewChar)
}

type char struct {
	*character.Tmpl
	lastConstruct   int
	skillAttackInfo core.AttackInfo
	skillSnapshot   core.Snapshot
	bloomSnapshot   core.Snapshot
	icdSkill        int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Geo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5

	c.icdSkill = 0

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.skillHook()

	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons == 6 {
		c.c6()
	}
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

/**

a1: skill tick deal 25% more dmg if enemy hp < 50%

a4: burst increase party em by 125 for 10s

c1: skill tick regen 1.2 energy

c2: skill tick grant stacks, lasts 30s; each stack increase burst dmg by 30% of def, stack up to 4 times

c4: active member +30% plunge attack in skill field

c6: active protected by crystallize +17% dmg

**/

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-1, f-1)
	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[0][c.TalentLvlAttack()],
	}
	//TODO: damage frame
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-15, f-15)
	ai.Mult = charge[1][c.TalentLvlAttack()]
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-5, f-5)

	return f, a
}

func (c *char) newConstruct(dur int) core.Construct {
	return &construct{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
	}
}

type construct struct {
	src    int
	expiry int
	char   *char
}

func (c *construct) Key() int {
	return c.src
}

func (c *construct) Type() core.GeoConstructType {
	return core.GeoConstructAlbedoSkill
}

func (c *construct) OnDestruct() {
	c.char.Tags["elevator"] = 0
}

func (c *construct) Expiry() int {
	return c.expiry
}

func (c *construct) IsLimited() bool {
	return true
}

func (c *construct) Count() int {
	return 1
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Abiogenesis: Solar Isotoma",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	//TODO: damage frame
	c.bloomSnapshot = c.Snapshot(&ai)
	c.Core.Combat.QueueAttackWithSnap(ai, c.bloomSnapshot, core.NewDefCircHit(3, false, core.TargettableEnemy), f)

	//snapshot for ticks
	ai.Abil = "Abiogenesis: Solar Isotoma (Tick)"
	ai.ICDTag = core.ICDTagElementalArt
	ai.Mult = skillTick[c.TalentLvlSkill()]
	ai.UseDef = true
	c.skillAttackInfo = ai
	c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)

	// Reset ICD
	c.icdSkill = c.Core.F - 1

	//create a construct
	// Construct is not fully formed until after the hit lands (exact timing unknown)
	c.AddTask(func() {
		c.Core.Constructs.New(c.newConstruct(1800), true)

		c.lastConstruct = c.Core.F

		c.Tags["elevator"] = 1
	}, "albedo-create-construct", f)

	c.SetCD(core.ActionSkill, 240)
	return f, a
}

func (c *char) skillHook() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)
		if c.Tags["elevator"] == 0 {
			return false
		}
		if c.Core.F < c.icdSkill {
			return false
		}
		// Can't be triggered by itself when refreshing
		if atk.Info.Abil == "Abiogenesis: Solar Isotoma" {
			return false
		}

		c.icdSkill = c.Core.F + 120 // every 2 seconds

		snap := c.skillSnapshot

		if c.Core.Flags.DamageMode && t.HP()/t.MaxHP() < .5 {
			snap.Stats[core.DmgP] += 0.25
			c.Core.Log.NewEvent("a1 proc'd, dealing extra dmg", core.LogCharacterEvent, c.Index, "hp %", t.HP()/t.MaxHP(), "final dmg", snap.Stats[core.DmgP])
		}

		c.Core.Combat.QueueAttackWithSnap(c.skillAttackInfo, snap, core.NewDefCircHit(3, false, core.TargettableEnemy), 1)

		//67% chance to generate 1 geo orb
		if c.Core.Rand.Float64() < 0.67 {
			c.QueueParticle("albedo", 1, core.Geo, 100)
		}

		//c1
		if c.Base.Cons >= 1 {
			c.AddEnergy("albedo-c1", 1.2)
			c.Core.Log.NewEvent("c1 restoring energy", core.LogCharacterEvent, c.Index)
		}

		//c2 add stacks
		if c.Base.Cons >= 2 {
			if c.Core.Status.Duration("albedoc2") == 0 {
				c.Tags["c2"] = 0
			}
			c.Core.Status.AddStatus("albedoc2", 1800) //lasts 30 seconds
			c.Tags["c2"]++
			if c.Tags["c2"] > 4 {
				c.Tags["c2"] = 4
			}
		}

		return false

	}, "albedo-skill")
}

func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)

	hits, ok := p["bloom"]
	if !ok {
		hits = 2 //default 2 hits
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rite of Progeniture: Tectonic Tide",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Geo,
		Durability: 25,
		Mult:       burst[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	//check stacks
	if c.Base.Cons >= 2 && c.Core.Status.Duration("albedoc2") > 0 {
		ai.FlatDmg += (snap.BaseDef*(1+snap.Stats[core.DEFP]) + snap.Stats[core.DEF]) * float64(c.Tags["c2"])
		c.Tags["c2"] = 0
	}

	//TODO: damage frame
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(3, false, core.TargettableEnemy), f)

	// Blooms are generated on a slight delay from initial hit
	// TODO: No precise frame data, guessing correct delay
	ai.Abil = "Rite of Progeniture: Tectonic Tide (Blossom)"
	ai.Mult = burstPerBloom[c.TalentLvlSkill()]
	for i := 0; i < hits; i++ {
		c.Core.Combat.QueueAttackWithSnap(ai, c.bloomSnapshot, core.NewDefCircHit(3, false, core.TargettableEnemy), f+30+i*5)
	}

	//Party wide EM buff
	for _, char := range c.Core.Chars {
		val := make([]float64, core.EndStatType)
		val[core.EM] = 125
		char.AddMod(core.CharStatMod{
			Key: "albedo-a4",
			Amount: func() ([]float64, bool) {
				return val, true
			},
			Expiry: c.Core.F + 600,
		})
	}

	c.SetCDWithDelay(core.ActionBurst, 720, 80)
	c.ConsumeEnergy(80)
	return f, a
}

func (c *char) c4() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.3
	for _, char := range c.Core.Chars {
		this := char
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "albedo-c4",
			Expiry: -1,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				if c.Core.ActiveChar != this.CharIndex() {
					return nil, false
				}
				if atk.Info.AttackTag != core.AttackTagPlunge {
					return nil, false
				}
				if c.Tags["elevator"] != 1 {
					return nil, false
				}
				return val, true
			},
		})
	}

}

func (c *char) c6() {

	c.AddMod(core.CharStatMod{
		Key:    "albedo-c6",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			val[core.DmgP] = 0.17
			if c.Tags["elevator"] != 1 {
				return nil, false
			}
			if c.Core.Shields.Get(core.ShieldCrystallize) == nil {
				return nil, false
			}
			return val, true
		},
	})
}
