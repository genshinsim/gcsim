package klee

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Klee, NewChar)
}

type char struct {
	*character.Tmpl
	c1Chance float64
	sparkICD int
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
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 3

	c.sparkICD = -1

	c.SetNumCharges(core.ActionSkill, 2)

	return &c, nil
}

func (c *char) Init() {
	c.Tmpl.Init()

	c.a4()
	c.onExitField()
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		if c.Core.Status.Duration("kleespark") > 0 {
			return 0
		}
		return 50
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}

}

func (c *char) a1() {
	if c.Core.F < c.sparkICD {
		return
	}
	if c.Core.Rand.Float64() < 0.5 {
		return
	}
	c.sparkICD = c.Core.F + 60*4
	c.Core.Status.AddStatus("kleespark", 60*30)
	c.Core.Log.NewEvent("klee gained spark", core.LogCharacterEvent, c.Index, "icd", c.sparkICD)
}

func (c *char) a4() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		if !crit {
			return false
		}
		for _, x := range c.Core.Chars {
			x.AddEnergy("klee-a4", 2)
		}
		return false
	}, "kleea1")
}

func (c *char) c1(delay int) {
	if c.Base.Cons < 1 {
		return
	}
	//0.1 base change, + 0.08 every failure
	if c.Core.Rand.Float64() > c.c1Chance {
		//failed
		c.c1Chance += 0.08
		return
	}
	c.c1Chance = 0.1

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sparks'n'Splash C1",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       1.2 * burst[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, delay)

}

// clear klee burst when she leaves the field and handle c4
func (c *char) onExitField() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		// check if burst is active
		if c.Core.Status.Duration("kleeq") <= 0 {
			return false
		}
		c.Core.Status.DeleteStatus("kleeq")

		if c.Base.Cons >= 4 {
			//blow up
			ai := core.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Sparks'n'Splash C4",
				AttackTag:  core.AttackTagNone,
				ICDTag:     core.ICDTagNone,
				ICDGroup:   core.ICDGroupDefault,
				Element:    core.Pyro,
				Durability: 50,
				Mult:       5.55,
			}

			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, 0)
		}
		return false
	}, "klee-exit")
}
