package ayato

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

type char struct {
	*character.Tmpl
	stacks            int
	stacksMax         int
	shunsuikenCounter int
	particleICD       int
	c6ready           bool
}

func init() {
	core.RegisterCharFunc(core.Ayato, NewChar)
}

// test auto build
func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Hydro
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.CharZone = core.ZoneInazuma
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 5
	c.shunsuikenCounter = 3
	c.stacksMax = 4
	c.particleICD = 0
	c.c6ready = false
	c.a2()
	c.a4()
	c.waveFlash()
	c.soukaiKankaHook()

	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.stacksMax = 5
		c.c2()
	}
	if c.Base.Cons >= 6 {
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
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.CharIndex() {
			return false
		}
		if !c.c6ready {
			return false
		}
		atk := args[1].(*core.AttackEvent)
		if atk.Info.AttackTag != core.AttackTagNormal {
			return false
		}
		ai := core.AttackInfo{
			Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
			ActorIndex: c.Index,
			AttackTag:  core.AttackTagNormal,
			ICDTag:     core.ICDTagNormalAttack,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Hydro,
			Durability: 25,
			Mult:       3,
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 20, 20)
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 22, 22)
		// if c.Core.F > c.particleICD {
		// 	c.particleICD = c.Core.F + 112 //best info we have rn
		// 	c.QueueParticle("ayato", 1, core.Hydro, 80)
		// 	if c.Core.Rand.Float64() < 0.5 {
		// 		c.QueueParticle("ayato", 1, core.Hydro, 80)
		// 	}
		// }

		c.Core.Log.NewEvent("ayato c6 proc'd", core.LogCharacterEvent, c.Index)
		c.c6ready = false
		return false
	}, "ayato-c6")
}

func (c *char) a2() {
	c.Core.Events.Subscribe(core.PostSkill, func(args ...interface{}) bool {
		if c.Core.ActiveChar != c.CharIndex() {
			return false
		}
		c.stacks = 2
		c.Core.Log.NewEvent("ayato a2 proc'd", core.LogCharacterEvent, c.Index)
		return false
	}, "ayato-a2")
}

func (c *char) a4() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.03 * c.MaxHP()
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "ayato-a4",
		Expiry: -1,
		Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
			if a.Info.AttackTag != core.AttackTagElementalBurst {
				return nil, false
			}
			return val, true
		},
	})
}

func (c *char) c1() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.4
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "ayato-c1",
		Expiry: -1,
		Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
			if a.Info.AttackTag != core.AttackTagNormal || t.HP()/t.MaxHP() > 0.5 {
				return nil, false
			}
			return val, true
		},
	})
}

func (c *char) c2() {
	val := make([]float64, core.EndStatType)
	val[core.HPP] = 0.4
	c.AddMod(core.CharStatMod{
		Key:    "ayato-c2",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			if c.stacks >= 3 {
				return val, true
			} else {
				return nil, false
			}
		},
	})
}

func (c *char) Snapshot(ai *core.AttackInfo) core.Snapshot {
	ds := c.Tmpl.Snapshot(ai)

	if c.Core.Status.Duration("soukaikanka") > 0 {
		switch ai.AttackTag {
		case core.AttackTagNormal:
		case core.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = core.Hydro
	}
	return ds
}

func (c *char) AdvanceNormalIndex() {
	c.NormalCounter++

	if c.Core.Status.Duration("soukaikanka") > 0 {
		if c.NormalCounter == c.shunsuikenCounter {
			c.NormalCounter = 0
		}
	} else {
		if c.NormalCounter == c.NormalHitNum {
			c.NormalCounter = 0
		}

	}
}
