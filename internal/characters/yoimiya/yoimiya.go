package yoimiya

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Yoimiya, NewChar)
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
	c.Base.Element = core.Pyro
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = core.ZoneInazuma

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
	c.AddMod(core.CharStatMod{
		Key:    "yoimiya-a2",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if c.Core.Status.Duration("yoimiyaa2") > 0 {
				val[core.Pyro] = float64(c.a2stack) * 0.02
				return val, true
			}
			c.a2stack = 0
			return nil, false
		},
	})
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if c.Core.Status.Duration("yoimiyaskill") == 0 {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagNormal {
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

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	//infusion to normal attack only
	if c.Core.Status.Duration("yoimiyaskill") > 0 && ai.AttackTag == core.AttackTagNormal {
		ai.Element = core.Pyro
		// ds.ICDTag = core.ICDTagNone
		//multiplier
		c.Core.Log.Debugw("skill mult applied", "frame", c.Core.F, "event", core.LogCharacterEvent, "prev", ai.Mult, "next", skill[c.TalentLvlSkill()]*ai.Mult, "char", c.Index)
		ai.Mult = skill[c.TalentLvlSkill()] * ai.Mult
	}
	return ds
}
