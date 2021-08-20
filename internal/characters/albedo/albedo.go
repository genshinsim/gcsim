package albedo

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("albedo", NewChar)
}

type char struct {
	*character.Tmpl
	lastConstruct int
	skillSnapshot core.Snapshot
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 40
	c.EnergyMax = 40
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5

	c.skillHook()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 20
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}
}

/**

a2: skill tick deal 25% more dmg if enemy hp < 50%

a4: burst increase party em by 125 for 10s

c1: skill tick regen 1.2 energy

c2: skill tick grant stacks, lasts 30s; each stack increase burst dmg by 30% of def, stack up to 4 times

c4: active member +30% plunge attack in skill field

c6: active protected by crystallize +17% dmg

**/

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(core.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)
	c.AdvanceNormalIndex()

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {

	f := c.ActionFrames(core.ActionCharge, p)

	d := c.Snapshot(
		"Charge 1",
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		charge[0][c.TalentLvlAttack()],
	)
	d2 := d.Clone()
	d2.Abil = "Charge 2"
	d2.Mult = charge[1][c.TalentLvlAttack()]

	c.QueueDmg(&d, f-15) //TODO: damage frame
	c.QueueDmg(&d2, f-5) //TODO: damage frame

	return f
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

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)

	d := c.Snapshot(
		"Abiogenesis: Solar Isotoma",
		core.AttackTagElementalArt,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		skill[c.TalentLvlSkill()],
	)

	c.QueueDmg(&d, f)

	c.skillSnapshot = c.Snapshot(
		"Abiogenesis: Solar Isotoma (Tick)",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		skillTick[c.TalentLvlSkill()],
	)
	c.skillSnapshot.UseDef = true

	//create a construct
	c.Core.Constructs.NewConstruct(c.newConstruct(2100), true) //35 seconds

	c.lastConstruct = c.Core.F

	c.Tags["elevator"] = 1

	c.SetCD(core.ActionSkill, 240)
	return f
}

func (c *char) skillHook() {
	icd := 0
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		if c.Tags["elevator"] == 0 {
			return false
		}
		if c.Core.F < icd {
			return false
		}
		icd = c.Core.F + 120 // every 2 seconds

		d := c.skillSnapshot.Clone()

		if c.Core.Flags.DamageMode && t.HP()/t.MaxHP() < .5 {
			d.Stats[core.DmgP] += 0.25
			c.Core.Log.Debugw("a2 proc'd, dealing extra dmg", "frame", c.Core.F, "event", core.LogCharacterEvent, "hp %", t.HP()/t.MaxHP(), "final dmg", d.Stats[core.DmgP])
		}

		c.QueueDmg(&d, 1)

		//67% chance to generate 1 geo orb
		if c.Core.Rand.Float64() < 0.67 {
			c.QueueParticle("albedo", 1, core.Geo, 100)
		}

		//c1
		if c.Base.Cons >= 1 {
			c.AddEnergy(1.2)
			c.Core.Log.Debugw("c1 restoring energy", "frame", c.Core.F, "event", core.LogCharacterEvent)
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

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)

	hits, ok := p["bloom"]
	if !ok {
		hits = 2 //default 2 hits
	}

	d := c.Snapshot(
		"Rite of Progeniture: Tectonic Tide",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		burst[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f)

	d = c.Snapshot(
		"Rite of Progeniture: Tectonic Tide (Bloom)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		burstPerBloom[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	//check stacks
	if c.Base.Cons >= 2 && c.Core.Status.Duration("albedoc2") > 0 {
		d.FlatDmg += (d.BaseDef*(1+d.Stats[core.DEFP]) + d.Stats[core.DEF]) * float64(c.Tags["c2"])
		c.Tags["c2"] = 0
	}

	for i := 0; i < hits; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f)
	}

	//self buff EM
	for _, char := range c.Core.Chars {
		val := make([]float64, core.EndStatType)
		val[core.EM] = 120
		char.AddMod(core.CharStatMod{
			Key: "albedo-a4",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return val, true
			},
			Expiry: c.Core.F + 600,
		})
	}

	c.SetCD(core.ActionSkill, 720)
	c.Energy = 0
	return f
}

func (c *char) c4() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.3
	c.AddMod(core.CharStatMod{
		Key:    "albedo-c4",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a != core.AttackTagPlunge {
				return nil, false
			}
			if c.Tags["elevator"] != 1 {
				return nil, false
			}
			return val, true
		},
	})
}

func (c *char) c6() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.17
	c.AddMod(core.CharStatMod{
		Key:    "albedo-c6",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
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
