package yoimiya

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Yoimiya, NewChar)
}

type char struct {
	*character.Tmpl
	a1stack  int
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

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.NormalHitNum = 5
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = core.ZoneInazuma

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()
	c.InitCancelFrames()

	c.a1()
	c.onExit()
	c.burstHook()

	if c.Base.Cons > 0 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
}

func (c *char) a1() {
	c.AddMod(core.CharStatMod{
		Key:    "yoimiya-a1",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			val := make([]float64, core.EndStatType)
			if c.Core.Status.Duration("yoimiyaa1") > 0 {
				val[core.PyroP] = float64(c.a1stack) * 0.02
				return val, true
			}
			c.a1stack = 0
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
		if c.a1stack < 10 {
			c.a1stack++
		}
		c.Core.Status.AddStatus("yoimiyaa1", 180)
		// c.a1expiry = c.Core.F + 180 // 3 seconds
		return false
	}, "yoimiya-a1")
}

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	//infusion to normal attack only
	if c.Core.Status.Duration("yoimiyaskill") > 0 && ai.AttackTag == core.AttackTagNormal {
		ai.Element = core.Pyro
		// ds.ICDTag = core.ICDTagNone
		//multiplier
		c.Core.Log.NewEvent("skill mult applied", core.LogCharacterEvent, c.Index, "prev", ai.Mult, "next", skill[c.TalentLvlSkill()]*ai.Mult, "char", c.Index)
		ai.Mult = skill[c.TalentLvlSkill()] * ai.Mult
	}

	return ds
}
