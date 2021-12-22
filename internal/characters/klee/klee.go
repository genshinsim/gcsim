package klee

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Klee, NewChar)
}

type char struct {
	*character.Tmpl
	c1Chance     float64
	eCharge      int
	eChargeMax   int
	eNextRecover int
	eTickSrc     int
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
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 3
	c.eChargeMax = 2
	c.eCharge = 2

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
		if c.Tags["spark"] > 0 {
			return 0
		}
		return 50
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}

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
			x.AddEnergy(2)
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

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, delay)

}

func (c *char) c4() {
	c.Core.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		//if burst is active and klee no longer active char
		if c.Core.ActiveChar != c.Index && c.Core.Status.Duration("kleeq") > 0 {
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

			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, 0)
		}
		return false

	}, "klee-c4")
}
