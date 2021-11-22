package yoimiya

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("yoimiya", NewChar)
}

type char struct {
	*character.Tmpl
	a2stack  int
	lastPart int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 5
	c.SkillCon = 3

	c.a2()
	c.onExit()
	c.burstHook()

	if c.Base.Cons > 0 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	// if c.Base.Cons == 6 {
	// 	c.c6()
	// }

	//add effect for burst

	return &c, nil
}

func (c *char) a2() {
	val := make([]float64, core.EndStatType)
	c.AddMod(core.CharStatMod{
		Key:    "yoimiya-a2",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if c.Core.Status.Duration("yoimiyaa2") > 0 {
				val[core.Pyro] = float64(c.a2stack) * 0.02
				return val, true
			}
			c.a2stack = 0
			return nil, false
		},
	})
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.Index {
			return false
		}
		if c.Core.Status.Duration("yoimiyaskill") == 0 {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal {
			return false
		}
		//here we can add stacks up to 10
		if c.a2stack < 10 {
			c.a2stack++
		}
		c.Core.Status.AddStatus("yoimiyaa2", 180)
		// c.a2expiry = c.Core.F + 180 // 3 seconds
		return false
	}, "yoimiya-a2")
}

func (c *char) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	//infusion to normal attack only
	if c.Core.Status.Duration("yoimiyaskill") > 0 && ds.AttackTag == core.AttackTagNormal {
		ds.Element = core.Pyro
		// ds.ICDTag = core.ICDTagNone
		//multiplier
		c.Core.Log.Debugw("skill mult applied", "frame", c.Core.F, "event", core.LogCharacterEvent, "prev", ds.Mult, "next", skill[c.TalentLvlSkill()]*ds.Mult, "char", c.Index)
		ds.Mult = skill[c.TalentLvlSkill()] * ds.Mult
	}
	return ds
}
