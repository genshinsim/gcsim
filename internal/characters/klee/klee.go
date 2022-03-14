package klee

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Klee, NewChar)
}

type char struct {
	*character.Tmpl
	c1Chance float64
	sparkICD int
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
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

	c.SetNumCharges(core.ActionSkill, 2)
	c.sparkICD = -1

	c.a4()

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return &c, nil
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
		c.coretype.Log.NewEvent("ActionStam not implemented", coretype.LogActionEvent, c.Index, "action", a.String())
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
	c.Core.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		if !crit {
			return false
		}
		for _, x := range c.Core.Chars {
			x.AddEnergy("klee-a4", 2)
		}
		return false
	}, "kleea2")
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

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 0, delay)

}

func (c *char) c4() {
	c.Core.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		//if burst is active and klee no longer active char
		if c.Core.ActiveChar != c.Index && c.Core.StatusDuration("kleeq") > 0 {
			c.Core.Status.DeleteStatus("kleeq")
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

			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, coretype.TargettableEnemy), 0, 0)
		}
		return false

	}, "klee-c4")
}
