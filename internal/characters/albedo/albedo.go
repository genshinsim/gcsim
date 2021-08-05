package albedo

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("albedo", NewChar)
}

type char struct {
	*character.Tmpl
	lastConstruct int
	skillSnapshot def.Snapshot
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 40
	c.EnergyMax = 40
	c.Weapon.Class = def.WeaponClassSword
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

func (c *char) ActionStam(a def.ActionType, p map[string]int) float64 {
	switch a {
	case def.ActionDash:
		return 18
	case def.ActionCharge:
		return 20
	default:
		c.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
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

	f := c.ActionFrames(def.ActionAttack, p)
	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-1)
	c.AdvanceNormalIndex()

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {

	f := c.ActionFrames(def.ActionCharge, p)

	d := c.Snapshot(
		"Charge 1",
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
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

func (c *char) newConstruct(dur int) def.Construct {
	return &construct{
		src:    c.Sim.Frame(),
		expiry: c.Sim.Frame() + dur,
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

func (c *construct) Type() def.GeoConstructType {
	return def.GeoConstructAlbedoSkill
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
	f := c.ActionFrames(def.ActionSkill, p)

	d := c.Snapshot(
		"Abiogenesis: Solar Isotoma",
		def.AttackTagElementalArt,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		skill[c.TalentLvlSkill()],
	)

	c.QueueDmg(&d, f)

	c.skillSnapshot = c.Snapshot(
		"Abiogenesis: Solar Isotoma (Tick)",
		def.AttackTagElementalArt,
		def.ICDTagElementalArt,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		skillTick[c.TalentLvlSkill()],
	)
	c.skillSnapshot.UseDef = true

	//create a construct
	c.Sim.NewConstruct(c.newConstruct(2100), true) //35 seconds

	c.lastConstruct = c.Sim.Frame()

	c.Tags["elevator"] = 1

	c.SetCD(def.ActionSkill, 240)
	return f
}

func (c *char) skillHook() {
	icd := 0
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		if c.Tags["elevator"] == 0 {
			return
		}
		if c.Sim.Frame() < icd {
			return
		}
		icd = c.Sim.Frame() + 120 // every 2 seconds

		d := c.skillSnapshot.Clone()

		if c.Sim.Flags().HPMode && t.HP()/t.MaxHP() < .5 {
			d.Stats[def.DmgP] += 0.25
			c.Log.Debugw("a2 proc'd, dealing extra dmg", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "hp %", t.HP()/t.MaxHP(), "final dmg", d.Stats[def.DmgP])
		}

		c.QueueDmg(&d, 1)

		//67% chance to generate 1 geo orb
		if c.Sim.Rand().Float64() < 0.67 {
			c.QueueParticle("albedo", 1, def.Geo, 100)
		}

		//c1
		if c.Base.Cons >= 1 {
			c.AddEnergy(1.2)
			c.Log.Debugw("c1 restoring energy", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent)
		}

		//c2 add stacks
		if c.Base.Cons >= 2 {
			if c.Sim.Status("albedoc2") == 0 {
				c.Tags["c2"] = 0
			}
			c.Sim.AddStatus("albedoc2", 1800) //lasts 30 seconds
			c.Tags["c2"]++
			if c.Tags["c2"] > 4 {
				c.Tags["c2"] = 4
			}
		}

		return

	}, "albedo-skill")
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)

	hits, ok := p["bloom"]
	if !ok {
		hits = 2 //default 2 hits
	}

	d := c.Snapshot(
		"Rite of Progeniture: Tectonic Tide",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		burst[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, f)

	d = c.Snapshot(
		"Rite of Progeniture: Tectonic Tide (Bloom)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		burstPerBloom[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	//check stacks
	if c.Base.Cons >= 2 && c.Sim.Status("albedoc2") > 0 {
		d.FlatDmg += (d.BaseDef*(1+d.Stats[def.DEFP]) + d.Stats[def.DEF]) * float64(c.Tags["c2"])
		c.Tags["c2"] = 0
	}

	for i := 0; i < hits; i++ {
		x := d.Clone()
		c.QueueDmg(&x, f)
	}

	//self buff EM
	for _, char := range c.Sim.Characters() {
		val := make([]float64, def.EndStatType)
		val[def.EM] = 120
		char.AddMod(def.CharStatMod{
			Key: "albedo-a4",
			Amount: func(a def.AttackTag) ([]float64, bool) {
				return val, true
			},
			Expiry: c.Sim.Frame() + 600,
		})
	}

	c.SetCD(def.ActionSkill, 720)
	c.Energy = 0
	return f
}

func (c *char) c4() {
	val := make([]float64, def.EndStatType)
	val[def.DmgP] = 0.3
	c.AddMod(def.CharStatMod{
		Key:    "albedo-c4",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if a != def.AttackTagPlunge {
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
	val := make([]float64, def.EndStatType)
	val[def.DmgP] = 0.17
	c.AddMod(def.CharStatMod{
		Key:    "albedo-c6",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if c.Tags["elevator"] != 1 {
				return nil, false
			}
			if c.Sim.GetShield(def.ShieldCrystallize) == nil {
				return nil, false
			}
			return val, true
		},
	})
}
