package noelle

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
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

func NewChar(s core.Sim, log *zap.SugaredLogger, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassClaymore
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

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
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
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionSkill:
		return 41 //TODO: not ok
	case core.ActionBurst:
		return 111 //ok
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) a2() {
	icd := 0
	c.Sim.AddOnHurt(func(s core.Sim) {
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
			core.AttackTagNone,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.NoElement,
			0,
			0,
		)

		//add shield
		x := d.BaseDef*(1+d.Stats[core.DEFP]) + d.Stats[core.DEF]
		s.AddShield(&shield.Tmpl{
			Src:        s.Frame(),
			ShieldType: core.ShieldNoelleA2,
			HP:         4 * x,
			Ele:        core.Cryo,
			Expires:    c.Sim.Frame() + 1200, //20 sec
		})
	})
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
	c.AddTask(func() {
		c.Sim.ApplyDamage(&d)
		//check for healing
		if c.Sim.GetShield(core.ShieldNoelleSkill) != nil {
			prob := healChance[c.TalentLvlSkill()]
			if c.Base.Cons >= 1 && c.Sim.Status("noelleq") > 0 {
				prob = 1
			}
			if c.Sim.Rand().Float64() < prob {
				//heal target
				x := d.BaseDef*(1+d.Stats[core.DEFP]) + d.Stats[core.DEF]
				heal := (shieldHeal[c.TalentLvlSkill()]*x + shieldHealFlat[c.TalentLvlSkill()]) * (1 + d.Stats[core.Heal])
				c.Sim.HealAll(heal)
			}
		}
	}, "noelle auto", f-1)

	c.AdvanceNormalIndex()

	c.a4Counter++
	if c.a4Counter == 4 {
		c.a4Counter = 0
		if c.Cooldown(core.ActionSkill) > 0 {
			c.ReduceActionCooldown(core.ActionSkill, 60)
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

func (n *noelleShield) OnDamage(dmg float64, ele core.EleType, bonus float64) (float64, bool) {
	taken, ok := n.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok && n.c.Base.Cons >= 4 {
		n.c.explodeShield()
	}
	return taken, ok
}

func (c *char) newShield(base float64, t core.ShieldType, dur int) *noelleShield {
	n := &noelleShield{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.Src = c.Sim.Frame()
	n.Tmpl.ShieldType = t
	n.Tmpl.HP = base
	n.Tmpl.Expires = c.Sim.Frame() + dur
	return n
}

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)

	d := c.Snapshot(
		"Breastplate (Shield)",
		core.AttackTagNone,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.NoElement,
		0,
		0,
	)
	d.Targets = core.TargetAll

	//add shield first
	x := d.BaseDef*(1+d.Stats[core.DEFP]) + d.Stats[core.DEF]
	shield := shieldFlat[c.TalentLvlSkill()] + shieldDef[c.TalentLvlSkill()]*x

	c.Sim.AddShield(c.newShield(shield, core.ShieldNoelleSkill, 720))

	//activate shield timer, on expiry explode
	c.shieldTimer = c.Sim.Frame() + 720 //12 seconds

	//deal dmg on cast
	d = c.Snapshot(
		"Breastplate",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
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

	c.SetCD(core.ActionSkill, 24*60)
	return f //TODO: frame count
}

func (c *char) explodeShield() {
	c.shieldTimer = 0
	d := c.Snapshot(
		"Breastplate",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		50,
		4,
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, 1)
}

func (c *char) Burst(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)

	c.Sim.AddStatus("noelleq", 900)

	d := c.Snapshot(
		"Sweeping Time (Burst)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		burst[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f-10)

	d = c.Snapshot(
		"Sweeping Time (Skill)",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeBlunt,
		core.Geo,
		25,
		burstskill[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f-5)

	c.SetCD(core.ActionBurst, 900)
	c.Energy = 0
	return f //TODO: frame count
}

func (c *char) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	if c.Sim.Status("noelleq") > 0 {

		x := c.Base.Def*(1+ds.Stats[core.DEFP]) + ds.Stats[core.DEF]
		mult := defconv[c.TalentLvlBurst()]
		if c.Base.Cons == 6 {
			mult += 0.5
		}
		fa := mult * x
		c.Log.Debugw("noelle burst", "frame", c.Sim.Frame(), "event", core.LogSnapshotEvent, "total def", x, "atk added", fa, "mult", mult)

		ds.Stats[core.ATK] += fa
		//infusion to attacks only
		switch ds.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagPlunge:
		case core.AttackTagExtra:
		default:
			return ds
		}
		ds.Element = core.Geo
		ds.Targets = core.TargetAll
	}
	return ds
}
