package noelle

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"github.com/genshinsim/gsim/pkg/shield"

	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("noelle", NewChar)
}

type char struct {
	*character.Tmpl
	shieldTimer int
	a4Counter   int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = def.WeaponClassClaymore
	c.NormalHitNum = 4

	c.a2()

	return &c, nil
}

/**

a2: shielding if fall below hp threshold, not implemented

a4: every 4 hit decrease breastplate cd by 1; implement as hook

c1: 100% healing, not implemented

c2: decrease stam consumption, to be implemented

c4: explodes for 400% when expired or destroyed; how to implement expired?

c6: sweeping time increase additional 50%; add 1s up to 10s everytime opponent killed (NOT IMPLEMENTED, NOTHING DIES)

**/

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			f = 28 //frames from keqing lib
		case 1:
			f = 70 - 28
		case 2:
			f = 116 - 70
		case 3:
			f = 174 - 116
		}
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionSkill:
		return 41 //TODO: not ok
	case def.ActionBurst:
		return 111 //ok
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) a2() {
	icd := 0
	c.Sim.AddOnHurt(func(s def.Sim) {
		if c.Sim.Frame() < icd {
			return
		}
		char, _ := s.CharByPos(s.ActiveCharIndex())
		if char.HP()/char.MaxHP() >= 0.3 {
			return
		}
		icd = s.Frame() + 3600
		d := c.Snapshot(
			"A2 Shield",
			def.AttackTagNone,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.NoElement,
			0,
			0,
		)

		//add shield
		x := d.BaseDef*(1+d.Stats[def.DEFP]) + d.Stats[def.DEF]
		s.AddShield(&shield.Tmpl{
			Src:        s.Frame(),
			ShieldType: def.ShieldNoelleA2,
			HP:         4 * x,
			Ele:        def.Cryo,
			Expires:    c.Sim.Frame() + 1200, //20 sec
		})
	})
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
	c.AddTask(func() {
		c.Sim.ApplyDamage(&d)
		//check for healing
		if c.Sim.GetShield(def.ShieldNoelleSkill) != nil {
			prob := healChance[c.TalentLvlSkill()]
			if c.Base.Cons >= 1 && c.Sim.Status("noelleq") > 0 {
				prob = 1
			}
			if c.Sim.Rand().Float64() < prob {
				//heal target
				x := d.BaseDef*(1+d.Stats[def.DEFP]) + d.Stats[def.DEF]
				heal := (shieldHeal[c.TalentLvlSkill()]*x + shieldHealFlat[c.TalentLvlSkill()]) * (1 + d.Stats[def.Heal])
				c.Sim.HealAll(heal)
			}
		}
	}, "noelle auto", f-1)

	c.AdvanceNormalIndex()

	c.a4Counter++
	if c.a4Counter == 4 {
		c.a4Counter = 0
		if c.Cooldown(def.ActionSkill) > 0 {
			c.ReduceActionCooldown(def.ActionSkill, 60)
		}
	}

	return f
}

type noelleShield struct {
	*shield.Tmpl
	c *char
}

func (n *noelleShield) OnExpire() {
	if n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
}

func (n *noelleShield) OnDamage(dmg float64, ele def.EleType, bonus float64) (float64, bool) {
	taken, ok := n.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok && n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
	return taken, ok
}

func (c *char) newShield(base float64, t def.ShieldType, dur int) *noelleShield {
	n := &noelleShield{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.Src = c.Sim.Frame()
	n.Tmpl.ShieldType = t
	n.Tmpl.HP = base
	n.Tmpl.Expires = c.Sim.Frame() + dur
	return n
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)

	d := c.Snapshot(
		"Breastplate (Shield)",
		def.AttackTagNone,
		def.ICDTagNone,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.NoElement,
		0,
		0,
	)
	d.Targets = def.TargetAll

	//add shield first
	x := d.BaseDef*(1+d.Stats[def.DEFP]) + d.Stats[def.DEF]
	shield := shieldFlat[c.TalentLvlSkill()] + shieldDef[c.TalentLvlSkill()]*x

	c.Sim.AddShield(c.newShield(shield, def.ShieldNoelleSkill, 720))

	//activate shield timer, on expiry explode
	c.shieldTimer = c.Sim.Frame() + 720 //12 seconds

	//deal dmg on cast
	d = c.Snapshot(
		"Breastplate",
		def.AttackTagElementalArt,
		def.ICDTagElementalArt,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		50,
		shieldDmg[c.TalentLvlSkill()],
	)
	d.UseDef = true //TODO: test if this is working ok

	c.a4Counter = 0

	c.QueueDmg(&d, f+1)

	if c.Base.Cons >= 4 {
		c.AddTask(func() {
			if c.shieldTimer == c.Sim.Frame() {
				//deal damage
				c.explodeShield()
			}
		}, "noelle shield", 720)
	}

	c.SetCD(def.ActionSkill, 24*60)
	return f //TODO: frame count
}

func (c *char) explodeShield() {
	c.shieldTimer = 0
	d := c.Snapshot(
		"Breastplate",
		def.AttackTagElementalArt,
		def.ICDTagElementalArt,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		50,
		4,
	)
	d.Targets = def.TargetAll
	c.QueueDmg(&d, 1)
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(def.ActionSkill, p)

	c.Sim.AddStatus("noelleq", 900)

	d := c.Snapshot(
		"Sweeping Time (Burst)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, f-10)

	d = c.Snapshot(
		"Sweeping Time (Skill)",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeBlunt,
		def.Geo,
		25,
		burstskill[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, f-5)

	c.SetCD(def.ActionBurst, 900)
	c.Energy = 0
	return f //TODO: frame count
}

func (c *char) Snapshot(name string, a def.AttackTag, icd def.ICDTag, g def.ICDGroup, st def.StrikeType, e def.EleType, d def.Durability, mult float64) def.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	if c.Sim.Status("noelleq") > 0 {

		x := c.Base.Def*(1+ds.Stats[def.DEFP]) + ds.Stats[def.DEF]
		mult := defconv[c.TalentLvlBurst()]
		if c.Base.Cons == 6 {
			mult += 0.5
		}
		fa := mult * x
		c.Log.Debugw("noelle burst", "frame", c.Sim.Frame(), "event", def.LogSnapshotEvent, "total def", x, "atk added", fa, "mult", mult)

		ds.Stats[def.ATK] += fa
		//infusion to attacks only
		switch ds.AttackTag {
		case def.AttackTagNormal:
		case def.AttackTagPlunge:
		case def.AttackTagExtra:
		default:
			return ds
		}
		ds.Element = def.Geo
		ds.Targets = def.TargetAll
	}
	return ds
}
